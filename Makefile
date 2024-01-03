mocks:
	mockgen -source=internal/platform/scan/scan.go -destination=internal/platform/scan/mock/scan.go -package=mock
	mockgen -source=internal/platform/storage/storage.go -destination=internal/platform/storage/mock/storage.go -package=mock
	mockgen -source=cmd/server/server.go -destination=cmd/server/mock/server.go -package=mock
	mockgen -source=cmd/server/api/handler.go -destination=cmd/server/api/mock/handler.go -package=mock
	mockgen -source=cmd/server/api/rpcmethods.go -destination=cmd/server/api/mock/rpcmethods.go -package=mock
	mockgen -source=internal/platform/state/state.go -destination=internal/platform/state/mock/state.go -package=mock
	mockgen -source=internal/platform/state/enumerated/tree.go -destination=internal/platform/state/enumerated/mock/tree.go -package=mock
	mockgen -source=internal/platform/state/enumeratedtotal/tree.go -destination=internal/platform/state/enumeratedtotal/mock/tree.go -package=mock
	mockgen -source=internal/platform/state/ownership/tree.go -destination=internal/platform/state/ownership/mock/tree.go -package=mock
	mockgen -source=internal/platform/core/processor/evolution/processor.go -destination=internal/platform/core/processor/evolution/mock/processor.go -package=mock
