mocks:
	mockgen -source=internal/platform/blockchain/ethclient.go -destination=internal/platform/blockchain/mock/ethclient.go -package=mock
	mockgen -source=internal/rpc/server.go -destination=internal/rpc/mock/server.go -package=mockrpc
	mockgen -source=internal/scan/scan.go -destination=internal/scan/mock/scan.go -package=mock
