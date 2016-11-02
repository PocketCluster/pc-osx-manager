package service

type SlaveAgentService interface {
    // this is main loop of locating service to master
    MonitorLocatingService() error
}

func NewSlaveLocatingService() SlaveAgentService {
    return &slaveAgent {
    }
}

type slaveAgent struct {
}

func (sa *slaveAgent) MonitorLocatingService() error {
    return nil
}
