package event

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	EventTransferName                  = "Transfer"
	EventNewERC721Universal            = "NewERC721Universal"
	EventMintedWithExternalURI         = "MintedWithExternalURI"
	EventEvolvedWithExternalURI        = "EvolvedWithExternalURI"
	EventNewERC721UniversalSigHash     = generateEventSignatureHash(EventNewERC721Universal, "address", "string")
	EventTransferSigHash               = generateEventSignatureHash(EventTransferName, "address", "address", "uint256")
	EventMintedWithExternalURISigHash  = generateEventSignatureHash(EventMintedWithExternalURI, "address", "uint96", "uint256", "string")
	EventEvolvedWithExternalURISigHash = generateEventSignatureHash(EventEvolvedWithExternalURI, "uint256", "string")
	EventTopicsError                   = fmt.Errorf("unexpected topics length")
	ERC721TransferEventSigHash         = generateEventSignatureHash("Transfer", "address", "address", "uint256")
)

func generateEventSignatureHash(event string, params ...string) string {
	eventSig := []byte(fmt.Sprintf("%s(%s)", event, strings.Join(params, ",")))

	return crypto.Keccak256Hash(eventSig).Hex()
}
