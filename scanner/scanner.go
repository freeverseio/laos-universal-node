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
	Close()
}

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

func ScanEvents(cli EthClient, contract common.Address, fromBlock *big.Int, toBlock *big.Int) ([]interface{}, error) {
	events, err := filterEvents(fromBlock, toBlock, contract, cli)
	if err != nil {
		return nil, fmt.Errorf("error filtering events: %w", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(ERC721.ERC721MetaData.ABI)))
	if err != nil {
		return nil, fmt.Errorf("error instantiating ABI: %w", err)
	}
	var parsedEvents []interface{}
	for _, e := range events {
		slog.Info("scanning event", "block", e.BlockNumber, "txHash", e.TxHash)
		switch e.Topics[0].Hex() {
		case eventTransferSigHash:
			var transfer EventTransfer
			err := contractAbi.UnpackIntoInterface(&transfer, eventTransferName, e.Data)
			if err != nil {
				return nil, fmt.Errorf("error unpacking the event %s: %w", eventTransferName, err)
			}
			// TODO check e.Topics length?
			transfer.From = common.HexToAddress(e.Topics[1].Hex())
			transfer.To = common.HexToAddress(e.Topics[2].Hex())
			transfer.TokenId = e.Topics[3].Big()

			parsedEvents = append(parsedEvents, transfer)
			slog.Info("received event", eventTransferName, transfer)
		case eventApprovalSigHash:
			var approval EventApproval
			e.Data = nil
			err := contractAbi.UnpackIntoInterface(&approval, eventApprovalName, e.Data)
			if err != nil {
				return nil, fmt.Errorf("error unpacking the event %s: %w", eventApprovalName, err)
			}
			approval.Owner = common.HexToAddress(e.Topics[1].Hex())
			approval.Approved = common.HexToAddress(e.Topics[2].Hex())
			approval.TokenId = e.Topics[3].Big()

			parsedEvents = append(parsedEvents, approval)
			slog.Info("received event", eventApprovalName, approval)
		case eventApprovalForAllSigHash:
			var approvalForAll EventApprovalForAll
			err := contractAbi.UnpackIntoInterface(&approvalForAll, eventApprovalForAllName, e.Data)
			if err != nil {
				return nil, fmt.Errorf("error unpacking the event %s: %w", eventApprovalForAllName, err)
			}
			approvalForAll.Owner = common.HexToAddress(e.Topics[1].Hex())
			approvalForAll.Operator = common.HexToAddress(e.Topics[2].Hex())

			parsedEvents = append(parsedEvents, approvalForAll)
			slog.Info("received event", eventApprovalForAllName, approvalForAll)
		default:
			slog.Warn("unrecognized event", "eventType", e.Topics[0].String())
		}
	}

	return parsedEvents, nil
}

func filterEvents(firstBlock, lastBlock *big.Int, address common.Address, cli EthClient) ([]types.Log, error) {
	return cli.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: firstBlock,
		ToBlock:   lastBlock,
		Addresses: []common.Address{address},
	})
}
