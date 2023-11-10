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
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain/contract"
	"github.com/freeverseio/laos-universal-node/internal/platform/model"
)

var (
	eventTransferName              = "Transfer"
	eventApprovalName              = "Approval"
	eventApprovalForAllName        = "ApprovalForAll"
	eventNewERC721Universal        = "NewERC721Universal"
	eventTransferSigHash           = generateEventSignatureHash(eventTransferName, "address", "address", "uint256")
	eventApprovalSigHash           = generateEventSignatureHash(eventApprovalName, "address", "address", "uint256")
	eventApprovalForAllSigHash     = generateEventSignatureHash(eventApprovalForAllName, "address", "address", "bool")
	eventNewERC721UniversalSigHash = generateEventSignatureHash(eventNewERC721Universal, "address", "string")
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

// EventNewERC721Universal is the ERC721 event emitted when a new Universal contract is deployed
type EventNewERC721Universal struct {
	NewContractAddress common.Address
	BaseURI            string
}

// TODO decide where this is supposed to go

func generateEventSignatureHash(event string, params ...string) string {
	eventSig := []byte(fmt.Sprintf("%s(%s)", event, strings.Join(params, ",")))

	return crypto.Keccak256Hash(eventSig).Hex()
}

// Scanner is responsible for scanning and retrieving the ERC721 events
type Scanner interface {
	ScanNewUniversalEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]model.ERC721UniversalContract, error)
	ScanEvents(ctx context.Context, fromBlock *big.Int, toBlock *big.Int, contracts []string) ([]Event, error)
}

type scanner struct {
	client    EthClient
	contracts []common.Address
}

// NewScanner instantiates the default implementation for the Scanner interface
func NewScanner(client EthClient, contracts ...string) Scanner {
	scan := scanner{
		client: client,
	}
	for _, c := range contracts {
		scan.contracts = append(scan.contracts, common.HexToAddress(c))
	}
	return scan
}

// ScanEvents returns the ERC721 events between fromBlock and toBlock
func (s scanner) ScanNewUniversalEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]model.ERC721UniversalContract, error) {
	eventLogs, err := s.filterEventLogs(ctx, fromBlock, toBlock, s.contracts...)
	if err != nil {
		return nil, fmt.Errorf("error filtering events: %w", err)
	}

	if len(eventLogs) == 0 {
		slog.Debug("no events found for block range", "from_block", fromBlock.Int64(), "to_block", toBlock.Int64())
		return nil, nil
	}

	contractAbi, err := abi.JSON(strings.NewReader(contract.Erc721universalMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("error instantiating ABI: %w", err)
	}

	contracts := make([]model.ERC721UniversalContract, 0)
	for i := range eventLogs {
		slog.Info("scanning event", "block", eventLogs[i].BlockNumber, "txHash", eventLogs[i].TxHash)
		if len(eventLogs[i].Topics) == 0 {
			continue
		}

		switch eventLogs[i].Topics[0].Hex() {
		case eventNewERC721UniversalSigHash:
			newERC721Universal, err := parseNewERC721Universal(&eventLogs[i], &contractAbi)
			if err != nil {
				return nil, err
			}
			slog.Info("received event", eventNewERC721Universal, newERC721Universal)

			c := model.ERC721UniversalContract{
				Address: newERC721Universal.NewContractAddress,
				BaseURI: newERC721Universal.BaseURI,
			}
			contracts = append(contracts, c)

		default:
			slog.Debug("no new universal contracts found")
		}
	}

	return contracts, nil
}

// ScanEvents returns the ERC721 events between fromBlock and toBlock
func (s scanner) ScanEvents(ctx context.Context, fromBlock, toBlock *big.Int, contracts []string) ([]Event, error) {
	addresses := make([]common.Address, 0)
	for _, c := range contracts {
		addresses = append(addresses, common.HexToAddress(c))
	}

	eventLogs, err := s.filterEventLogs(ctx, fromBlock, toBlock, addresses...)
	if err != nil {
		return nil, fmt.Errorf("error filtering events: %w", err)
	}
	if len(eventLogs) == 0 {
		slog.Debug("no events found for block range", "from_block", fromBlock.Int64(), "to_block", toBlock.Int64())
		return nil, nil
	}

	contractAbi, err := abi.JSON(strings.NewReader(contract.Erc721universalMetaData.ABI))
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
			slog.Warn("unrecognized event", "event_type", eventLogs[i].Topics[0].String())
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

func parseNewERC721Universal(eL *types.Log, contractAbi *abi.ABI) (EventNewERC721Universal, error) {
	var newERC721Universal EventNewERC721Universal
	err := unpackIntoInterface(&newERC721Universal, eventNewERC721Universal, contractAbi, eL)
	if err != nil {
		return newERC721Universal, err
	}

	return newERC721Universal, nil
}

func unpackIntoInterface(e Event, eventName string, contractAbi *abi.ABI, eL *types.Log) error {
	err := contractAbi.UnpackIntoInterface(e, eventName, eL.Data)
	if err != nil {
		return fmt.Errorf("error unpacking the event %s: %w", eventName, err)
	}

	return nil
}
