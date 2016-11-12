package locator

type DebugCommChannel struct {
    LastMcastMessage []byte
    LastUcastMessage []byte
    LastUcastHost    string
    MCommCount       uint
    UCommCount       uint
}

func (dc *DebugCommChannel) McastSend(data []byte) error {
    dc.LastMcastMessage = data
    dc.MCommCount++
    return nil
}

func (dc *DebugCommChannel) UcastSend(data []byte, target string) error {
    dc.LastUcastMessage = data
    dc.LastUcastHost = target
    dc.UCommCount++
    return nil
}
