# netifaces
Portable network interface information - Golang port of Python [netifaces 0.10.5](https://pypi.python.org/pypi/netifaces) 

> 1. What is this?
> 
> It’s been annoying me for some time that there’s no easy way to get the address(es) of the machine’s network interfaces from Python. There is a good reason for this difficulty, which is that it is virtually impossible to do so in a portable manner. However, it seems to me that there should be a package you can easy_install that will take care of working out the details of doing so on the machine you’re using, then you can get on with writing Python code without concerning yourself with the nitty gritty of system-dependent low-level networking APIs.

### Python Source

[Source Codes](pysrc/)

### Task

- [ ] ifaddresses : `HAVE_SOCKET_IOCTLS` section needs tests. Need Go binding.  
- [ ] interfaces : Need Go binding.  
- [x] gateways : IP6 portion `AF_INET6` should be tested thoroughly.  

### Native Tests

- OSX  
  run xcode

- Linux 

  ```sh
  cd xcode/netifaces
  gcc  ../../netifaces.c ./main.c
  ./a.out
  ```
  
### GO test

```sh
go test ./...
```

### API

1. Gateway
  - `func FindAllSystemGateways() (*Gateway, error)`
  - `func (g *Gateway) Release()`
  - `func (g *Gateway) DefaultIP4Gateway() (address string, ifname string, err error)`