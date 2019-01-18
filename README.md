Tezos Remote Signer
===================

![pipeline status](https://gitlab.com/polychain/tezos-remote-signer/badges/master/pipeline.svg) ![coverage](https://gitlab.com/polychain/tezos-remote-signer/badges/master/coverage.svg)

Implement the Tezos HTTP signing interface, backed by an HSM over PKCS#11.

### Usage

```shell 
# Identify the keys and their HSM slots or labels
$ vi keys.yaml
# Start the signer
go run main.go  \
    --bind "localhost:6732" \
    --hsm-so "/usr/local/lib/softhsm/libsofthsm2.so" \
    --hsm-pin "1234" \
    --key-file "./keys.yaml"

# Import keys to your client managed by this signer
tezos-client import secret key remote http://localhost:6732/tz...
# Sign an operation with the remote signer
tezos-client transfer 1 from remote to remote
```

### Testing

```shell
# Unit Testing
go test ./...
```

**Future Work**

* Improve request parsing
* Validate signatures before returning
* Finish functional testing w/ SoftHSM in Gitlab CI
* Better testing of file and HSM locking
