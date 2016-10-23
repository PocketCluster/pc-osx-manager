**This is where interface maintaiing policy resides**

## Selection Logic

1. find primary interface & address from wifi and ethernet.  
2. Assume that the client connection would come though this addresses. 
3. When a client request comes through, find an interface and address that matches with the client.      
    - look for gateway & netmask and if master and client match.  
    - `TODO` : check if two addresses falls into the same network segment.  
4. The one with most client attached is the one being primary.
5. The primary interface & address will be remembered until new, different change comes along. 
6. Whenever new client checks in, tell it to connect to the primary interface & address.
