mocks:
	mockgen -source=internal/blockchain/ethclient.go -destination=internal/blockchain/mock/ethclient.go -package=mock
	mockgen -source=internal/rpc/server.go -destination=internal/rpc/mock/server.go -package=mockrpc