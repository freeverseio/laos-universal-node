mocks:
	mockgen -source=internal/platform/scan/scan.go -destination=internal/platform/scan/mock/scan.go -package=mock
	mockgen -source=internal/platform/blockchain/client.go -destination=internal/platform/blockchain/mock/client.go -package=mock
	mockgen -source=internal/platform/storage/storage.go -destination=internal/platform/storage/mock/storage.go -package=mock
	mockgen -source=cmd/server/server.go -destination=cmd/server/mock/server.go -package=mock
	mockgen -source=cmd/server/api/handler.go -destination=cmd/server/api/mock/handler.go -package=mock
	mockgen -source=cmd/server/api/rpcmethods.go -destination=cmd/server/api/mock/rpcmethods.go -package=mock
	mockgen -source=internal/platform/state/state.go -destination=internal/platform/state/mock/state.go -package=mock
	mockgen -source=internal/platform/state/tree/enumerated/tree.go -destination=internal/platform/state/tree/enumerated/mock/tree.go -package=mock
	mockgen -source=internal/platform/state/tree/enumeratedtotal/tree.go -destination=internal/platform/state/tree/enumeratedtotal/mock/tree.go -package=mock
	mockgen -source=internal/platform/state/tree/ownership/tree.go -destination=internal/platform/state/tree/ownership/mock/tree.go -package=mock
	mockgen -source=internal/core/processor/evolution/processor.go -destination=internal/core/processor/evolution/mock/processor.go -package=mock
	mockgen -source=internal/core/processor/evolution/client.go -destination=internal/core/processor/evolution/mock/client.go -package=mock
	mockgen -source=internal/core/processor/universal/discoverer/validator/validator.go -destination=internal/core/processor/universal/discoverer/validator/mock/validator.go -package=mock
	mockgen -source=internal/core/processor/universal/discoverer/discoverer.go -destination=internal/core/processor/universal/discoverer/mock/discoverer.go -package=mock
	mockgen -source=internal/core/processor/universal/updater/updater.go -destination=internal/core/processor/universal/updater/mock/updater.go -package=mock
	mockgen -source=internal/core/processor/universal/processor.go -destination=internal/core/processor/universal/mock/processor.go -package=mock

	
