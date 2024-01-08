package scan

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain"
	"github.com/freeverseio/laos-universal-node/internal/platform/blockchain/contract"
)

var (
	eventTransferName                  = "Transfer"
	eventApprovalName                  = "Approval"
	eventApprovalForAllName            = "ApprovalForAll"
	eventNewERC721Universal            = "NewERC721Universal"
	eventNewCollection                 = "NewCollection"
	eventMintedWithExternalURI         = "MintedWithExternalURI"
	eventEvolvedWithExternalURI        = "EvolvedWithExternalURI"
	eventTransferSigHash               = generateEventSignatureHash(eventTransferName, "address", "address", "uint256")
	eventApprovalSigHash               = generateEventSignatureHash(eventApprovalName, "address", "address", "uint256")
	eventApprovalForAllSigHash         = generateEventSignatureHash(eventApprovalForAllName, "address", "address", "bool")
	eventNewERC721UniversalSigHash     = generateEventSignatureHash(eventNewERC721Universal, "address", "string")
	eventNewCollectionSigHash          = generateEventSignatureHash(eventNewCollection, "address", "address")
	eventMintedWithExternalURISigHash  = generateEventSignatureHash(eventMintedWithExternalURI, "address", "uint96", "uint256", "string")
	eventEvolvedWithExternalURISigHash = generateEventSignatureHash(eventEvolvedWithExternalURI, "uint256", "string")
	eventTopicsError                   = fmt.Errorf("unexpected topics length")
)

// Event is an alias of interface{} and it represents the ERC721 events
type Event interface{}

// EventTransfer is the ERC721 Transfer event
type EventTransfer struct {
	From        common.Address
	To          common.Address
	TokenId     *big.Int
	BlockNumber uint64
	Contract    common.Address
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
	BlockNumber        uint64
}

func (e EventNewERC721Universal) GlobalConsensus() (string, error) {
	// Define a regular expression pattern to match the desired content between parentheses
	pattern := `GlobalConsensus\(([^)]+)\)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find the match in the input string
	match := re.FindStringSubmatch(e.BaseURI)

	if len(match) != 2 {
		return "", fmt.Errorf("no global consensus ID found in base URI: %s", e.BaseURI)
	}

	return match[1], nil
}

func (e EventNewERC721Universal) Parachain() (uint64, error) {
	// Define a regular expression pattern to match the desired content between parentheses
	pattern := `Parachain\(([^)]+)\)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find the match in the input string
	match := re.FindStringSubmatch(e.BaseURI)

	if len(match) != 2 {
		return 0, fmt.Errorf("no parachain ID found in base URI: %s", e.BaseURI)
	}
	parachain, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing parachain value to uint: %w", err)
	}

	return parachain, nil
}

func (e EventNewERC721Universal) CollectionAddress() (common.Address, error) {
	// Define a regular expression pattern to match the desired content between parentheses
	pattern := `AccountKey20\(([^)]+)\)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find the match in the input string
	match := re.FindStringSubmatch(e.BaseURI)

	if len(match) != 2 {
		return common.Address{}, fmt.Errorf("no collection address found in base URI: %s", e.BaseURI)
	}

	return common.HexToAddress(match[1]), nil
}

// EventNewCollecion is the LaosEvolution event emitted when a new collection is created
type EventNewCollecion struct {
	CollectionAddress common.Address
	Owner             common.Address
}

// EventMintedWithExternalURI is the LaosEvolution event emitted when a token is minted
type EventMintedWithExternalURI struct {
	Slot        *big.Int
	To          common.Address
	TokenURI    string
	TokenId     *big.Int
	Contract    common.Address
	BlockNumber uint64
	Timestamp   uint64
}

// EventEvolvedWithExternalURI is the LaosEvolution event emitted when a token metadata is updated
type EventEvolvedWithExternalURI struct {
	TokenId  *big.Int
	TokenURI string
}

func generateEventSignatureHash(event string, params ...string) string {
	eventSig := []byte(fmt.Sprintf("%s(%s)", event, strings.Join(params, ",")))

	return crypto.Keccak256Hash(eventSig).Hex()
}

// Scanner is responsible for scanning and retrieving the ERC721 events
type Scanner interface {
	ScanNewUniversalEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]EventNewERC721Universal, error)
	ScanEvents(ctx context.Context, fromBlock *big.Int, toBlock *big.Int, contracts []string) ([]Event, error)
}

type scanner struct {
	client    blockchain.EthClient
	contracts []common.Address
}

// NewScanner instantiates the default implementation for the Scanner interface
func NewScanner(client blockchain.EthClient, contracts ...string) Scanner {
	scan := scanner{
		client: client,
	}
	for _, c := range contracts {
		scan.contracts = append(scan.contracts, common.HexToAddress(c))
	}
	return scan
}

// TODO decide whether contracts should be a variadic parameter of ScanNewUniversalEvents or not (if so, conversion from []string to []common.Address should be done in config.go)

// ScanEvents returns the ERC721 events between fromBlock and toBlock
func (s scanner) ScanNewUniversalEvents(ctx context.Context, fromBlock, toBlock *big.Int) ([]EventNewERC721Universal, error) {
	slog.Info("scanning universal events", "from_block", fromBlock, "to_block", toBlock)
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

	contracts := make([]EventNewERC721Universal, 0)
	for i := range eventLogs {
		if len(eventLogs[i].Topics) == 0 {
			continue
		}

		if eventLogs[i].Topics[0].Hex() == eventNewERC721UniversalSigHash {
			newERC721Universal, err := parseNewERC721Universal(&eventLogs[i], &contractAbi)
			if err != nil {
				return nil, err
			}
			slog.Info("received event", eventNewERC721Universal, newERC721Universal)

			contracts = append(contracts, newERC721Universal)
		}
	}

	if len(contracts) > 0 {
		slog.Info("universal contracts found", "contracts", len(contracts))
	}

	return contracts, nil
}

// ScanEvents returns the ERC721 events between fromBlock and toBlock
func (s scanner) ScanEvents(ctx context.Context, fromBlock, toBlock *big.Int, contracts []string) ([]Event, error) { // TODO change contracts from []string to ...string
	var addresses []common.Address
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

	erc721UniversalAbi, err := abi.JSON(strings.NewReader(contract.Erc721universalMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("error instantiating ABI: %w", err)
	}

	collectionAbi, err := abi.JSON(strings.NewReader(contract.CollectionMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("error instantiating ABI: %w", err)
	}

	evoAbi, err := abi.JSON(strings.NewReader(contract.EvolutionMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("error instantiating ABI: %w", err)
	}

	var parsedEvents []Event
	for i := range eventLogs {
		if len(eventLogs[i].Topics) > 0 {
			switch eventLogs[i].Topics[0].Hex() {
			// Ownership events
			case eventTransferSigHash:
				transfer, err := parseTransfer(&eventLogs[i], &erc721UniversalAbi)
				if err != nil {
					if err != eventTopicsError {
						return nil, err
					}
					slog.Warn("incorrect number of topics found in Transfer event",
						"topics_found", len(eventLogs[i].Topics),
						"topics_expected", 4)
				} else {
					parsedEvents = append(parsedEvents, transfer)
					slog.Info("received event", eventTransferName, transfer)
				}
			case eventApprovalSigHash:
				approval, err := parseApproval(&eventLogs[i], &erc721UniversalAbi)
				if err != nil {
					if err != eventTopicsError {
						return nil, err
					}
					slog.Warn("incorrect number of topics found in Approval event",
						"topics_found", len(eventLogs[i].Topics),
						"topics_expected", 4)
				} else {
					parsedEvents = append(parsedEvents, approval)
					slog.Info("received event", eventApprovalName, approval)
				}
			case eventApprovalForAllSigHash:
				approvalForAll, err := parseApprovalForAll(&eventLogs[i], &erc721UniversalAbi)
				if err != nil {
					if err != eventTopicsError {
						return nil, err
					}
					slog.Warn("incorrect number of topics found in ApprovalForAll event",
						"topics_found", len(eventLogs[i].Topics),
						"topics_expected", 3)
				} else {
					parsedEvents = append(parsedEvents, approvalForAll)
					slog.Info("received event", eventApprovalForAllName, approvalForAll)
				}
			// Collection event
			case eventNewCollectionSigHash:
				ev, err := parseNewCollection(&eventLogs[i], &collectionAbi)
				if err != nil {
					return nil, err
				}

				parsedEvents = append(parsedEvents, ev)
				slog.Info("received event", eventNewCollection, ev)

			// Evolution events
			case eventMintedWithExternalURISigHash:
				ev, err := parseMintedWithExternalURI(&eventLogs[i], &evoAbi)
				if err != nil {
					return nil, err
				}

				blockNum := eventLogs[i].BlockNumber
				h, err := s.client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
				if err != nil {
					return nil, err
				}

				ev.Contract = eventLogs[i].Address
				ev.BlockNumber = blockNum
				ev.Timestamp = h.Time

				parsedEvents = append(parsedEvents, ev)
				slog.Info("received event", eventMintedWithExternalURI, ev)

			case eventEvolvedWithExternalURISigHash:
				ev, err := parseEvolvedWithExternalURI(&eventLogs[i], &evoAbi)
				if err != nil {
					return nil, err
				}

				parsedEvents = append(parsedEvents, ev)
				slog.Info("received event", eventEvolvedWithExternalURI, ev)
			default:
				slog.Debug("unrecognized event", "event_type", eventLogs[i].Topics[0].String())
			}
		}
	}

	return parsedEvents, nil
}

func (s scanner) filterEventLogs(ctx context.Context, firstBlock, lastBlock *big.Int, contracts ...common.Address) ([]types.Log, error) {
	// TODO optionally filter by topics?
	return s.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: firstBlock,
		ToBlock:   lastBlock,
		Addresses: contracts,
	})
}

func parseTransfer(eL *types.Log, contractAbi *abi.ABI) (EventTransfer, error) {
	var transfer EventTransfer
	if len(eL.Topics) != 4 {
		return transfer, eventTopicsError
	}
	err := unpackIntoInterface(&transfer, eventTransferName, contractAbi, eL)
	if err != nil {
		return transfer, err
	}
	transfer.From = common.HexToAddress(eL.Topics[1].Hex())
	transfer.To = common.HexToAddress(eL.Topics[2].Hex())
	transfer.TokenId = eL.Topics[3].Big()
	transfer.BlockNumber = eL.BlockNumber
	transfer.Contract = eL.Address

	return transfer, nil
}

func parseApproval(eL *types.Log, contractAbi *abi.ABI) (EventApproval, error) {
	var approval EventApproval
	if len(eL.Topics) != 4 {
		return approval, eventTopicsError
	}
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
	if len(eL.Topics) != 3 {
		return approvalForAll, eventTopicsError
	}
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
	newERC721Universal.BlockNumber = eL.BlockNumber

	return newERC721Universal, nil
}

func parseNewCollection(eL *types.Log, contractAbi *abi.ABI) (EventNewCollecion, error) {
	var newCollection EventNewCollecion
	err := unpackIntoInterface(&newCollection, eventNewCollection, contractAbi, eL)
	if err != nil {
		return newCollection, err
	}
	newCollection.Owner = common.HexToAddress(eL.Topics[1].Hex())

	return newCollection, nil
}

func parseEvolvedWithExternalURI(eL *types.Log, contractAbi *abi.ABI) (EventEvolvedWithExternalURI, error) {
	var evolveWithExternalURI EventEvolvedWithExternalURI
	err := unpackIntoInterface(&evolveWithExternalURI, eventEvolvedWithExternalURI, contractAbi, eL)
	if err != nil {
		return evolveWithExternalURI, err
	}
	evolveWithExternalURI.TokenId = eL.Topics[1].Big()

	return evolveWithExternalURI, nil
}

func parseMintedWithExternalURI(eL *types.Log, contractAbi *abi.ABI) (EventMintedWithExternalURI, error) {
	var mintWithExternalURI EventMintedWithExternalURI
	err := unpackIntoInterface(&mintWithExternalURI, eventMintedWithExternalURI, contractAbi, eL)
	if err != nil {
		return mintWithExternalURI, err
	}
	mintWithExternalURI.To = common.HexToAddress(eL.Topics[1].Hex())

	return mintWithExternalURI, nil
}

func unpackIntoInterface(e Event, eventName string, contractAbi *abi.ABI, eL *types.Log) error {
	err := contractAbi.UnpackIntoInterface(e, eventName, eL.Data)
	if err != nil {
		return fmt.Errorf("error unpacking the event %s: %w", eventName, err)
	}

	return nil
}
