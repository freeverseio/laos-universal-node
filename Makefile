mocks:
	mockgen -source=internal/scan/scan.go -destination=internal/scan/mock/scan.go -package=mock
	mockgen -source=cmd/server/server.go -destination=cmd/server/mock/server.go -package=mock
	mockgen -source=cmd/server/api/handler.go -destination=cmd/server/api/mock/handler.go -package=mock