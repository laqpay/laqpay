package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/laqpay/laqpay/src/cipher"
)

// Note: Address_gen generates public keys and addresses
// address, pubkey, privatekey
// -n=5 for number of addresses
// -seed to set wallet seed. Prompt will ask
// for seed to prevent seed from being stored in .bashrc

// -json for json output
// -add option to password the secret key
// -let people add the key from the command line

var (
	// HideSeckey whether need hide secret key
	HideSeckey = false
	seed       = ""
	genCount   = 1
)

func registerFlags() {

	flag.IntVar(&genCount, "n", genCount,
		"number of addresses to generate")

	flag.BoolVar(&HideSeckey, "s", HideSeckey,
		"only generate publickey and address, hide seckey")

	flag.StringVar(&seed, "seed", seed,
		"seed for deterministic key generation")

	//flag.StringVar(&outFile, "o", outFile,
	//    "If present, will create a new wallet entry and write to disk. "+
	//        "For safety, it will not overwrite an existing keypair")
	//flag.BoolVar(&printSecret, "print-secret", printSecret,
	//    "Print the wallet entry's secret key")
	//flag.StringVar(&inFile, "i", inFile,
	//    "Will read a wallet entry from this file for printing info")
}

func parseFlags() {
	flag.Parse()
}

// Wallet represents the wallet
type Wallet struct {
	Meta    map[string]string `json:"meta"`
	Entries []KeyEntry        `json:"entries"`
}

// KeyEntry represents the key entry in wallet
type KeyEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

func getKeyEntry(pub cipher.PubKey, sec cipher.SecKey) KeyEntry {

	var e KeyEntry

	//laqpay address
	e = KeyEntry{
		Address: cipher.AddressFromPubKey(pub).String(),
		Public:  pub.Hex(),
		Secret:  sec.Hex(),
	}

	//hide the secret key
	if HideSeckey == true {
		e.Secret = ""
	}

	return e
}

func main() {
	registerFlags()
	parseFlags()

	w := Wallet{
		Meta:    make(map[string]string), //map[string]string
		Entries: make([]KeyEntry, genCount),
	}

	w.Meta = map[string]string{"coin": "laqpay"}

	if seed == "" { //generate a new seed, as hex string
		seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	}

	w.Meta["seed"] = seed

	seckeys, _ := cipher.GenerateDeterministicKeyPairs([]byte(seed), genCount)

	for i, sec := range seckeys {
		pub, _ := cipher.PubKeyFromSecKey(sec)
		w.Entries[i] = getKeyEntry(pub, sec)
	}

	output, err := json.MarshalIndent(w, "", "    ")
	if err != nil {
		fmt.Printf("Error formating wallet to JSON. Error : %s\n", err.Error())
		return
	}
	fmt.Printf("%s\n", string(output))

}

/*
   if outFile != "" {
       w := createWalletEntry(outFile, testNetwork)
       if w != nil {
           printWalletEntry(w, labelStdout, PrintAddress, printPublic,
               printSecret)
       }
   }
   if inFile != "" {
       printWalletEntryFromFile(inFile, labelStdout, PrintAddress,
           printPublic, printSecret)
   }
*/
