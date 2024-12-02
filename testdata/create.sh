mkdir -p private certs

certtool --generate-privkey --key-type ecc --curve secp256r1 --no-text --outfile private/rootca-key.pem
certtool --generate-self-signed --load-privkey private/rootca-key.pem --outfile certs/rootca.pem --template rootca.cfg --hash=sha256

certtool --generate-privkey --key-type ecc --curve secp256r1 --no-text --outfile private/intermca-key.pem
certtool --generate-certificate --load-ca-privkey private/rootca-key.pem --load-ca-certificate certs/rootca.pem \
    --load-privkey private/intermca-key.pem --outfile certs/intermca.pem --template intermca.cfg --hash=sha256

certtool --generate-privkey --key-type rsa --bits 3072 --no-text --outfile private/key.pem
certtool --generate-certificate --load-ca-privkey private/intermca-key.pem --load-ca-certificate certs/intermca.pem \
    --load-privkey private/key.pem --outfile certs/cert.pem --template ts.cfg --hash=sha256

cat certs/intermca.pem certs/rootca.pem > certs/ca.pem
