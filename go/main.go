package main

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	internalRpc "github.com/freeverseio/laos-universal-node/cmd/rpc"
)

// System is a type that will be exported as an RPC service.
type System int

type Args struct{}

// SystemResponse holds the result of the Multiply method.
type SystemResponse struct {
	Up int
}

// Up increments the Up field of the SystemResponse struct by 1.
func (a *System) Up(_ *Args, reply *SystemResponse) error {
	reply.Up = 1
	return nil
}

func main() {
	ctx, _ := context.WithCancel(context.Background())


		ethcli, err := ethclient.Dial("url")
		if err != nil {
			log.Printf("failed to connect to Ethereum node: %v", err)
		}
		rpcServer, err := internalRpc.NewServer(ctx, ethcli, common.HexToAddress("0xc4d9faef49ec1e604a76ee78bc992abadaa29527"),  80001)
		
		if err != nil {
			log.Printf("failed to create RPC server: %v", err)
	
		}
		rpcServer.ListenAndServe(ctx, "0.0.0.0:5001")

}             
