package rpc

import (
	"log"
)

type System struct{}

type Args struct{}

// SystemResponse holds the result of the Multiply method.
type SystemResponse struct {
	Up int
}

func NewSystemHealthService() *System {
	return &System{}
}

// Up increments the Up field of the SystemResponse struct by 1.
func (a *System) Up(_ *Args) (int, error) {
	log.Println("Up")
	
	return 1, nil
}