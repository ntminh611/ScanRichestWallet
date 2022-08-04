package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var eth map[string]string
var btc map[string]string
var wallet map[string]string
var result map[string]string

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
			writeFile("wallet")
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

	//count = 0
	//for i := length - 2; i > 0; i-- {
	//	if s[i] == pin-1 {
	//		count++
	//		pin = s[i]
	//	} else {
	//		break
	//	}
	//}
	//if count >= Nice {
	//	return true
	//}
	return false
}

func generateEthWallet() Wallet {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())
	return Wallet{strings.ToLower(address), privateKey}
}

func getEthWallets(page int) {
	fmt.Println(page)
	res, err := http.Get(fmt.Sprintf("https://etherscan.io/accounts/%d?ps=100", page))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	doc.Find(`.table tbody tr`).Each(func(i int, selection *goquery.Selection) {
		walletAddr := selection.Find(`td`).Eq(1).Text()
		eth[strings.ToLower(walletAddr)] = "1"
	})

	_, bool := doc.Find(`.pagination [title='Go to Next'] a`).First().Attr("href")
	if bool {
		time.Sleep(time.Second)
		getEthWallets(page + 1)
	}
}

func writeFile(name string) {
	var jsonData []byte
	var err error
	switch name {
	case "eth":
		jsonData, err = json.Marshal(eth)
	case "btc":
		jsonData, err = json.Marshal(btc)
	case "wallet":
		jsonData, err = json.Marshal(wallet)
	case "result":
		jsonData, err = json.Marshal(result)
	}

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

func readFile(name string) {
	jsonFile, err := os.Open(fmt.Sprintf("./%s.json", name))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Successfully Opened %s.json", name))
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	switch name {
	case "eth":
		json.Unmarshal(byteValue, &eth)
	case "btc":
		json.Unmarshal(byteValue, &btc)
	case "wallet":
		json.Unmarshal(byteValue, &wallet)
	case "result":
		json.Unmarshal(byteValue, &result)
	}
}
