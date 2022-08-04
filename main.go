package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"strings"
)

type Wallet struct {
	publicKey, privateKey string
}

//var tlBot = service.NewTelegramBot("5378904643:AAHR48MZY1EwNeZlHhW076fsrsccyxzROtk", -1001341233186)

func main() {
	wallet := make(map[string]string)
	for {
		ethWallet := generateEthWallet()
		public := strings.ToLower(ethWallet.publicKey)
		if isNice(public) {
			wallet[public] = ethWallet.privateKey
			writeWallets(wallet, "wallet")
		}
	}
}

const Nice = 6

func isNice(s string) bool {
	length := len(s)
	count := 0
	pin := s[length-1]
	for i := length - 2; i > 0; i-- {
		if s[i] == pin {
			count++
		} else {
			break
		}
	}
	if count >= Nice {
		return true
	}

	count = 0
	for i := length - 2; i > 0; i-- {
		if s[i] == pin-1 {
			count++
			pin = s[i]
		} else {
			break
		}
	}
	if count >= Nice {
		return true
	}
	return false
}

func generateEthWallet() Wallet {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())
	return Wallet{strings.ToLower(address), privateKey}
}

func writeWallets(data map[string]string, name string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	jsonFile, err := os.Create(fmt.Sprintf("./%s.json", name))

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}
