openssl genpkey -algorithm RSA -out ca.key -pkeyopt rsa_keygen_bits:2048

openssl req -new -x509 -key ca.key -out ca.crt --days 10950