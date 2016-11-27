package pcrypto

const TestKeySignature string = "natR922oRMExBFdXFAu9wYT6m0mufe15D5hLLT+K3JIjsc8UM4gn3luBSdqzH0UJJ3ysSztea7eiBlkFM6+PV845iKxACl8LlHg5Fhm4GxIljXCcQQOypMZyqnYG9Iyhggc3lYMAqZHFivM0QuVkK1Ti3SN6341HM+FEcWpR37A="

func TestMasterPublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`)
}

func TestMasterPrivateKey() []byte {
    return []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`)
}

func TestSlavePublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDG0PyjODHigKj6jKnISMAykXgy
2ehOSinHfX5+tRwTAaqe511XnxrPqVB1xygRR3KcjOGNhNa6O7RPJhJ7KA9J7jVV
NaZwAArFNA1CuaHnO+ZXToHC2lcGLcBfTrtgVw1JptZBvH8btWujBHpKD4KTjiDu
u2YYaCeigvLbyUGx+wIDAQAB
-----END PUBLIC KEY-----`)
}

func TestSlavePrivateKey() []byte {
    return []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDG0PyjODHigKj6jKnISMAykXgy2ehOSinHfX5+tRwTAaqe511X
nxrPqVB1xygRR3KcjOGNhNa6O7RPJhJ7KA9J7jVVNaZwAArFNA1CuaHnO+ZXToHC
2lcGLcBfTrtgVw1JptZBvH8btWujBHpKD4KTjiDuu2YYaCeigvLbyUGx+wIDAQAB
AoGAVUOHNVCCREs9LMZqgdSBaK5uSBCfygOQS1eMijaNpbEPRTqgE1XOn8RTF0+j
5VUo1+6rRI/1rsSwHUmMn3icpSxlPw16KtTpAhYfY+KU463HlQ+M2LTwD0JgHFV8
7Jy4STErguyRJ3HoAKImVefLwS8OqOsCaqf6V5qXOtovErECQQDXQQa24Hz2+Va0
06Im2hL7e1vXZVhyPcog2f2nnadWBokxPWiWP6TkSJrV+rslybdhyc41boxZSvq/
CPIYIro/AkEA7HNo8HiZJOedd7l9lEkXwb7O4Oor6Cxu6Fg6zNTfeySuxp5nS0N9
a4K4Ohfw3P3f2HjKtQuNWQy3Bb3Qr4TBRQJADbgeRm+eZ1tS9Gl8rz888HxXSS4z
aeyYQmnCafl5XdlCyzmfvdvGlaou/C5j2S+3GWt0UiF+nn5R5vUaAQHNnwJAOM9Z
zT0MfoNvoA5fD7uoC5LOnddliTjzxLs+FWyn7SxZGbuBUeH7RlN38+1An7gXiikr
eug1o8mcR7Ldau5YiQJAMXdZMl4zTO/FSwvwjAAnh5M+g8qvGVuIKhZ5SB7NT7sA
D1aAZ3gE6zFmOmH4cenuHZ1ha82Np4CEVnRaee91YA==
-----END RSA PRIVATE KEY-----`)
}

var TestAESKey []byte = []byte("longer means more possible keys ")
var TestAESCryptor, _ = NewAESCrypto(TestAESKey)

var TestMasterRSAEncryptor, _ = NewEncryptorFromKeyData(TestSlavePublicKey(), TestMasterPrivateKey())
var TestMasterRSADecryptor, _ = NewDecryptorFromKeyData(TestSlavePublicKey(), TestMasterPrivateKey())

var TestSlaveRSAEncryptor, _ = NewEncryptorFromKeyData(TestMasterPublicKey(), TestSlavePrivateKey())
var TestSlaveRSADecryptor, _ = NewEncryptorFromKeyData(TestMasterPublicKey(), TestSlavePrivateKey())