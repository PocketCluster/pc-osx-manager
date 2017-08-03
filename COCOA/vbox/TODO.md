# TODO

- [ ] pass appropriate host network interface name  
- [x] port forwarding host : (127.0.01:3022) <-> guest : (x.x.x.15:3022)
- [x] sata host cache
- [x] pass a folder to be shared
- [x] pass a name for shared foler
- [x] add network interfaces  
- [x] add `NAT` or `Host Adapter` in the future for communication.
- [x] add **Shared Folder**  
- [x] make sure we put enough time buffer + progress monitor when hard disk to be created  
- [x] wrap in `Object-C` static library
- [x] add test cases  
- [x] check options from `boot2docker`

  * VBoxInternal/CPUM/EnableHVP
  * m.Flag |= F_pae
  * m.Flag |= F_longmode // important: use x86-64 processor
  * m.Flag |= F_rtcuseutc
  * m.Flag |= F_acpi
  * m.Flag |= F_ioapic
  * m.Flag |= F_hpet
  * m.Flag |= F_hwvirtex
  * m.Flag |= F_vtxvpid
  * m.Flag |= F_largepages
  * m.Flag |= F_nestedpaging
