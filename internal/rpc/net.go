package rpc

import "strconv"

// Net is an interface that defines network-related operations.
type Net interface {
	Version() string
}

// netService is an implementation of the Net interface.
type netService struct {
	networkID uint64
}

// NewNetService creates a new instance of netService with the given networkID,
// and returns it as a Net interface.
func NewNetService(networkID uint64) Net {
	return &netService{networkID: networkID}
}

// Version returns the networkID of the netService as a string.
func (b *netService) Version() string {
	return strconv.Itoa(int(b.networkID))
}
