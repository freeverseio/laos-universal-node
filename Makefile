mocks:
	mockgen -source=scan/scan.go -destination=scan/mock/scan.go -package=mock
