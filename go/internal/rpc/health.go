package rpc

type System struct{}

type Args struct{}

type SystemResponse struct {
	Up int
}

func NewSystemHealthService() *System {
	return &System{}
}

func (a *System) Up(_ *Args) (int, error) {
	return 1, nil
}
