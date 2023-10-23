mocks:
	mockgen -source=internal/scan/scan.go -destination=internal/scan/mock/scan.go -package=mock
	mockgen -source=internal/scan/store.go -destination=internal/scan/mock/store.go -package=mock
