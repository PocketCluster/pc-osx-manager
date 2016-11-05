package service

import (
    "time"

    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
)

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
    locator.NewSlaveLocator(onSucess, onIdle, onFail)
    return nil
}

// On sucess happens at the moment successful state transition takes place
func onSucess(state locator.SlaveLocatingState) {
    switch state {
        case locator.SlaveInquired: {

            break
        }
        case locator.SlaveKeyExchange: {

            break
        }
        case locator.SlaveCryptoCheck: {

            break
        }
        case locator.SlaveBounded: {

            break
        }
        default:
    }
}

// On Idle happens as locator awaits
func onIdle(state locator.SlaveLocatingState, lastIdleAction time.Time, slaveTimestamp time.Time, trialCount int) bool {
    switch state {
        case locator.SlaveUnbounded: {
            if (time.Second * 3) < slaveTimestamp.Sub(lastIdleAction) {
                ua, err := slagent.UnboundedMasterDiscovery()
                if err != nil {
                    return false
                }
                _, err = slagent.UnboundedMasterDiscoveryMeta(ua)
                if err != nil {
                    return false
                }

                // TODO : broadcast slave meta
                return true
            }
        }
        case locator.SlaveInquired: {
            if (time.Second * 3) < slaveTimestamp.Sub(lastIdleAction) {
                agent, err := slagent.AnswerMasterInquiryStatus(slaveTimestamp)
                if err != nil {
                    return false
                }
                _, err = slagent.AnswerMasterInquiryMeta(agent)
                if err != nil {
                    return false
                }

                // TODO : send answer to master
                return true
            }
        }
        case locator.SlaveKeyExchange: {
            if (time.Second * 3) < slaveTimestamp.Sub(lastIdleAction) {
                slctx := slcontext.SharedSlaveContext()

                masterAgentName, err := slctx.GetMasterAgent()
                if err != nil {
                    return false
                }
                agent, err := slagent.KeyExchangeStatus(masterAgentName, slaveTimestamp)
                if err != nil {
                    return false
                }
                _, err = slagent.KeyExchangeMeta(agent, slctx.GetPublicKey())
                if err != nil {
                    return false
                }

                // TODO : send answer to master
                return true
            }
        }
        case locator.SlaveCryptoCheck: {
            if (time.Second * 3) < slaveTimestamp.Sub(lastIdleAction) {
                slctx := slcontext.SharedSlaveContext()

                masterAgentName, err := slctx.GetMasterAgent()
                if err != nil {
                    return false
                }
                slaveAgentName, err := slctx.GetMasterAgent()
                if err != nil {
                    return false
                }
                aesCryptor, err := slctx.AESCryptor()
                if err != nil {
                    return false
                }
                sa, err := slagent.CheckSlaveCryptoStatus(masterAgentName, slaveAgentName, slaveTimestamp)
                if err != nil {
                    return false
                }
                _, err = slagent.CheckSlaveCryptoMeta(sa, aesCryptor)
                if err != nil {
                    return false
                }

                // TODO : send answer to master
                return true
            }
        }
        case locator.SlaveBounded: {
            if (time.Second * 10) < slaveTimestamp.Sub(lastIdleAction) {
                slctx := slcontext.SharedSlaveContext()

                masterAgentName, err := slctx.GetMasterAgent()
                if err != nil {
                    return false
                }
                slaveAgentName, err := slctx.GetMasterAgent()
                if err != nil {
                    return false
                }
                aesCryptor, err := slctx.AESCryptor()
                if err != nil {
                    return false
                }
                sa, err := slagent.SlaveBoundedStatus(masterAgentName, slaveAgentName, slaveTimestamp)
                if err != nil {
                    return false
                }
                _, err = slagent.SlaveBoundedMeta(sa, aesCryptor)
                if err != nil {
                    return false
                }

                // TODO : send answer to master
                return true
            }
        }
        case locator.SlaveBindBroken: {
            if (time.Second * 3) < slaveTimestamp.Sub(lastIdleAction) {
                slctx := slcontext.SharedSlaveContext()

                masterAgentName, err := slctx.GetMasterAgent()
                if err != nil {
                    return false
                }
                ba, err := slagent.BrokenBindDiscovery(masterAgentName)
                if err != nil {
                    return false
                }
                _, err = slagent.BrokenBindMeta(ba)
                if err != nil {
                    return false
                }

                // TODO : send answer to master
                return true
            }
        }
    default:
    }

    return false
}

// OnFail happens at the moment state transition fails to happen
func onFail(state locator.SlaveLocatingState) {
    switch state {
        case locator.SlaveUnbounded: {

            break
        }
        case locator.SlaveInquired: {

            break
        }
        case locator.SlaveKeyExchange: {

            break
        }
        case locator.SlaveCryptoCheck: {

            break
        }
        case locator.SlaveBounded: {

            break
        }
        case locator.SlaveBindBroken: {

            break
        }
    default:
    }
}




