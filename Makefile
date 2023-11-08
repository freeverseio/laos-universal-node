mocks:
	mockgen -source=internal/scan/scan.go -destination=internal/scan/mock/scan.go -package=mock
	mockgen -source=internal/scan/store.go -destination=internal/scan/mock/store.go -package=mock
	mockgen -source=internal/platform/storage/storage.go -destination=internal/platform/storage/mock/storage.go -package=mock
	mockgen -source=cmd/server/server.go -destination=cmd/server/mock/server.go -package=mock
	mockgen -source=cmd/server/api/handler.go -destination=cmd/server/api/mock/handler.go -package=mock
