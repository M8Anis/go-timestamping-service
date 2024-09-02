mkdir -p test-root

certtool --generate-privkey --key-type ecc --curve secp256r1 --no-text --outfile test-root/ca-key.pem
certtool --generate-self-signed --load-privkey test-root/ca-key.pem --outfile test-root/ca-crt.pem --template test-root/ca.cfg --hash=sha256

mkdir -p test-ts-ca

certtool --generate-privkey --key-type ecc --curve secp256r1 --no-text --outfile test-ts-ca/ca-key.pem
certtool --generate-certificate --load-ca-privkey test-root/ca-key.pem --load-ca-certificate test-root/ca-crt.pem \
    --load-privkey test-ts-ca/ca-key.pem --outfile test-ts-ca/ca-crt.pem --template test-ts-ca/ca.cfg --hash=sha256

certtool --generate-privkey --key-type rsa --bits 3072 --no-text --outfile ts-key.pem
certtool --generate-certificate --load-ca-privkey test-ts-ca/ca-key.pem --load-ca-certificate test-ts-ca/ca-crt.pem \
    --load-privkey ts-key.pem --outfile ts-crt.pem --template ts.cfg --hash=sha256

cat ts-crt.pem test-ts-ca/ca-crt.pem test-root/ca-crt.pem > ts-crt_chain.pem
