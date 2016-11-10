package locator

type DebugCommChannel struct {
    LastMcastMessage []byte
    LastUcaseMessage []byte
    LastUcaseHost    string
    MCommCount       uint
    UCommCount       uint
}

func (dc *DebugCommChannel) McastSend(data []byte) error {
    dc.LastMcastMessage = data
    dc.MCommCount++
    return nil
}

func (dc *DebugCommChannel) UcastSend(data []byte, target string) error {
    dc.LastUcaseMessage = data
    dc.LastUcaseHost = target
    dc.UCommCount++
    return nil
}
