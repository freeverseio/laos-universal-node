package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	internalRpc "github.com/freeverseio/laos-universal-node/cmd/rpc"
)

const (
	ethereumEndpoint = "http://localhost:8545"
	rpcAddress       = "0xc4d9faef49ec1e604a76ee78bc992abadaa29527"
	listenAddress    = "0.0.0.0:5001"
	networkID        = 80001
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure resources are cleaned up

	// Set up signal catching
	signals := make(chan os.Signal, 1)
	// Catch all signals since we are okay with a graceful shut down
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle the shutdown signal
	go func() {
		<-signals
		cancel()
	}()

	ethcli, err := ethclient.Dial(ethereumEndpoint)
	if err != nil {
		slog.Error("failed to connect to Ethereum node: %v", err)
	}
	rpcServer, err := internalRpc.NewServer(
		internalRpc.WithEthService(ethcli, common.HexToAddress(rpcAddress), networkID),
		internalRpc.WithNetService(networkID),
		internalRpc.WithSystemHealthService(),
	)
	if err != nil {
		slog.Error("failed to create RPC server: %v", err)
	}

	slog.Info("Starting RPC server", "listenAddress", listenAddress)
	if err := rpcServer.ListenAndServe(ctx, listenAddress); err != nil {
		slog.Error("failed to start RPC server: %v", err)
	}
}
