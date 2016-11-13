package beacon

type DebugCommChannel struct {
    LastUcastMessage []byte
    LastUcastHost    string
    UCommCount       uint
}

func (dc *DebugCommChannel) UcastSend(data []byte, target string) error {
    dc.LastUcastMessage = data
    dc.LastUcastHost = target
    dc.UCommCount++
    return nil
}
