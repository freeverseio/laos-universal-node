mocks:
	mockgen -source=internal/scan/scan.go -destination=internal/scan/mock/scan.go -package=mock
