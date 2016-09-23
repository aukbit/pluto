### Generate private key  
[reference](https://gist.github.com/denji/12b3a568f092ab951456)
```
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out private.key 2048
# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out private.key
```

### Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
Note: for testing purposes only
```
openssl req -new -x509 -sha256 -key private.key -out server.crt -days 3650
```

### Generating the Certficate Signing Request
[reference](https://digitalelf.net/2016/02/creating-ssl-certificates-in-3-easy-steps/)
#### create a file csr.cnf with the following minimum config
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost:8443

#### And then create the csr
``` 
 openssl req -new -sha256 -key private.key -out certificate.csr \
    -subj "/C=UK/ST=London/O=Pluto/CN=localhost:8443" \
    -config csr.cnf
```
### Signing your certificate with letsencrypt
#### Set up a letsencrypt in a virtualenv
```
pip install virtualenv
virtualenv letsencrypt
cd letsencrypt; source bin/activate
pip install letsencrypt
```
#### run letsencrypt against the certificate
```
letsencrypt -n certonly --agree-tos \
    --email 'nobody@pluto-micro.net' \
    --csr /opt/local/etc/certs/server.csr \
    --cert-path /opt/local/etc/certs/cert.pem \
    --fullchain-path /opt/local/etc/cert/fullchain.pem \
    --webroot -w /opt/www/pluto-micro.net \
    -d digitalelf.net -d www.pluto-micro.net
```