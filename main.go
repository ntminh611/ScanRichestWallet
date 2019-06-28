package main

import (
	"./service"
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

type Wallet struct {
	publicKey, privateKey string
}

var tlBot = service.NewTelegramBot("523919774:AAEL5bsetD3yhIMpjgiQSvWFY9dfvs1ibAQ", -1001341233186)

func main() {
	//"https://etherscan.io/accounts"
	//"https://99bitcoins.com/bitcoin-rich-list-top1000/"

	eth = make(map[string]string)
	btc = make(map[string]string)

	readFile("eth")
	fmt.Println(len(eth))
	if len(eth) <= 0 {
		getEthWallets(1)
		writeFile("eth")
	}

	//readFile("btc")
	//fmt.Println(len(btc))
	//if len(btc) <= 0{
	//	getBtcWallets()
	//	writeFile("btc")
	//}

	for i := 0; i < 4; i++ {
		go func() {
			for {
				ethWallet := generateEthWallet()
				if eth[ethWallet.publicKey] == "1" {
					//log.Println(ethWallet)
					fmt.Println(ethWallet)
					//log.Println("public: " + ethWallet.publicKey +
					//	" private: " + ethWallet.privateKey)
					tlBot.SendMessage("public: " + ethWallet.publicKey +
						" private: " + ethWallet.privateKey)
				}
			}
		}()
	}

	for {
		fmt.Println("scan richest wallets Ping!")
		//tlBot.SendMessage("scan richest wallets Ping!")
		time.Sleep(30 * time.Minute)
	}
}

func generateEthWallet() Wallet {
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())
	return Wallet{strings.ToLower(address), privateKey}
}

func getEthWallets(page int) {
	fmt.Println(page)
	res, err := http.Get(fmt.Sprintf("https://etherscan.io/accounts/%d?", page))
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
		getEthWallets(page + 1)
	}
}

func writeFile(name string) {
	var jsonData []byte
	var err error
	if name == "eth" {
		jsonData, err = json.Marshal(eth)
	} else {
		jsonData, err = json.Marshal(btc)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonData))

	jsonFile, err := os.Create(fmt.Sprintf("./%s.json", name))

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())
}

func readFile(name string) {
	jsonFile, err := os.Open(fmt.Sprintf("./%s.json", name))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Successfully Opened %s.json", name))
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if name == "eth" {
		json.Unmarshal([]byte(byteValue), &eth)
	} else {
		json.Unmarshal([]byte(byteValue), &btc)
	}
}
