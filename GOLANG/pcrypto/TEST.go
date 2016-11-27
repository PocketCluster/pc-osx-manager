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

func TestSlaveNodePublicKey() []byte {
    return []byte(`MacPro:pki almightykim$ cat node.pub
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyO2QqsbF+KYro90d79PB
54cezOXNMDn5WlVov7X0HIUAI9i74HcU9zWUgx1tWimueYJ568XxzuB2Gr6sFWJy
0ZDNXjmLmGvx3hCXd7cfJlm/QvBAOAZy6HV22R2ujnqzEdoAjzzrayLcSySDslkY
00uvDVkky+syL+gxywSXAzYijA2rnQFxFEP9bkdzdNgiPAAbQ7Xi5pM6eQnhbKsU
i5apRgy3+rRJI0NjIpn86pcw7SIM46vscPV18RB8BakPkmVjUa4vDk/w/+/xrezp
Fo0n4myExIywcbqaFoVMzWh8j0STm4cQOjspRLxBkSWdvVHkbK+pjeAQQQZtrFo9
BQIDAQAB
-----END PUBLIC KEY-----`)
}

func TestSlaveNodePrivateKey() []byte {
    return []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAyO2QqsbF+KYro90d79PB54cezOXNMDn5WlVov7X0HIUAI9i7
4HcU9zWUgx1tWimueYJ568XxzuB2Gr6sFWJy0ZDNXjmLmGvx3hCXd7cfJlm/QvBA
OAZy6HV22R2ujnqzEdoAjzzrayLcSySDslkY00uvDVkky+syL+gxywSXAzYijA2r
nQFxFEP9bkdzdNgiPAAbQ7Xi5pM6eQnhbKsUi5apRgy3+rRJI0NjIpn86pcw7SIM
46vscPV18RB8BakPkmVjUa4vDk/w/+/xrezpFo0n4myExIywcbqaFoVMzWh8j0ST
m4cQOjspRLxBkSWdvVHkbK+pjeAQQQZtrFo9BQIDAQABAoIBAQCmyVS5/fAdu8oj
5otQJeYii24MxYDy1FzhGF6wLJirB9ga6XDjHdZAcrCJueao3kqfQKh2B0T25ioD
f10XDzaiMOHYokn3RztpizpAPLjVu8/g/88+8lN2FPOHvHTGfqGgYubt/7KnpzNY
CMJtTDooQv3XRbUetGhfjg2vKWs1VMSGHGi6Pq3Z8wTlY38NT0lBJNO49JNPsurA
DLJiC7Rb6SM/YwWh+XEp+Hm0wX52JTXD31S7LUlzOCnBMq9XAM2dC/TNJXEErG8K
8+F40NUpf8AEYpSqu3aBiyfLTRK+mfdFTtFFkfWBJxJy6nXM08mKNKJvI8V12EPA
W7/hgebhAoGBAO7bfHFljsFAXrYhWbUfP/tXmkZndfUXYZV+FNAt3OxUc455dJpG
uOcDrGWY/dnAxklHPoqQWs7g2rHS9uhozGz/64x4QeVO4RDMhkm2FxnuDeqULynb
xcUdxKxOE4SoPu6IdTlW4xRxdqVDD74KYogU8eQm3iX0eY38M0jyivOJAoGBANdZ
NFrF3fNF3t/piSVHfU6yTKuT0khbY06r5KMsW7EAJJJ4RkSdhCIQ8AooSwlhyc8X
v+vHquYsXfXT8aQEw+ExYh44vN1ZZ5VIYdGofnl1mhab9mVrqdYYa3kr1Ny00Ui5
Gk3FvP/LBhG7Jgc2SRSVtzhQq0RpXqIdX60c2VKdAoGBAKLWgX0xVmRLRQZ3sBe5
qT3p2CRdTl57xSxMW1YdnjqDzI/6H1M6Gb5sk7Bj39P/B29XobyHc1EMnCuU/n0t
TQiWZHhMV+hDoU55kKdZ+1/TGiutQIYR7T9X7wfk5ouOw/CMmRYxNPhv7gn2sRnH
LKtHVC1Nji9j/yacJD58E9y5AoGBAKnQzlhGcB/GmVo47s1W8pl8QLmMd+ZXKph/
NGz4LdYGJtDZx4+UJv42HRPlckaTtnB4af+kFEAt/Go+F+8fUtfh+V2boFNsjSJL
Udfi5tkgw8HQexy/Kc6KszV6OwFQFTkjvnpV1BRiJQcWbYaCaF6zMShXdLcd4GI2
h5wbg8SBAoGAE3vx0D5KqMKyrUzNlvNKORqYLxCl0yACl+hl5GwMXqr774VCQZ2A
x/6Caylpij1+rFP2SNU4tP++2icJUXUdm7fxV7N3mUrVxcfQnpkZii16yfud+mZN
BiMNQ1We6oTEqGxIKwSq4VjojwUOqaby5SLRs552a2wBWf2kwAv/3Ec=
-----END RSA PRIVATE KEY-----`)
}

var TestAESKey []byte = []byte("longer means more possible keys ")
var TestAESCryptor, _ = NewAESCrypto(TestAESKey)

var TestMasterRSAEncryptor, _ = NewRsaEncryptorFromKeyData(TestSlavePublicKey(), TestMasterPrivateKey())
var TestMasterRSADecryptor, _ = NewRsaDecryptorFromKeyData(TestSlavePublicKey(), TestMasterPrivateKey())

var TestSlaveRSAEncryptor, _ = NewRsaEncryptorFromKeyData(TestMasterPublicKey(), TestSlavePrivateKey())
var TestSlaveRSADecryptor, _ = NewRsaEncryptorFromKeyData(TestMasterPublicKey(), TestSlavePrivateKey())