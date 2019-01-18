package signer

type testOperation struct {
	Operation      string
	HsmResponse    string
	SignerResponse string
	PublicKeyHash  string
	OpType         int
	Level          string
}

// Test Transactions
var (
	testSecp256k1Tx = testOperation{
		// tezos-client transfer 1 from remote-secp256k1 to remote-secp256k1
		OpType:         opTypeTx,
		Operation:      "\"0380270c97773c117d71d95e96f3a5292f7949f571a31cbc80994f5fea61b1546608000154f5d8f71ce18f9f05bb885a4120e64c667bc1b4fb09b9b037d84f00c0843d000154f5d8f71ce18f9f05bb885a4120e64c667bc1b400\"",
		HsmResponse:    "31ccb1d176e80b7caa2164d3c18f5c3ae257e68e44b93851687d2be2b8d0725f8ae4458e8e7174ade426ef57d08970184b4261bd8b65eca100110e246e30b722",
		SignerResponse: "{\"signature\":\"spsig1CKrXpQWRoyKxcJHFXGT3sc9ZpdpBEQwLmjoJQitLQCg8hSxrcoMwuZw4bfaC44K4k4U57QBhneeNy389vNFuS7oNtTCwF\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
	}
	testP256Tx = testOperation{
		// tezos-client transfer 1 from remote-secp256r1 to remote-secp256r1
		OpType:         opTypeTx,
		Operation:      "\"0307456de90f901440e17e76d95a79b74827cc5663ca36994d8603992bea6d6637080002d0ea30de52fb4806d075ab8d312d19be7d0c23e9fb09be8e35d84f00c0843d0002d0ea30de52fb4806d075ab8d312d19be7d0c23e900\"",
		HsmResponse:    "385321c63d21c65009fb0cd8c1845bfb7f2e69048a844040176c4178488f1315c6d8970d6f356c05c1ec13864e21d9a5e0f627e276f50126f38a4bce2de1ffa6",
		SignerResponse: "{\"signature\":\"p2sigUfup3yJF6tQUAzzztLFyAtSwXHiVm6TinFEgB858JAeeopgJ5Ns4iX34i63N7N3hyxVtuXHmUAVj4KqY13renR5L3PAMx\"}",
		PublicKeyHash:  "tz3fNgiRyEZeXD5eh6rEocSp8PBzii2w38Ku",
	}
)

// Test Endorsements
var (
	testEndorse = testOperation{
		// tezos-client endorse for remote-secp256k1
		OpType:         opTypeEndorsement,
		Operation:      "\"027a06a770e6cebe5b3e39483a13ac35f998d650e8b864696e31520922c7242b88c8d2ac55000003eb6d\"",
		HsmResponse:    "f41956681a9a17e4d48ee8e62ccd179f9d12a29155858b5993b013fcb570b10951d25c52ed0b84f0a548a6bf7968e0e77bbc2d190f2a14c2bbfe3a97512c1311",
		SignerResponse: "{\"signature\":\"spsig1dkD3k1tKoyiwno2cLJB9tgTgFJzW9tAXzDn5NbvDaamKggVRSnCRsCfBu8j7K5xoZmEmijstVhit1Z9A4mpggpemq2zBs\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Level:          "256877",
	}

	testEndorseLevel259938 = testOperation{
		// tezos-client endorse for remote-secp256k1
		OpType:         opTypeEndorsement,
		Operation:      "\"027a06a7706ed859dc7394d7216ed6ab088b51089d0dba1707e8544b5c898ef084df727aa4000003f762\"",
		HsmResponse:    "2f63016c1c9638e2630dc0056f3f625903efbcac26d5978aa3752d6050319068f6641148fda3d0a591c9a4913c863b5c90ecb029ee737e28aeed19795d62eeb8",
		SignerResponse: "{\"signature\":\"spsig1C1YcyDsYwiV2F1YimwQUDPuz1AuCj5UVb6rfZ2Dm1iCj7k1aKY31Nxnikx13W3NGjf9BbbWaPpZWJx3qq8MNLp2YX3bvU\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Level:          "259938",
	}

	testEndorseLevel259939 = testOperation{
		// tezos-client endorse for remote-secp256k1
		OpType:         opTypeEndorsement,
		Operation:      "\"027a06a7705aa1570d2cd7d36982adcc4d58c08601f32654131017e19b6b36fdbfd1cbc9da000003f763\"",
		HsmResponse:    "715daf2be170b827df8e71352939f5fda7e920aaa1f21332d3ee2dd9ea46cf1b3b4ee3e86834857acfd0779ad7988c339d76d24016d26603dd0a057f7be285a9",
		SignerResponse: "{\"signature\":\"spsig1LeCXtYt7Ru24o3EyEuHcnxSfDbVDrUtkf9RXwJ23DwXZBrpspdG3S9TP842Bopb6jSEKNViMGSDLGeX6ejrdHyNcsjb1Z\"}",
		PublicKeyHash:  "tz2G4TwEbsdFrJmApAxJ1vdQGmADnBp95n9m",
		Level:          "259939",
	}
)
