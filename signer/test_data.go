package signer

type testOperation struct {
	Operation      string
	HsmResponse    string
	SignerResponse string
	PublicKeyHash  string
	OpWatermark    uint8
	Level          string
	ChainID        string
}

// Test Transactions
var (
	testSecp256k1Tx = testOperation{
		// tezos-client transfer 1 from remote-secp256k1 to remote-secp256k1
		OpWatermark:    opWatermarkGeneric,
		Operation:      "\"0380270c97773c117d71d95e96f3a5292f7949f571a31cbc80994f5fea61b1546608000154f5d8f71ce18f9f05bb885a4120e64c667bc1b4fb09b9b037d84f00c0843d000154f5d8f71ce18f9f05bb885a4120e64c667bc1b400\"",
		HsmResponse:    "31ccb1d176e80b7caa2164d3c18f5c3ae257e68e44b93851687d2be2b8d0725f8ae4458e8e7174ade426ef57d08970184b4261bd8b65eca100110e246e30b722",
		SignerResponse: "{\"signature\":\"spsig1CKrXpQWRoyKxcJHFXGT3sc9ZpdpBEQwLmjoJQitLQCg8hSxrcoMwuZw4bfaC44K4k4U57QBhneeNy389vNFuS7oNtTCwF\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		ChainID:        "NetXeSG6ShTTieu",
	}
	testP256Tx = testOperation{
		// tezos-client transfer 1 from remote-secp256r1 to remote-secp256r1
		OpWatermark:    opWatermarkGeneric,
		Operation:      "\"0307456de90f901440e17e76d95a79b74827cc5663ca36994d8603992bea6d6637080002d0ea30de52fb4806d075ab8d312d19be7d0c23e9fb09be8e35d84f00c0843d0002d0ea30de52fb4806d075ab8d312d19be7d0c23e900\"",
		HsmResponse:    "385321c63d21c65009fb0cd8c1845bfb7f2e69048a844040176c4178488f1315c6d8970d6f356c05c1ec13864e21d9a5e0f627e276f50126f38a4bce2de1ffa6",
		SignerResponse: "{\"signature\":\"p2sigUfup3yJF6tQUAzzztLFyAtSwXHiVm6TinFEgB858JAeeopgJ5Ns4iX34i63N7N3hyxVtuXHmUAVj4KqY13renR5L3PAMx\"}",
		PublicKeyHash:  "tz3fNgiRyEZeXD5eh6rEocSp8PBzii2w38Ku",
		ChainID:        "NetXJDZUe2asiD2",
	}
)

// Test Endorsements
var (
	testEndorse = testOperation{
		// tezos-client endorse for remote-secp256k1
		OpWatermark:    opWatermarkEndorsement,
		Operation:      "\"027a06a770e6cebe5b3e39483a13ac35f998d650e8b864696e31520922c7242b88c8d2ac55000003eb6d\"",
		HsmResponse:    "f41956681a9a17e4d48ee8e62ccd179f9d12a29155858b5993b013fcb570b10951d25c52ed0b84f0a548a6bf7968e0e77bbc2d190f2a14c2bbfe3a97512c1311",
		SignerResponse: "{\"signature\":\"spsig1dkD3k1tKoyiwno2cLJB9tgTgFJzW9tAXzDn5NbvDaamKggVRSnCRsCfBu8j7K5xoZmEmijstVhit1Z9A4mpggpemq2zBs\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Level:          "256877",
		ChainID:        "NetXdQprcVkpaWU",
	}

	testEndorseLevel259938 = testOperation{
		// tezos-client endorse for remote-secp256k1
		OpWatermark:    opWatermarkEndorsement,
		Operation:      "\"027a06a7706ed859dc7394d7216ed6ab088b51089d0dba1707e8544b5c898ef084df727aa4000003f762\"",
		HsmResponse:    "2f63016c1c9638e2630dc0056f3f625903efbcac26d5978aa3752d6050319068f6641148fda3d0a591c9a4913c863b5c90ecb029ee737e28aeed19795d62eeb8",
		SignerResponse: "{\"signature\":\"spsig1C1YcyDsYwiV2F1YimwQUDPuz1AuCj5UVb6rfZ2Dm1iCj7k1aKY31Nxnikx13W3NGjf9BbbWaPpZWJx3qq8MNLp2YX3bvU\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Level:          "259938",
		ChainID:        "NetXdQprcVkpaWU",
	}

	testEndorseLevel259939 = testOperation{
		// tezos-client endorse for remote-secp256k1
		OpWatermark:    opWatermarkEndorsement,
		Operation:      "\"027a06a7705aa1570d2cd7d36982adcc4d58c08601f32654131017e19b6b36fdbfd1cbc9da000003f763\"",
		HsmResponse:    "715daf2be170b827df8e71352939f5fda7e920aaa1f21332d3ee2dd9ea46cf1b3b4ee3e86834857acfd0779ad7988c339d76d24016d26603dd0a057f7be285a9",
		SignerResponse: "{\"signature\":\"spsig1LeCXtYt7Ru24o3EyEuHcnxSfDbVDrUtkf9RXwJ23DwXZBrpspdG3S9TP842Bopb6jSEKNViMGSDLGeX6ejrdHyNcsjb1Z\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Level:          "259939",
		ChainID:        "NetXdQprcVkpaWU",
	}
)

// Test Blocks
var (
	testBlock = testOperation{
		// tezos-client endorse for remote-secp256k1
		OpWatermark:    opWatermarkBlock,
		Operation:      "\"018eceda2f00023df201a43530a7ac35eb2b0000c51c29e21b943355242b73d76cfee53dd80908e3e843000000005c562260047d63043ed1c565061d83ef6ce2fed35410549d87350808ee050c25c1a2dc7795000000110000000100000000080000000000435a1af44a000a5e6a5297b6531706fed4438e7e1c134a9e70d239f4c17c3d81411258000000000003dcf012f600\"",
		HsmResponse:    "428fa4f31d7e6c4ec1a100618abd4ac0c8f100d67fb754226c185c0bf93f60562c60592f5189a0797b23d519d67babd2ad379055a1f639fdad8af1daaf0ba333",
		SignerResponse: "{\"signature\":\"spsig1EX3PsUAHsQQUYpztfrV5w1GEPsDwJmLBhE2JSUCinH9hBgbL2fwbG73ZYfSB4pJ6aW98gTGh1VMBBU7YcGPQiBmX2o7kM\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Level:          "146930",
		ChainID:        "NetXgtSLGNJvNye",
	}
)
