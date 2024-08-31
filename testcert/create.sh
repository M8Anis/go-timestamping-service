#!/bin/bash

certtool --generate-privkey --key-type ecc --curve secp256r1 --no-text --outfile ts-key.pem
certtool --key-info --infile ts-key.pem --outder --outfile ts-key.der

certtool --generate-self-signed --load-privkey ts-key.pem --outfile ts-crt.pem --template ts.cfg --hash=sha256
certtool --certificate-info --infile ts-crt.pem --outder --outfile ts-crt.der
