package context

func MasterIPAddress() (string, error) {
    return "192.168.1.236", nil
}

func MasterLiveInterface() ([]string, error) {
    return []string{"en0"}, nil
}
