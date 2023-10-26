package scan

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain/ERC721BridgelessMinting"
)

var (
	eventTransferName                      = "Transfer"
	eventApprovalName                      = "Approval"
	eventApprovalForAllName                = "ApprovalForAll"
	eventNewERC721BridgelessMintingName    = "NewERC721BridgelessMinting"
	eventTransferSig                       = []byte(fmt.Sprintf("%s(address,address,uint256)", eventTransferName))
	eventApprovalSig                       = []byte(fmt.Sprintf("%s(address,address,uint256)", eventApprovalName))
	eventApprovalForAllSig                 = []byte(fmt.Sprintf("%s(address,address,bool)", eventApprovalForAllName))
	eventNewERC721BridgelessMintingSig     = []byte(fmt.Sprintf("%s(address,string)", eventNewERC721BridgelessMintingName))
	eventTransferSigHash                   = crypto.Keccak256Hash(eventTransferSig).Hex()
	eventApprovalSigHash                   = crypto.Keccak256Hash(eventApprovalSig).Hex()
	eventApprovalForAllSigHash             = crypto.Keccak256Hash(eventApprovalForAllSig).Hex()
	eventNewERC721BridgelessMintingSigHash = crypto.Keccak256Hash(eventNewERC721BridgelessMintingSig).Hex()
)

// EthClient is an interface for interacting with Ethereum.
// https://github.com/ethereum/go-ethereum/pull/23884
type EthClient interface {
	ethereum.ChainReader
	ethereum.TransactionReader
	ethereum.ChainSyncReader
	ethereum.ContractCaller
	ethereum.LogFilterer
	ethereum.TransactionSender
	ethereum.GasPricer
	ethereum.PendingContractCaller
	ethereum.GasEstimator
	bind.ContractBackend
	ChainID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	Close()
}

// Event is an alias of interface{} and it represents the ERC721 events
type Event interface{}

// EventTransfer is the ERC721 Transfer event
type EventTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

// EventApproval is the ERC721 Approval event
type EventApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

// EventApprovalForAll is the ERC721 ApprovalForAll event
type EventApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
}

// EventNewERC721BridgelessMinting is the ERC721 event emitted when a new Bridgeless Minting contract is deployed
type EventNewERC721BridgelessMinting struct {
	NewContractAddress common.Address
	BaseURI            string
}

// Scanner is responsible for scanning and retrieving the ERC721 events
type Scanner interface {
	ScanNewBridgelessMintingEvents(ctx context.Context, fromBlock, toBlock *big.Int) error
	ScanEvents(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) ([]Event, error)
}

type scanner struct {
	client    EthClient
	contracts []common.Address
	storage   Storage
}

// NewScanner instantiates the default implementation for the Scanner interface
func NewScanner(client EthClient, s Storage, contracts ...string) Scanner {
	scan := scanner{
		client:  client,
		storage: s,
	}
	for _, c := range contracts {
		scan.contracts = append(scan.contracts, common.HexToAddress(c))
	}
	return scan
}

// ScanEvents returns the ERC721 events between fromBlock and toBlock
func (s scanner) ScanNewBridgelessMintingEvents(ctx context.Context, fromBlock, toBlock *big.Int) error {
	triggerDiscovery, err := s.triggerDiscovery(ctx)
	if err != nil {
		return err
	}
	if !triggerDiscovery {
		return nil
	}
	eventLogs, err := s.filterEventLogs(ctx, fromBlock, toBlock, s.contracts...)
	if err != nil {
		return fmt.Errorf("error filtering events: %w", err)
	}
	if len(eventLogs) == 0 {
		slog.Debug("no events found for block range", "from_block", fromBlock.Int64(), "to_block", toBlock.Int64())
		return nil
	}

	contractAbi, err := abi.JSON(strings.NewReader(ERC721BridgelessMinting.ERC721BridgelessMintingMetaData.ABI))
	if err != nil {
		return fmt.Errorf("error instantiating ABI: %w", err)
	}
	for i := range eventLogs {
		slog.Info("scanning event", "block", eventLogs[i].BlockNumber, "txHash", eventLogs[i].TxHash)
		if len(eventLogs[i].Topics) == 0 {
			continue
		}

		switch eventLogs[i].Topics[0].Hex() {
		case eventNewERC721BridgelessMintingSigHash:
			newERC721BridgelessMinting, err := parseNewERC721BridgelessMinting(&eventLogs[i], &contractAbi)
			if err != nil {
				return err
			}
			slog.Info("received event", eventNewERC721BridgelessMintingName, newERC721BridgelessMinting)

			c := ERC721BridgelessContract{
				Address: newERC721BridgelessMinting.NewContractAddress,
				Block:   eventLogs[i].BlockNumber,
				BaseURI: newERC721BridgelessMinting.BaseURI,
			}
			if err := s.storage.Store(ctx, c); err != nil {
				return err
			}
		default:
			slog.Debug("no new bridgeless minting contracts found")
		}
	}

	return nil
}

func (s scanner) triggerDiscovery(ctx context.Context) (bool, error) {
	if len(s.contracts) == 0 {
		return true, nil
	}
	storageContracts, err := s.storage.ReadAll(ctx)
	if err != nil {
		return false, fmt.Errorf("error reading contracts from storage: %w", err)
	}
	/*
	 * When a user provides a list of contracts via flag, we have to discover and
	 * scan those contracts only. For this reason, when we have to determine whether we have
	 * to discover infos about those contracts or not, we will compare if those user-provided contracts
	 * exist in the list of stored contracts (i.e. infos about those contracts, like starting block,
	 * have already been found).
	 * For now, as we don't have a database yet, we only compare that the number of user-provided
	 * contracts matches with the number of stored contracts.
	 */
	if len(storageContracts) == len(s.contracts) {
		return false, nil
	}
	return true, nil
}

// ScanEvents returns the ERC721 events between fromBlock and toBlock
func (s scanner) ScanEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]Event, error) {
	contracts, err := s.storage.ReadAll(context.Background())
	if err != nil {
		return nil, err
	}

	if len(contracts) == 0 {
		slog.Debug("no contracts found", "from_block", fromBlock, "to_block", toBlock)
		return nil, nil
	}

	addresses := make([]common.Address, 0)
	for _, c := range contracts {
		addresses = append(addresses, c.Address)
	}

	eventLogs, err := s.filterEventLogs(ctx, fromBlock, toBlock, addresses...)
	if err != nil {
		return nil, fmt.Errorf("error filtering events: %w", err)
	}
	if len(eventLogs) == 0 {
		slog.Debug("no events found for block range", "from_block", fromBlock.Int64(), "to_block", toBlock.Int64())
		return nil, nil
	}

	contractAbi, err := abi.JSON(strings.NewReader(ERC721BridgelessMinting.ERC721BridgelessMintingMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("error instantiating ABI: %w", err)
	}
	var parsedEvents []Event
	for i := range eventLogs {
		slog.Info("scanning event", "block", eventLogs[i].BlockNumber, "txHash", eventLogs[i].TxHash)
		switch eventLogs[i].Topics[0].Hex() {
		case eventTransferSigHash:
			transfer, err := parseTransfer(&eventLogs[i], &contractAbi)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, transfer)
			slog.Info("received event", eventTransferName, transfer)
		case eventApprovalSigHash:
			approval, err := parseApproval(&eventLogs[i], &contractAbi)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, approval)
			slog.Info("received event", eventApprovalName, approval)
		case eventApprovalForAllSigHash:
			approvalForAll, err := parseApprovalForAll(&eventLogs[i], &contractAbi)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, approvalForAll)
			slog.Info("received event", eventApprovalForAllName, approvalForAll)
		default:
			slog.Warn("unrecognized event", "eventType", eventLogs[i].Topics[0].String())
		}
	}

	return parsedEvents, nil
}

func (s scanner) filterEventLogs(ctx context.Context, firstBlock, lastBlock *big.Int, contracts ...common.Address) ([]types.Log, error) {
	return s.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: firstBlock,
		ToBlock:   lastBlock,
		Addresses: contracts,
	})
}

func parseTransfer(eL *types.Log, contractAbi *abi.ABI) (EventTransfer, error) {
	var transfer EventTransfer
	err := unpackIntoInterface(&transfer, eventTransferName, contractAbi, eL)
	if err != nil {
		return transfer, err
	}
	transfer.From = common.HexToAddress(eL.Topics[1].Hex())
	transfer.To = common.HexToAddress(eL.Topics[2].Hex())
	transfer.TokenId = eL.Topics[3].Big()

	return transfer, nil
}

func parseApproval(eL *types.Log, contractAbi *abi.ABI) (EventApproval, error) {
	var approval EventApproval
	err := unpackIntoInterface(&approval, eventApprovalName, contractAbi, eL)
	if err != nil {
		return approval, err
	}
	approval.Owner = common.HexToAddress(eL.Topics[1].Hex())
	approval.Approved = common.HexToAddress(eL.Topics[2].Hex())
	approval.TokenId = eL.Topics[3].Big()

	return approval, nil
}

func parseApprovalForAll(eL *types.Log, contractAbi *abi.ABI) (EventApprovalForAll, error) {
	var approvalForAll EventApprovalForAll
	err := unpackIntoInterface(&approvalForAll, eventApprovalForAllName, contractAbi, eL)
	if err != nil {
		return approvalForAll, err
	}
	approvalForAll.Owner = common.HexToAddress(eL.Topics[1].Hex())
	approvalForAll.Operator = common.HexToAddress(eL.Topics[2].Hex())

	return approvalForAll, nil
}

func parseNewERC721BridgelessMinting(eL *types.Log, contractAbi *abi.ABI) (EventNewERC721BridgelessMinting, error) {
	var newERC721BridgelessMinting EventNewERC721BridgelessMinting
	err := unpackIntoInterface(&newERC721BridgelessMinting, eventNewERC721BridgelessMintingName, contractAbi, eL)
	if err != nil {
		return newERC721BridgelessMinting, err
	}

	return newERC721BridgelessMinting, nil
}

func unpackIntoInterface(e Event, eventName string, contractAbi *abi.ABI, eL *types.Log) error {
	err := contractAbi.UnpackIntoInterface(e, eventName, eL.Data)
	if err != nil {
		return fmt.Errorf("error unpacking the event %s: %w", eventName, err)
	}

	return nil
}
