mocks:
	mockgen -source=internal/scan/scan.go -destination=internal/scan/mock/scan.go -package=mock
	mockgen -source=internal/platform/storage/storage.go -destination=internal/platform/storage/mock/storage.go -package=mock
	mockgen -source=cmd/server/server.go -destination=cmd/server/mock/server.go -package=mock
	mockgen -source=cmd/server/api/handler.go -destination=cmd/server/api/mock/handler.go -package=mock
	mockgen -source=internal/state/state.go -destination=internal/state/mock/state.go -package=mock
	mockgen -source=internal/state/enumerated/tree.go -destination=internal/state/enumerated/mock/tree.go -package=mock
	mockgen -source=internal/state/enumeratedtotal/tree.go -destination=internal/state/enumeratedtotal/mock/tree.go -package=mock
	mockgen -source=internal/state/ownership/tree.go -destination=internal/state/ownership/mock/tree.go -package=mock