package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"

	"github.com/freeverseio/laos-universal-node/cmd/server"
	"github.com/freeverseio/laos-universal-node/internal/config"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	badgerStorage "github.com/freeverseio/laos-universal-node/internal/platform/storage/badger"
	"github.com/freeverseio/laos-universal-node/internal/repository"
	"github.com/freeverseio/laos-universal-node/internal/scan"
	"github.com/freeverseio/laos-universal-node/internal/state"
	v1 "github.com/freeverseio/laos-universal-node/internal/state/v1"
)

var version = "undefined"

const (
	historyLength = 256
	klaosChainID  = 2718
)

var mu sync.Mutex

type ReorgError struct {
	block       uint64
	chainHash   common.Hash
	storageHash common.Hash
}

func (e ReorgError) Error() string {
	return "reorg error"
}

func main() {
	if err := run(); err != nil {
		slog.Error("error occurred", "err", err)
	}
}

func run() error {
	c := config.Load()

	setLogger(c.Debug)
	c.LogFields()

	db, err := badger.Open(badger.DefaultOptions(c.Path).WithLoggingLevel(badger.ERROR))
	if err != nil {
		return fmt.Errorf("error initializing storage: %w", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			slog.Error("error closing db", "err", err)
		}
	}()

	// Disclaimer
	slog.Info("******************************************************************************")
	slog.Info("This is a beta version of the Laos Universal Node. It is not intended for production use. Use at your own risk.")
	slog.Info("You are now running the Universal Node Docker Image. Please be aware that this version currently does not handle blockchain reorganizations (reorgs). As a precaution, we strongly encourage operating with a heightened safety margin in your ownership chain management.")
	slog.Info("******************************************************************************")

	storageService := badgerStorage.NewService(db)
	// TODO merge repositoryService and stateService into a single service
	repositoryService := repository.New(storageService)
	stateService := v1.NewStateService(storageService)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)

	// Badger DB garbage collection
	group.Go(func() error {
		numIterations := 3
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				// garbage collection cleans up at most 1 file per iteration
				// https://dgraph.io/docs/badger/get-started/#garbage-collection
				for i := 0; i < numIterations; i++ {
					err := db.RunValueLogGC(0.5)
					if err != nil {
						if err != badger.ErrNoRewrite {
							slog.Error("error occurred while running badger GC", "err", err.Error())
						}
						break
					}
				}
			}
		}
	})

	// Ownership chain scanner
	group.Go(func() error {
		client, err := ethclient.Dial(c.Rpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}
		if err := compareChainIDs(ctx, client, repositoryService); err != nil {
			return err
		}
		s := scan.NewScanner(client, c.Contracts...)
		return scanUniversalChain(ctx, c, client, s, stateService)
	})

	// Evolution chain scanner
	group.Go(func() error {
		client, err := ethclient.Dial(c.EvoRpc)
		if err != nil {
			return fmt.Errorf("error instantiating eth client: %w", err)
		}

		chainID, err := client.ChainID(ctx)
		if err != nil {
			return err
		}
		if chainID.Cmp(big.NewInt(klaosChainID)) == 0 {
			slog.Info("***********************************************************************************************")
			slog.Info("The KLAOS Parachain on Kusama is a test chain for the LAOS Parachain on Polkadot.")
			slog.Info("KLAOS is not endorsed by the LAOS Foundation nor Freeverse")
			slog.Info("for real-value transactions involving the KLAOS token https://www.laosfoundation.io/disclaimer-klaos")
			slog.Info("***********************************************************************************************")
		}

		// TODO check if chain ID match with the one in DB (call "compareChainIDs")
		s := scan.NewScanner(client)
		return scanEvoChain(ctx, c, client, s, stateService)
	})

	// Universal node RPC server
	group.Go(func() error {
		rpcServer, err := server.New()
		if err != nil {
			return fmt.Errorf("failed to create RPC server: %w", err)
		}
		addr := fmt.Sprintf("0.0.0.0:%v", c.Port)
		slog.Info("starting RPC server", "listen_address", addr)
		return rpcServer.ListenAndServe(ctx, c.Rpc, addr, stateService)
	})

	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func compareChainIDs(ctx context.Context, client scan.EthClient, repositoryService repository.Service) error {
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("error getting chain ID from Ethereum client: %w", err)
	}

	chainIdDB, err := repositoryService.GetChainID()
	if err != nil {
		return fmt.Errorf("error getting chain ID from database: %w", err)
	}

	if chainIdDB == "" {
		if err = repositoryService.SetChainID(chainId.String()); err != nil {
			return fmt.Errorf("error setting chain ID in database: %w", err)
		}
	} else if chainId.String() != chainIdDB {
		return fmt.Errorf("mismatched chain IDs: database has %s, eth client reports %s", chainIdDB, chainId.String())
	}
	return nil
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

func scanUniversalChain(ctx context.Context, c *config.Config, client scan.EthClient,
	s scan.Scanner, stateService state.Service,
) error {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	startingBlockDB, err := tx.GetCurrentOwnershipBlock()
	if err != nil {
		return fmt.Errorf("error retrieving the current block from storage: %w", err)
	}
	startingBlock, err := getStartingBlock(ctx, startingBlockDB, c.StartingBlock, client)
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
				l1LatestBlock, err = getL1LatestBlock(ctx, client)
				if err != nil {
					slog.Error("error retrieving the latest block", "err", err.Error())
					break
				}
				lastBlock = calculateLastBlock(startingBlock, l1LatestBlock, c.BlocksRange, c.BlocksMargin)
				if lastBlock < startingBlock {
					slog.Debug("last calculated block is behind starting block, waiting...",
						"lastBlock", lastBlock, "startingBlock", startingBlock)
					waitBeforeNextScan(ctx, c.WaitingTime)
					break
				}
				lastOwnershipBlockTimestamp, err = getTimestampForBlockNumber(ctx, client, lastBlock)
				if err != nil {
					break
				}
			}

			evoSynced, err = isEvoSyncedWithOwnership(stateService, lastOwnershipBlockTimestamp)
			if err != nil {
				slog.Error("error occurred while checking if evolution chain is synced with ownership chain", "err", err.Error())
				break
			}

			if !evoSynced {
				slog.Debug("evolution chain is not synced with ownership chain, waiting...")
				waitBeforeNextScan(ctx, c.WaitingTime)
				break
			}

			err = processUniversalBlockRange(ctx, c, client, stateService, s, startingBlock, lastBlock, lastOwnershipBlockTimestamp)
			if err != nil {
				var reorgErr ReorgError
				if errors.As(err, &reorgErr) {
					slog.Error("ownership chain reorganization detected", "block number", reorgErr.block, "chain hash", reorgErr.chainHash.String(), "storage hash", reorgErr.storageHash.String())
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

func processUniversalBlockRange(ctx context.Context, c *config.Config, client scan.EthClient, stateService state.Service, s scan.Scanner, startingBlock, lastBlock, lastOwnershipBlockTimestamp uint64) error {
	mu.Lock()
	defer mu.Unlock()
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
		return err
	}
	// Store the final block hash to verify in next iteration if a reorganization has taken place.
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
	err = scanAndDigest(ctx, c, s, tx, startingBlock, lastBlock, lastOwnershipBlockTimestamp, client)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("error committing transaction", "err", err.Error())
		return nil
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
		prevIterLastBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(prevIterLastBlockNumber)))
		if err != nil {
			slog.Error("error occurred while retrieving new start range block", "err", err.Error())
			return err
		}

		if prevIterLastBlock.Hash().Cmp(prevLastBlockHash) != 0 {
			return ReorgError{block: startingBlock - 1, chainHash: prevIterLastBlock.Hash(), storageHash: prevLastBlockHash}
		}
	}

	return nil
}

func isEvoSyncedWithOwnership(stateService state.Service, lastOwnershipBlockTimestamp uint64) (bool, error) {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	evoCurrentTimestamp, err := tx.GetCurrentEvoBlockTimestamp()
	if err != nil {
		return false, err
	}

	slog.Debug("check if evo chain is synced with ownership chain",
		"evoCurrentTimestamp", evoCurrentTimestamp, "lastOwnershipBlockTimestamp", lastOwnershipBlockTimestamp)
	if evoCurrentTimestamp < lastOwnershipBlockTimestamp {
		return false, nil
	}

	return true, nil
}

func scanAndDigest(ctx context.Context, c *config.Config, s scan.Scanner, tx state.Tx, startingBlock, lastBlock, lastBlockTimestamp uint64, client scan.EthClient) error {
	shouldDiscover, err := shouldDiscover(tx, c.Contracts)
	if err != nil {
		slog.Error("error occurred reading contracts from storage", "err", err.Error())
		return err
	}
	if shouldDiscover {
		if errDiscover := discoverContracts(ctx, client, s, startingBlock, lastBlock, tx); errDiscover != nil {
			return errDiscover
		}
	}

	var contractsAddress []string

	if len(c.Contracts) > 0 {
		var existingContracts []string
		existingContracts, err = tx.GetExistingERC721UniversalContracts(c.Contracts)
		if err != nil {
			slog.Error("error occurred checking if user-provided contracts exist in storage", "err", err.Error())
			return err
		}
		contractsAddress = append(contractsAddress, existingContracts...)
	} else {
		dbContracts := tx.GetAllERC721UniversalContracts()
		contractsAddress = append(contractsAddress, dbContracts...)
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

func scanEvoChain(ctx context.Context, c *config.Config, client scan.EthClient, s scan.Scanner, stateService state.Service) error {
	tx := stateService.NewTransaction()
	defer tx.Discard()
	startingBlockDB, err := tx.GetCurrentEvoBlock()
	if err != nil {
		return fmt.Errorf("error retrieving the current block from storage: %w", err)
	}
	startingBlock, err := getStartingBlock(ctx, startingBlockDB, c.EvoStartingBlock, client)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		default:
			l1LatestBlock, err := getL1LatestBlock(ctx, client)
			if err != nil {
				slog.Error("error retrieving the latest block", "err", err.Error())
				break
			}
			lastBlock := calculateLastBlock(startingBlock, l1LatestBlock, c.EvoBlocksRange, c.EvoBlocksMargin)
			if lastBlock < startingBlock {
				slog.Debug("scanEvochain, last calculated block is behind starting block, waiting...",
					"lastBlock", lastBlock, "startingBlock", startingBlock)
				waitBeforeNextScan(ctx, c.WaitingTime)
				break
			}

			err = processEvoBlockRange(ctx, client, stateService, s, startingBlock, lastBlock)
			if err != nil {
				var reorgErr ReorgError
				if errors.As(err, &reorgErr) {
					slog.Error("evolution chain reorganization detected", "block number", reorgErr.block, "chain hash", reorgErr.chainHash.String(), "storage hash", reorgErr.storageHash.String())
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

func processEvoBlockRange(ctx context.Context, client scan.EthClient, stateService state.Service, s scan.Scanner, startingBlock, lastBlock uint64) error {
	mu.Lock()
	defer mu.Unlock()
	tx := stateService.NewTransaction()
	defer tx.Discard()

	// retrieve the hash of the final block of the previous iteration.
	prevLastBlockHash, err := tx.GetEvoEndRangeBlockHash()
	if err != nil {
		slog.Error("error occurred while reading LaosEvolution end range block hash", "err", err.Error())
		return nil
	}

	err = verifyChainConsistency(ctx, client, prevLastBlockHash, startingBlock)
	if err != nil {
		return err
	}

	// Retrieve information about the final block in the current block range
	endRangeBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while fetching LaosEvolution end range block", "err", err.Error())
		return nil
	}
	// Store the final block hash to verify in next iteration if a reorganization has taken place.
	if err = tx.SetEvoEndRangeBlockHash(endRangeBlock.Hash()); err != nil {
		slog.Error("error occurred while storing LaosEvolution end range block hash", "err", err.Error())
		return nil
	}

	events, err := s.ScanEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)), nil)
	if err != nil {
		slog.Error("error occurred while scanning LaosEvolution events", "err", err.Error())
		return nil
	}

	err = storeMintEventsAndUpdateBlock(ctx, tx, events, big.NewInt(int64(lastBlock)), client)
	if err != nil {
		return nil
	}

	if err = tx.Commit(); err != nil {
		slog.Error("error committing transaction", "err", err.Error())
		return nil
	}

	return nil
}

func storeMintEventsAndUpdateBlock(ctx context.Context, tx state.Tx, events []scan.Event, lastBlock *big.Int, client scan.EthClient) error {
	err := storeMintedWithExternalURIEventsByContract(tx, events)
	if err != nil {
		slog.Error("error occurred while storing minted events", "err", err.Error())
		return err
	}

	nextStartingBlock := lastBlock.Uint64() + 1
	if err = tx.SetCurrentEvoBlock(nextStartingBlock); err != nil {
		slog.Error("error occurred while storing current block", "err", err.Error())
		return err
	}

	// asking for timestamp of lastBlock as nextStartingBlock does not exist yet
	timestamp, err := getTimestampForBlockNumber(ctx, client, lastBlock.Uint64())
	if err != nil {
		slog.Error("error retrieving block headers", "err", err.Error())
		return err
	}

	if err = tx.SetCurrentEvoBlockTimestamp(timestamp); err != nil {
		slog.Error("error storing block headers", "err", err.Error())
		return err
	}

	return nil
}

func storeMintedWithExternalURIEventsByContract(tx state.Tx, events []scan.Event) error {
	groupedMintEvents := groupEventsMintedWithExternalURIByContract(events)

	for contract, scannedEvents := range groupedMintEvents {
		// fetch current storedEvents stored for this specific contract address
		storedEvents, err := tx.GetMintedWithExternalURIEvents(contract.String())
		if err != nil {
			return err
		}

		ev := make([]model.MintedWithExternalURI, 0)
		if storedEvents != nil {
			ev = append(ev, storedEvents...)
		}
		ev = append(ev, scannedEvents...)
		if err := tx.StoreMintedWithExternalURIEvents(contract.String(), ev); err != nil {
			return err
		}
	}

	return nil
}

// groups events that are of type scan.EventMintedWithExternalURI by contract address
func groupEventsMintedWithExternalURIByContract(events []scan.Event) map[common.Address][]model.MintedWithExternalURI {
	groupMintEvents := make(map[common.Address][]model.MintedWithExternalURI, 0)
	for _, event := range events {
		if e, ok := event.(scan.EventMintedWithExternalURI); ok {
			groupMintEvents[e.Contract] = append(groupMintEvents[e.Contract], model.MintedWithExternalURI{
				Slot:        e.Slot,
				To:          e.To,
				TokenURI:    e.TokenURI,
				TokenId:     e.TokenId,
				BlockNumber: e.BlockNumber,
				Timestamp:   e.Timestamp,
			})
		}
	}
	return groupMintEvents
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

func setLogger(debug bool) {
	// Default slog.Level is Info (0)
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("version", version),
	}))
	slog.SetDefault(logger)
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

func discoverContracts(ctx context.Context, client scan.EthClient, s scan.Scanner, startingBlock, lastBlock uint64, tx state.Tx) error {
	contracts, err := s.ScanNewUniversalEvents(ctx, big.NewInt(int64(startingBlock)), big.NewInt(int64(lastBlock)))
	if err != nil {
		slog.Error("error occurred while discovering new universal events", "err", err.Error())
		return err
	}

	if len(contracts) > 0 {
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
	// TODO if transfer event's timestamp is ahead of global evo chain current block's timestamp => wait X seconds and read again the global evo chain current block from DB
	// we must wait because there might have been mint events whose timestamp is < this transfer event
	// maybe it is worth storing the global evo chain current block's timestamp also?!

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
