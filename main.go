package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/polychain/tezos-remote-signer/signer"
)

var (
	keyfile    = flag.String("keyfile", "./keys.yaml", "Yaml file that identifies keys preloaded in your HSM")
	hsmPin     = flag.String("hsm-pin", "", "User PIN to log into the HSM")
	hsmPinFile = flag.String("hsm-pin-file", "", "Text file containing the user PIN to log into the HSM")
	hsmSO      = flag.String("hsm-so", "", "Shared object used to access the HSM")
	bind       = flag.String("bind", "localhost:6732", "Host:Port for the signer to bind to")
	enableTx   = flag.Bool("enable-tx", false, "WARNING: Allows the signer to sign transactions that move funds.")
	debug      = flag.Bool("debug", false, "Enable debug mode")
)

// getBlockfile to lock the block height.
func getHeightLockDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getPinFromHsmFile(file string) *string {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Error reading %v\n", file)
	}
	pin := strings.Replace(string(contents), "\n", "", 1)
	return &pin
}

func main() {
	flag.Parse()

	if len(*hsmPinFile) > 0 && len(*hsmPin) > 0 {
		log.Fatal("Only one of --hsm-pin and --hsm-pin-file can be set")
	}

	if len(*hsmPinFile) > 0 {
		hsmPin = getPinFromHsmFile(*hsmPinFile)
	}

	signingServer := signer.CreateServer(*keyfile, *hsmPin, *hsmSO, *bind, *enableTx, *debug, getHeightLockDir())
	signingServer.Serve()
}
