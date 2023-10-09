mocks:
	mockgen -source=scanner/scanner.go -destination=scanner/mock/scanner.go -package=mock
