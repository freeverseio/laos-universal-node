package scanner

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
	"github.com/freeverseio/laos-universal-node/ERC721"
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

type Event interface{}

type EventTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

type EventApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

type EventApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
}

func ScanEvents(cli EthClient, contract common.Address, fromBlock *big.Int, toBlock *big.Int) ([]Event, error) {
	eventLogs, err := filterEventLogs(fromBlock, toBlock, contract, cli)
	if err != nil {
		return nil, fmt.Errorf("error filtering events: %w", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(ERC721.ERC721MetaData.ABI)))
	if err != nil {
		return nil, fmt.Errorf("error instantiating ABI: %w", err)
	}
	var parsedEvents []Event
	for _, eL := range eventLogs {
		slog.Info("scanning event", "block", eL.BlockNumber, "txHash", eL.TxHash)
		switch eL.Topics[0].Hex() {
		case eventTransferSigHash:
			transfer, err := parseTransfer(eL, contractAbi)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, transfer)
			slog.Info("received event", eventTransferName, transfer)
		case eventApprovalSigHash:
			approval, err := parseApproval(eL, contractAbi)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, approval)
			slog.Info("received event", eventApprovalName, approval)
		case eventApprovalForAllSigHash:
			approvalForAll, err := parseApprovalForAll(eL, contractAbi)
			if err != nil {
				return nil, err
			}
			parsedEvents = append(parsedEvents, approvalForAll)
			slog.Info("received event", eventApprovalForAllName, approvalForAll)
		default:
			slog.Warn("unrecognized event", "eventType", eL.Topics[0].String())
		}
	}

	return parsedEvents, nil
}

func filterEventLogs(firstBlock, lastBlock *big.Int, address common.Address, cli EthClient) ([]types.Log, error) {
	return cli.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: firstBlock,
		ToBlock:   lastBlock,
		Addresses: []common.Address{address},
	})
}

func parseTransfer(eL types.Log, contractAbi abi.ABI) (EventTransfer, error) {
	var transfer EventTransfer
	err := unpackIntoInterface(&transfer, contractAbi, eL)
	if err != nil {
		return transfer, err
	}
	transfer.From = common.HexToAddress(eL.Topics[1].Hex())
	transfer.To = common.HexToAddress(eL.Topics[2].Hex())
	transfer.TokenId = eL.Topics[3].Big()
	return transfer, nil
}

func parseApproval(eL types.Log, contractAbi abi.ABI) (EventApproval, error) {
	var approval EventApproval
	err := unpackIntoInterface(&approval, contractAbi, eL)
	if err != nil {
		return approval, err
	}
	approval.Owner = common.HexToAddress(eL.Topics[1].Hex())
	approval.Approved = common.HexToAddress(eL.Topics[2].Hex())
	approval.TokenId = eL.Topics[3].Big()
	return approval, nil
}

func parseApprovalForAll(eL types.Log, contractAbi abi.ABI) (EventApprovalForAll, error) {
	var approvalForAll EventApprovalForAll
	err := unpackIntoInterface(&approvalForAll, contractAbi, eL)
	if err != nil {
		return approvalForAll, err
	}
	approvalForAll.Owner = common.HexToAddress(eL.Topics[1].Hex())
	approvalForAll.Operator = common.HexToAddress(eL.Topics[2].Hex())
	return approvalForAll, nil
}

func unpackIntoInterface(e Event, contractAbi abi.ABI, eL types.Log) error {
	err := contractAbi.UnpackIntoInterface(e, eventTransferName, eL.Data)
	if err != nil {
		return fmt.Errorf("error unpacking the event %s: %w", eventTransferName, err)
	}
	return nil
}
