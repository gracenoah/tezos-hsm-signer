Tezos HSM Signer
================

![pipeline status](https://gitlab.com/polychain/tezos-hsm-signer/badges/master/pipeline.svg) ![coverage](https://gitlab.com/polychain/tezos-hsm-signer/badges/master/coverage.svg)

Implement the Tezos HTTP signing interface, backed by an HSM over PKCS#11.

### Usage

Install and start the signer:

```shell
go get -u gitlab.com/polychain/tezos-hsm-signer

# Identify HSM keys and slots/labels
$ vi keys.yaml

# Launch a server backed by SoftHSM that can send 
# up to 500 XTZ per day to the listed tz address
tezos-hsm-signer \
    --bind "localhost:6732" \
    --hsm-so "/usr/local/lib/softhsm/libsofthsm2.so" \
    --hsm-pin "1234" \
    --enable-tx \
    --tx-daily-max 500 \
    --tx-whitelist-addresses "tz1...,tz2..." \
    --key-file "./keys.yaml"
```

Interact with the signer from tezos-client:

```shell
# Import keys to your client managed by this signer
tezos-client import secret key remote http://localhost:6732/tz...
# Sign an operation with the hsm signer
tezos-client transfer 1 from remote to remote
```

### Development

```shell 
go test ./...
go run main.go
```

**Future Work**

* Improve request parsing
* Validate signatures before returning
* Finish functional testing w/ SoftHSM in Gitlab CI
* Better testing of file and HSM locking
