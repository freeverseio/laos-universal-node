package model

import "encoding/json"

type ERC721UniversalContract struct {
	BaseURI      string `json:"base_uri"`
	CurrentBlock uint64 `json:"block"`
}

func MarshalERC721UniversalContract(baseURI string, currentBlock uint64) ([]byte, error) {
	c := ERC721UniversalContract{
		BaseURI:      baseURI,
		CurrentBlock: currentBlock,
	}
	return json.Marshal(c)
}

func UnmarshalERC721UniversalContract(d []byte) (*ERC721UniversalContract, error) {
	c := ERC721UniversalContract{}
	err := json.Unmarshal(d, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
