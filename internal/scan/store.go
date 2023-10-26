package scan

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
)

type ERC721BridgelessContract struct {
	Address common.Address `json:"address"`
	Block   uint64         `json:"block"` // TODO how to use this when scanning new events?
	BaseURI string         `json:"base_uri"`
}

type Storage interface {
	Store(ctx context.Context, c ERC721BridgelessContract) error
	ReadAll(ctx context.Context) ([]ERC721BridgelessContract, error)
}

type fsStorage struct {
	file string
}

func NewFSStorage(filename string) (Storage, error) {
	var file *os.File
	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}
	if err = file.Close(); err != nil {
		slog.Warn("error closing storage file", "err", err.Error())
	}

	// Change the file permissions to read and write
	err = os.Chmod(filename, 0o600)
	if err != nil {
		return nil, err
	}

	return fsStorage{file: filename}, nil
}

// Store adds an ERC721BridgelessContract struct in JSON format to the storage file
func (fs fsStorage) Store(ctx context.Context, c ERC721BridgelessContract) error {
	buf, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("error occurred when marshaling ERC721BridgelessContract struct: %w", err)
	}
	f, err := os.OpenFile(fs.file, os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Warn("error closing storage file", "err", err.Error())
		}
	}()

	buf = append(buf, '\n')
	if _, err := f.Write(buf); err != nil {
		return err
	}

	return nil
}

// ReadAll implements FSStorage.
func (fs fsStorage) ReadAll(ctx context.Context) ([]ERC721BridgelessContract, error) {
	f, err := os.Open(fs.file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Warn("error closing storage file", "err", err.Error())
		}
	}()

	contracts := make([]ERC721BridgelessContract, 0)
	var contract ERC721BridgelessContract

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &contract); err != nil {
			return nil, err
		}

		contracts = append(contracts, contract)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}
