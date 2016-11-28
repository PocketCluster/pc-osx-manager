## How to sign keys with certificate

### Create CA certificate

1. Create a private key called `ca-priv-key.pem` for the CA:
	
  > openssl genrsa -out ca-priv-key.pem 2048

2. Create a public key called `ca.pem` for the CA:

  > openssl req -new -key ca-priv-key.pem -x509 -days 1825 -out ca.pem

### Create Swarm Private/Certificate

1. Create a private key `swarm-priv-key.pem` for your Swarm Manager:
	
  > openssl genrsa -out swarm-priv-key.pem 2048

2. Generate a **certificate signing request** (CSR) `swarm.csr` using the private key you create in the previous step:

  > openssl req -subj "/CN=swarm" -new -key swarm-priv-key.pem -out swarm.csr

3. Create the certificate `swarm-cert.pem` based on the CSR created in the previous step.

  > openssl x509 -req -days 1825 -in swarm.csr -CA ca.pem -CAkey ca-priv-key.pem -CAcreateserial -out swarm-cert.pem -extensions v3_req 

### Create Node Private/Certificate

1. Create a private key `node-priv-key.pem` for your Swarm Manager:

  > openssl genrsa -out node-priv-key.pem 2048

2. Generate a **certificate signing request** (CSR) `node.csr` using the private key you create in the previous step:

  > openssl req -subj "/CN=pc-node1" -new -key node-priv-key.pem -out node.csr

3. Create the certificate `node.cert` based on the CSR created in the previous step.

  > openssl x509 -req -days 1825 -in node.csr -CA ca.pem -CAkey ca-priv-key.pem -CAcreateserial -out node.cert -extensions v3_req 
