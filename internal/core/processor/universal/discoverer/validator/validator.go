package validator

import (
	"fmt"
	"log/slog"

	"github.com/freeverseio/laos-universal-node/internal/platform/model"
	"github.com/freeverseio/laos-universal-node/internal/platform/scan"
)

type Validator interface {
	Validate(event scan.EventNewERC721Universal) (model.ERC721UniversalContract, error)
}

type validator struct {
	globalConsensus string
	parachain       uint64
}

func New(globalConsensus string, parachain uint64) Validator {
	return &validator{
		globalConsensus: globalConsensus,
		parachain:       parachain,
	}
}

func (v *validator) Validate(event scan.EventNewERC721Universal) (model.ERC721UniversalContract, error) {
	contractGlobalConsensus, err := event.GlobalConsensus()
	if err != nil {
		slog.Warn("error parsing global consensus for contract", "contract", event.NewContractAddress,
			"base_uri", event.BaseURI)
		return model.ERC721UniversalContract{}, err
	}
	contractParachain, err := event.Parachain()
	if err != nil {
		slog.Warn("error parsing parachain for contract", "contract", event.NewContractAddress,
			"base_uri", event.BaseURI)
		return model.ERC721UniversalContract{}, err
	}
	collectionAddress, err := event.CollectionAddress()
	if err != nil {
		slog.Warn("error parsing collection address for contract", "contract", event.NewContractAddress,
			"base_uri", event.BaseURI)
		return model.ERC721UniversalContract{}, err
	}

	if contractGlobalConsensus != v.globalConsensus || contractParachain != v.parachain {
		slog.Debug("universal contract's base URI points to a collection in a different evochain, contract discarded",
			"base_uri", event.BaseURI, "chain_global_consensus", v.globalConsensus, "chain_parachain", v.parachain)
		return model.ERC721UniversalContract{}, fmt.Errorf("universal contract's base URI points to a collection in a different evochain, contract discarded")
	}

	return model.ERC721UniversalContract{
		Address:           event.NewContractAddress,
		CollectionAddress: collectionAddress,
		BlockNumber:       event.BlockNumber,
	}, nil
}
