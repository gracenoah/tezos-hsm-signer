package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"gitlab.com/polychain/tezos-hsm-signer/signer"
	"gitlab.com/polychain/tezos-hsm-signer/signer/watermark"
)

var (
	// Server Flags
	bind    = flag.String("bind", "localhost:6732", "Host:Port for the signer to bind to")
	keyfile = flag.String("keyfile", "./keys.yaml", "Yaml file that identifies keys preloaded in your HSM")
	debug   = flag.Bool("debug", false, "Enable debug mode")
	// Operation Filter Flags
	enableGeneric        = flag.Bool("enable-generic", false, "WARNING: Enables all generic operations including transmitting funds")
	enableTx             = flag.Bool("enable-tx", false, "Enable transferring funds")
	enableVoting         = flag.Bool("enable-voting", false, "Enable voting on proposals")
	txWhitelistAddresses = flag.String("tx-whitelist-addresses", "", "tz... addresses that transfers are enabled to")
	txDailyMax           = flag.String("tx-daily-max", "", "Max amount of XTZ that can be sent offsite in a 24 hour period")
	// HSM Flags
	hsmPin     = flag.String("hsm-pin", "", "User PIN to log into the HSM")
	hsmPinFile = flag.String("hsm-pin-file", "", "Text file containing the user PIN to log into the HSM")
	hsmSO      = flag.String("hsm-so", "", "Shared object used to access the HSM")
	// Watermark Flags
	watermarkType  = flag.String("watermark-type", "file", "Location to store high-watermark.  One of \"ignore\", \"session\", \"file\" or \"dynamodb\"")
	watermarkTable = flag.String("watermark-table", "tezos-hsm-signer", "If --watermark-type is \"dynamodb\", the DynamoDB table to store high-watermarks in")
	watermarkFile  = flag.String("watermark-file", "", "If --watermark-type is \"file\", the file to store high-watermarks in.  Default is ${HOME}/.hsm-signer-watermarks")
)

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

	// Process HSM flags
	if len(*hsmPinFile) > 0 && len(*hsmPin) > 0 {
		log.Fatal("Only one of --hsm-pin and --hsm-pin-file can be set")
	}
	if len(*hsmPinFile) > 0 {
		hsmPin = getPinFromHsmFile(*hsmPinFile)
	}

	// Process Watermark Flags
	var wm watermark.Watermark
	if *watermarkType == "ignore" {
		wm = watermark.GetIgnoreWatermark()
	} else if *watermarkType == "session" {
		wm = watermark.GetSessionWatermark()
	} else if *watermarkType == "file" {
		wm = watermark.GetFileWatermark(*watermarkFile)
	} else if *watermarkType == "dynamodb" {
		wm = watermark.GetDynamoWatermark(*watermarkTable)
	} else {
		panic("Invalid --watermark-type provided")
	}

	// Process Operation Flags
	opFilter := signer.OperationFilter{
		EnableGeneric: *enableGeneric,
		EnableTx:      *enableTx,
		EnableVoting:  *enableVoting,
	}
	if len(*txDailyMax) > 0 {
		opFilter.TxDailyMax, _ = new(big.Int).SetString(*txDailyMax, 10)
		// Convert from XTZ to uXTZ
		opFilter.TxDailyMax.Mul(opFilter.TxDailyMax, new(big.Int).SetInt64(1000000))
	}
	if len(*txWhitelistAddresses) > 0 {
		opFilter.TxWhitelistAddresses = strings.Split(*txWhitelistAddresses, ",")
	}

	signingServer := signer.CreateServer(*keyfile, *hsmPin, *hsmSO, *bind, opFilter, *debug, wm)
	signingServer.Serve()
}
