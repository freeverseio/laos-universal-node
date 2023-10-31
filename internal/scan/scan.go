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
	ERC721 "github.com/freeverseio/laos-universal-node/internal/platform/blockchain/erc721"
)

var (
	eventTransferName          = "Transfer"
	eventApprovalName          = "Approval"
	eventApprovalForAllName    = "ApprovalForAll"
	eventTransferSig           = []byte(fmt.Sprintf("%s(address,address,uint256)", eventTransferName))
	eventApprovalSig           = []byte(fmt.Sprintf("%s(address,address,uint256)", eventApprovalName))
	eventApprovalForAllSig     = []byte(fmt.Sprintf("%s(address,address,bool)", eventApprovalForAllName))
	eventTransferSigHash       = crypto.Keccak256Hash(eventTransferSig).Hex()
	eventApprovalSigHash       = crypto.Keccak256Hash(eventApprovalSig).Hex()
	eventApprovalForAllSigHash = crypto.Keccak256Hash(eventApprovalForAllSig).Hex()
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

// Scanner is responsible for scanning and retrieving the ERC721 events
type Scanner interface {
	ScanEvents(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) ([]Event, error)
}

type scanner struct {
	client   EthClient
	contract common.Address
}

// NewScanner instantiates the default implementation for the Scanner interface
func NewScanner(client EthClient, contract common.Address) Scanner {
	return scanner{
		client:   client,
		contract: contract,
	}
}

// ScanEvents returns the ERC721 events between fromBlock and toBlock
func (s scanner) ScanEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]Event, error) {
	eventLogs, err := s.filterEventLogs(ctx, fromBlock, toBlock)
	if err != nil {
		return nil, fmt.Errorf("error filtering events: %w", err)
	}
	if len(eventLogs) == 0 {
		slog.Debug("no events found for block range", "from_block", fromBlock.Int64(), "to_block", toBlock.Int64())
		return nil, nil
	}

	contractAbi, err := abi.JSON(strings.NewReader(ERC721.ERC721MetaData.ABI))
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

func (s scanner) filterEventLogs(ctx context.Context, firstBlock, lastBlock *big.Int) ([]types.Log, error) {
	return s.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: firstBlock,
		ToBlock:   lastBlock,
		Addresses: []common.Address{s.contract},
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

func unpackIntoInterface(e Event, eventName string, contractAbi *abi.ABI, eL *types.Log) error {
	err := contractAbi.UnpackIntoInterface(e, eventName, eL.Data)
	if err != nil {
		return fmt.Errorf("error unpacking the event %s: %w", eventTransferName, err)
	}

	return nil
}
