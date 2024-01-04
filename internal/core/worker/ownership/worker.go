package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
	"github.com/freeverseio/laos-universal-node/internal/platform/state"
)

const historyLength = 256

type ReorgError struct {
	Block       uint64
	ChainHash   common.Hash
	StorageHash common.Hash
}

func (e ReorgError) Error() string {
	return fmt.Errorf("reorg error: blockNumber %d, chainHash %s, storageHash %s",
		e.Block, e.ChainHash.String(), e.StorageHash.String()).Error()
}

type Worker interface {
	Run(ctx context.Context) error
}

type worker struct {
	waitingTime         time.Duration
	client              scan.EthClient
	stateService        state.Service
	scanner             scan.Scanner
	configStartingBlock uint64
	configBlocksRange   uint64
	configBlocksMargin  uint64
	contracts           []string
	GlobalConsensus     string
	Parachain           uint64
}

func NewWorker(c *config.Config, client scan.EthClient, scanner scan.Scanner, stateService state.Service) Worker {
	return &worker{
		waitingTime:         c.WaitingTime,
		client:              client,
		stateService:        stateService,
		scanner:             scanner,
		configStartingBlock: c.StartingBlock,
		configBlocksRange:   uint64(c.BlocksRange),
		configBlocksMargin:  uint64(c.BlocksMargin),
		contracts:           c.Contracts,
		GlobalConsensus:     c.GlobalConsensus,
		Parachain:           c.Parachain,
	}
}

func (w *worker) Run(ctx context.Context) error {
	tx := w.stateService.NewTransaction()
	defer tx.Discard()
	startingBlockDB, err := tx.GetCurrentOwnershipBlock()
	if err != nil {
		return fmt.Errorf("error retrieving the current block from storage: %w", err)
	}
	startingBlock, err := getStartingBlock(ctx, startingBlockDB, w.configStartingBlock, w.client)
	if err != nil {
		return err
	}
	evoSynced := true
	lastBlock := startingBlock
	var lastOwnershipBlockTimestamp uint64

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			if evoSynced {
				var l1LatestBlock uint64
				l1LatestBlock, err = getL1LatestBlock(ctx, w.client)
				if err != nil {
					slog.Error("error retrieving the latest block", "err", err.Error())
					break
				}
				lastBlock = calculateLastBlock(startingBlock, l1LatestBlock, uint(w.configBlocksRange), uint(w.configBlocksMargin))
				if lastBlock < startingBlock {
					slog.Debug("last calculated block is behind starting block, waiting...",
						"lastBlock", lastBlock, "startingBlock", startingBlock)
					waitBeforeNextScan(ctx, w.waitingTime)
					break
				}
				lastOwnershipBlockTimestamp, err = getTimestampForBlockNumber(ctx, w.client, lastBlock)
				if err != nil {
					break
				}
			}

			evoSynced, err = isEvoSyncedWithOwnership(w.stateService, lastOwnershipBlockTimestamp)
			if err != nil {
				slog.Error("error occurred while checking if evolution chain is synced with ownership chain", "err", err.Error())
				break
			}

			if !evoSynced {
				slog.Debug("evolution chain is not synced with ownership chain, waiting...")
				waitBeforeNextScan(ctx, w.waitingTime)
				break
			}

			err = processUniversalBlockRange(ctx, w.contracts, w.client, w.stateService, w.scanner, startingBlock, lastBlock, lastOwnershipBlockTimestamp, w.GlobalConsensus, w.Parachain)
			if err != nil {
				var reorgErr ReorgError
				if errors.As(err, &reorgErr) {
					slog.Error("ownership chain reorganization detected", "block number", reorgErr.Block, "chain hash", reorgErr.ChainHash.String(), "storage hash", reorgErr.StorageHash.String())
					slog.Info("***********************************************************************************************")
					slog.Info("Please wipe out the database before running the node again.")
					slog.Info("***********************************************************************************************")
					return reorgErr
				}
				break
			}

			startingBlock = lastBlock + 1
		}
	}
}

func shouldDiscover(tx state.Tx, contracts []string) (bool, error) {
	if len(contracts) == 0 {
		return true, nil
	}
	for i := 0; i < len(contracts); i++ {
		hasContract, err := tx.HasERC721UniversalContract(contracts[i])
		if err != nil {
			return false, err
		}
		if !hasContract {
			return true, nil
		}
	}
	return false, nil
}

func processUniversalBlockRange(ctx context.Context,
	contracts []string,
	client scan.EthClient,
	stateService state.Service,
	s scan.Scanner,
	startingBlock,
	lastBlock,
	lastOwnershipBlockTimestamp uint64,
	globalConsensus string,
	parachain uint64,
) error {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	// retrieve the hash of the final block of the previous iteration.
	prevLastBlockHash, err := tx.GetOwnershipEndRangeBlockHash()
	if err != nil {
		slog.Error("error occurred while reading ownership end range block", "err", err.Error())
		return err
	}

	err = verifyChainConsistency(ctx, client, prevLastBlockHash, startingBlock)
	if err != nil {
		return err
	}

	// Retrieve information about the final block in the current block range
	block, err := client.BlockByNumber(ctx, big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred retrieving ownership end range block", "lastBlock", lastBlock, "err", err.Error())
		return err
	}
	// Store the final block hash to verify in next iteration if a reorganization has taken place.
	slog.Debug("setting ownership end range block hash for block number",
		"blockNumber", block.Number(), "blockHash", block.Hash(), "parentHash", block.ParentHash())
	if err = tx.SetOwnershipEndRangeBlockHash(block.Hash()); err != nil {
		slog.Error("error occurred while storing end range block hash", "err", err.Error())
		return err
	}

	// discovering new contracts deployed on the ownership chain
	// choosing which contracts to scan
	// if contracts come from flag, consider only those that have been discovered (whether in this iteration or previously)
	// load merkle trees for all contracts whose events have to be scanned for
	// scanning contracts for events on the ownership chain
	// getting transfer events from scan events
	// retrieving minted events and update the state accordingly
	err = scanAndDigest(ctx, contracts, s, tx, startingBlock, lastBlock, lastOwnershipBlockTimestamp, client, globalConsensus,
		parachain)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("error committing transaction", "err", err.Error())
		return err
	}

	return nil
}

func verifyChainConsistency(ctx context.Context, client scan.EthClient, prevLastBlockHash common.Hash, startingBlock uint64) error {
	// During the initial iteration, no hash is stored in the database, so this code block is bypassed.
	// Verify whether the hash of the last block from the previous iteration remains unchanged; if it differs,
	// it indicates a reorganization has taken place.
	if prevLastBlockHash != (common.Hash{}) {
		var prevIterLastBlock *types.Block
		prevIterLastBlockNumber := startingBlock - 1
		slog.Debug("verifying chain consistency on block number", "lastBlock", prevIterLastBlockNumber)
		prevIterLastBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(prevIterLastBlockNumber)))
		if err != nil {
			slog.Error("error occurred while retrieving new start range block", "err", err.Error())
			return err
		}

		if prevIterLastBlock.Hash().Cmp(prevLastBlockHash) != 0 {
			return ReorgError{Block: startingBlock - 1, ChainHash: prevIterLastBlock.Hash(), StorageHash: prevLastBlockHash}
		}
	}

	return nil
}

func isEvoSyncedWithOwnership(stateService state.Service, lastOwnershipBlockTimestamp uint64) (bool, error) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	lastEvoBLockData, err := tx.GetLastEvoBlock()
	if err != nil {
		return false, err
	}

	slog.Debug("IsEvoSyncedWithOwnership", "evo_bock_number", lastEvoBLockData.Number, "evo_block_timestamp", lastEvoBLockData.Timestamp,
		"evo_block_hash", lastEvoBLockData.Hash)

	evoCurrentTimestamp := lastEvoBLockData.Timestamp
	slog.Debug("check if evo chain is synced with ownership chain",
		"evoCurrentTimestamp", evoCurrentTimestamp, "lastOwnershipBlockTimestamp", lastOwnershipBlockTimestamp)
	if evoCurrentTimestamp < lastOwnershipBlockTimestamp {
		return false, nil
	}

	return true, nil
}

func scanAndDigest(ctx context.Context,
	contracts []string,
	s scan.Scanner,
	tx state.Tx,
	startingBlock,
	lastBlock,
	lastBlockTimestamp uint64,
	client scan.EthClient,
	globalConsensus string,
	parachain uint64,
) error {
	shouldDiscover, err := shouldDiscover(tx, contracts)
	if err != nil {
		slog.Error("error occurred reading contracts from storage", "err", err.Error())
		return err
	}
	if shouldDiscover {
		errDiscover := discoverContracts(ctx, client, s, startingBlock, lastBlock, tx, globalConsensus, parachain)
		if errDiscover != nil {
			return errDiscover
		}
	}

	var contractsAddress []string

	if len(contracts) > 0 {
		var existingContracts []string
		existingContracts, err = tx.GetExistingERC721UniversalContracts(contracts)
		if err != nil {
			slog.Error("error occurred checking if user-provided contracts exist in storage", "err", err.Error())
			return err
		}
		contractsAddress = existingContracts
	} else {
		contractsAddress = tx.GetAllERC721UniversalContracts()
	}

	if err = loadMerkleTrees(tx, contractsAddress); err != nil {
		slog.Error("error creating merkle trees", "err", err)
		return err
	}

	if len(contractsAddress) > 0 {
		var scanEvents []scan.Event
		scanEvents, err = s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), contractsAddress)
		if err != nil {
			slog.Error("error occurred while scanning events", "err", err.Error())
			return err
		}

		var modelTransferEvents map[string][]model.ERC721Transfer
		modelTransferEvents, err = getModelTransferEvents(ctx, client, scanEvents)
		if err != nil {
			slog.Error("error parsing transfer events", "err", err.Error())
			return err
		}

		if err = readEventsAndUpdateState(ctx, client, contractsAddress, modelTransferEvents, tx, lastBlockTimestamp); err != nil {
			slog.Error("error occurred", "err", err.Error())
			return err
		}
	}

	nextStartingBlock := lastBlock + 1

	if err = tagRootsUntilBlock(tx, contractsAddress, nextStartingBlock); err != nil {
		slog.Error("error occurred while tagging roots", "err", err.Error())
		return err
	}

	if err = tx.SetCurrentOwnershipBlock(nextStartingBlock); err != nil {
		slog.Error("error occurred while storing current block", "err", err.Error())
		return err
	}

	return nil
}

func getStartingBlock(ctx context.Context, startingBlockDB, configStartingBlock uint64, client scan.EthClient) (uint64, error) {
	var startingBlock uint64
	var err error
	if startingBlockDB != 0 {
		startingBlock = startingBlockDB
		slog.Debug("ignoring user provided starting block, using last updated block from storage", "starting_block", startingBlock)
	}

	if startingBlock == 0 {
		startingBlock = configStartingBlock
		if startingBlock == 0 {
			startingBlock, err = getL1LatestBlock(ctx, client)
			if err != nil {
				return 0, fmt.Errorf("error retrieving the latest block from chain: %w", err)
			}
			slog.Debug("latest block found", "latest_block", startingBlock)
		}
	}
	return startingBlock, nil
}

func waitBeforeNextScan(ctx context.Context, waitingTime time.Duration) {
	timer := time.NewTimer(waitingTime)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}

func getL1LatestBlock(ctx context.Context, client scan.EthClient) (uint64, error) {
	l1LatestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return l1LatestBlock, nil
}

func calculateLastBlock(startingBlock, l1LatestBlock uint64, blocksRange, blocksMargin uint) uint64 {
	return min(startingBlock+uint64(blocksRange), l1LatestBlock-uint64(blocksMargin))
}

func readEventsAndUpdateState(ctx context.Context, client scan.EthClient, contractsAddress []string, modelTransferEvents map[string][]model.ERC721Transfer, tx state.Tx, lastBlockTimestamp uint64) error {
	for i := range contractsAddress {
		var mintedEvents []model.MintedWithExternalURI
		collectionAddress, err := tx.GetCollectionAddress(contractsAddress[i]) // get collection address from ownership address
		if err != nil {
			return fmt.Errorf("error occurred retrieving the collection address from the ownership contract %s: %w", contractsAddress[i], err)
		}

		// now we get the minted events from the evolution chain for the collection address
		mintedEvents, err = tx.GetMintedWithExternalURIEvents(collectionAddress.String())
		if err != nil {
			return fmt.Errorf("error occurred retrieving evochain minted events for ownership contract %s and collection address %s: %w",
				contractsAddress[i], collectionAddress.String(), err)
		}

		var evoIndex uint64
		evoIndex, err = updateState(ctx, client, mintedEvents, modelTransferEvents[contractsAddress[i]], contractsAddress[i], tx, lastBlockTimestamp)
		if err != nil {
			return fmt.Errorf("error updating state: %w", err)
		}

		if err = tx.SetCurrentEvoEventsIndexForOwnershipContract(contractsAddress[i], evoIndex); err != nil {
			return fmt.Errorf("error updating current evochain index %d for ownership contract %s: %w", evoIndex, contractsAddress[i], err)
		}
	}
	return nil
}

func tagRootsUntilBlock(tx state.Tx, contractsAddress []string, blockNumber uint64) error {
	for i := range contractsAddress {
		lastTaggedBlock, err := tx.GetLastTaggedBlock(common.HexToAddress(contractsAddress[i]))
		if err != nil {
			return err
		}

		for block := lastTaggedBlock + 1; block < int64(blockNumber); block++ {
			if err := tx.TagRoot(common.HexToAddress(contractsAddress[i]), block); err != nil {
				return err
			}
			if (block - historyLength) > 0 {
				if err := tx.DeleteRootTag(common.HexToAddress(contractsAddress[i]), block-historyLength); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func loadMerkleTrees(tx state.Tx, contractsAddress []string) error {
	for i := range contractsAddress {
		if err := loadMerkleTree(tx, common.HexToAddress(contractsAddress[i])); err != nil {
			return err
		}
	}
	return nil
}

func getValidUniversalContracts(globalConsensus string, parachain uint64, events []scan.EventNewERC721Universal) []model.ERC721UniversalContract {
	contracts := make([]model.ERC721UniversalContract, 0)
	for _, e := range events {
		contractGlobalConsensus, err := e.GlobalConsensus()
		if err != nil {
			slog.Warn("error parsing collection address for contract", "contract", e.NewContractAddress,
				"base_uri", e.BaseURI)
			continue
		}
		contractParachain, err := e.Parachain()
		if err != nil {
			slog.Warn("error parsing collection address for contract", "contract", e.NewContractAddress,
				"base_uri", e.BaseURI)
			continue
		}
		collectionAddress, err := e.CollectionAddress()
		if err != nil {
			slog.Warn("error parsing collection address for contract", "contract", e.NewContractAddress,
				"base_uri", e.BaseURI)
			continue
		}

		if contractGlobalConsensus != globalConsensus || contractParachain != parachain {
			slog.Debug("universal contract's base URI points to a collection in a different evochain, contract discarded",
				"base_uri", e.BaseURI, "chain_global_consensus", globalConsensus, "chain_parachain", parachain)
			continue
		}

		contract := model.ERC721UniversalContract{
			Address:           e.NewContractAddress,
			CollectionAddress: collectionAddress,
			BlockNumber:       e.BlockNumber,
		}

		contracts = append(contracts, contract)
	}

	return contracts
}

func discoverContracts(ctx context.Context,
	client scan.EthClient,
	s scan.Scanner,
	startingBlock,
	lastBlock uint64,
	tx state.Tx,
	globalConsensus string,
	parachain uint64,
) error {
	scannedContracts, err := s.ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while discovering new universal events", "err", err.Error())
		return err
	}

	contracts := make([]model.ERC721UniversalContract, 0)
	if len(scannedContracts) > 0 {
		contracts = getValidUniversalContracts(globalConsensus, parachain, scannedContracts)
		if err = tx.StoreERC721UniversalContracts(contracts); err != nil {
			slog.Error("error occurred while storing universal contract(s)", "err", err.Error())
			return err
		}
	}

	for i := range contracts {
		if err = loadMerkleTree(tx, contracts[i].Address); err != nil {
			slog.Error("error creating merkle trees for newly discovered universal contract(s)", "err", err)
			return err
		}

		// check if there are mint events for this contract
		mintEvents, err := tx.GetMintedWithExternalURIEvents(contracts[i].CollectionAddress.String())
		if err != nil {
			slog.Error("error occurred retrieving evochain minted events for ownership contract: %w", err)
			return err
		}

		timestampContract, err := getTimestampForBlockNumber(ctx, client, contracts[i].BlockNumber)
		if err != nil {
			return err
		}

		ownershipContractEvoEventIndex, err := updateStateWithMintEvents(contracts[i].Address, tx, mintEvents, timestampContract)
		if err != nil {
			slog.Error("error occurred updating state with mint events", "err", err)
			return err
		}

		if err = tx.SetCurrentEvoEventsIndexForOwnershipContract(contracts[i].Address.String(), ownershipContractEvoEventIndex); err != nil {
			return fmt.Errorf("error updating current evochain event index %d for ownership contract %s: %w",
				ownershipContractEvoEventIndex, strings.ToLower(contracts[i].Address.String()), err)
		}

		if err = tx.TagRoot(contracts[i].Address, int64(contracts[i].BlockNumber)); err != nil {
			slog.Error("error occurred tagging roots for newly discovered universal contract(s)", "err", err.Error())
			return err
		}
	}

	return nil
}

func updateStateWithMintEvents(contract common.Address, tx state.Tx, mintedEvents []model.MintedWithExternalURI, timestampContract uint64) (uint64, error) {
	for i := range mintedEvents {
		if mintedEvents[i].Timestamp > timestampContract {
			return uint64(i), nil
		}
		if err := tx.Mint(contract, &mintedEvents[i]); err != nil {
			return 0, fmt.Errorf("error updating mint state for contract %s and token id %d: %w",
				contract, mintedEvents[i].TokenId, err)
		}
	}
	return uint64(len(mintedEvents)), nil
}

func loadMerkleTree(tx state.Tx, contractAddress common.Address) error {
	if !tx.IsTreeSetForContract(contractAddress) {
		ownership, enumerated, enumeratedTotal, err := tx.CreateTreesForContract(contractAddress)
		if err != nil {
			return err
		}
		tx.SetTreesForContract(contractAddress, ownership, enumerated, enumeratedTotal)
	}
	return nil
}

func updateState(ctx context.Context, client scan.EthClient, mintedEvents []model.MintedWithExternalURI, modelTransferEvents []model.ERC721Transfer, contract string, tx state.Tx, lastBlockTimestamp uint64) (uint64, error) {
	ownershipContractEvoEventIndex, err := tx.GetCurrentEvoEventsIndexForOwnershipContract(contract)
	if err != nil {
		return 0, err
	}
	var transferIndex int
	for {
		switch {
		// all events have been processed => return
		case ownershipContractEvoEventIndex >= uint64(len(mintedEvents)) && transferIndex >= len(modelTransferEvents):
			return ownershipContractEvoEventIndex, nil
		// all minted events have been processed => process remaining transfer events
		case ownershipContractEvoEventIndex >= uint64(len(mintedEvents)):
			if err := updateStateWithTransfer(contract, tx, &modelTransferEvents[transferIndex]); err != nil {
				return 0, err
			}
			transferIndex++
		// all transfer events have been processed => process remaining minted events
		case transferIndex >= len(modelTransferEvents):
			if mintedEvents[ownershipContractEvoEventIndex].Timestamp < lastBlockTimestamp {
				err := updateStateWithMint(ctx, client, contract, tx, &mintedEvents[ownershipContractEvoEventIndex])
				if err != nil {
					return 0, err
				}
				ownershipContractEvoEventIndex++
			} else {
				return ownershipContractEvoEventIndex, nil
			}

		default:
			// if minted event's timestamp is behind transfer event's timestamp => process minted event
			if mintedEvents[ownershipContractEvoEventIndex].Timestamp < modelTransferEvents[transferIndex].Timestamp {
				err := updateStateWithMint(ctx, client, contract, tx, &mintedEvents[ownershipContractEvoEventIndex])
				if err != nil {
					return 0, err
				}
				ownershipContractEvoEventIndex++
			} else {
				if err := updateStateWithTransfer(contract, tx, &modelTransferEvents[transferIndex]); err != nil {
					return 0, err
				}
				transferIndex++
			}
		}
	}
}

func updateStateWithTransfer(contract string, tx state.Tx, modelTransferEvent *model.ERC721Transfer) error {
	lastTaggedBlock, err := tx.GetLastTaggedBlock(common.HexToAddress(contract))
	if err != nil {
		return err
	}

	for block := lastTaggedBlock + 1; block < int64(modelTransferEvent.BlockNumber); block++ {
		if err := tx.TagRoot(common.HexToAddress(contract), block); err != nil {
			return err
		}
		if (block - historyLength) > 0 {
			if err := tx.DeleteRootTag(common.HexToAddress(contract), block-historyLength); err != nil {
				return err
			}
		}
	}

	if err := tx.Transfer(common.HexToAddress(contract), modelTransferEvent); err != nil {
		return fmt.Errorf("error updating transfer state for contract %s and token id %d, from %s, to %s: %w",
			contract, modelTransferEvent.TokenId,
			modelTransferEvent.From, modelTransferEvent.To, err)
	}
	return nil
}

func updateStateWithMint(ctx context.Context, client scan.EthClient, contract string, tx state.Tx, mintedEvent *model.MintedWithExternalURI) error {
	lastTaggedBlock, err := tx.GetLastTaggedBlock(common.HexToAddress(contract))
	if err != nil {
		return err
	}

	for {
		blockToTag := lastTaggedBlock + 1
		timestamp, err := getTimestampForBlockNumber(ctx, client, uint64(blockToTag))
		if err != nil {
			return err
		}
		if timestamp >= mintedEvent.Timestamp {
			break
		}
		if err := tx.TagRoot(common.HexToAddress(contract), blockToTag); err != nil {
			return err
		}
		if err := tx.DeleteRootTag(common.HexToAddress(contract), blockToTag-historyLength); err != nil {
			return err
		}
		// update lastTaggedBlock
		lastTaggedBlock = blockToTag
	}

	if err := tx.Mint(common.HexToAddress(contract), mintedEvent); err != nil {
		return fmt.Errorf("error updating mint state for contract %s and token id %d: %w",
			contract, mintedEvent.TokenId, err)
	}

	return nil
}

func getModelTransferEvents(ctx context.Context, client scan.EthClient, scanEvents []scan.Event) (map[string][]model.ERC721Transfer, error) {
	modelTransferEvents := make(map[string][]model.ERC721Transfer)
	for i := range scanEvents {
		if scanEvent, ok := scanEvents[i].(scan.EventTransfer); ok {
			timestamp, err := getTimestampForBlockNumber(ctx, client, scanEvent.BlockNumber)
			if err != nil {
				return nil, fmt.Errorf("error retrieving timestamp for block number %d: %w", scanEvent.BlockNumber, err)
			}
			contractString := strings.ToLower(scanEvent.Contract.String())
			modelTransferEvents[contractString] = append(modelTransferEvents[contractString], model.ERC721Transfer{
				From:        scanEvent.From,
				To:          scanEvent.To,
				TokenId:     scanEvent.TokenId,
				BlockNumber: scanEvent.BlockNumber,
				Contract:    scanEvent.Contract,
				Timestamp:   timestamp,
			})
		}
	}
	return modelTransferEvents, nil
}

func getTimestampForBlockNumber(ctx context.Context, client scan.EthClient, blockNumber uint64) (uint64, error) {
	header, err := client.HeaderByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return 0, err
	}
	return header.Time, err
}
