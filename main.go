package main

import (
	"net/http"

	"fmt"
	crand "crypto/rand"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"log"
)

//错误处理
func handle(why string, e error) {
	if e != nil {
		fmt.Println(why, "错误为：", e)
	}
}

func main() {

	//如果有余额就将私钥和地址存入文件中
	file, e := os.OpenFile("./addr_amount.txt", os.O_WRONLY|os.O_CREATE, 0761)
	handle("文件打开失败！", e)

	defer file.Close()

	//连接以太坊浏览器查询生成地址余额是否大于0
	s := "https://etherscan.io/address/"

	for /* i := 0; i < 1000000; i++ */ {

		privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), crand.Reader)
		//fmt.Printf("%x",privateKeyECDSA.D.Bytes())
		if err != nil {
			return
		}
		address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
		addr := address.String()

		fmt.Println("priv :", privateKeyECDSA.D.String(), "\t", "addr:", addr)
		url := s + addr
		resp, e := http.Get(url)
		if e != nil {
			return
		}
		defer resp.Body.Close()

		go func() {
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			doc.Find(".col-md-8").First().Each(func(i int, selection *goquery.Selection) {
				text := selection.Text()
				b := []byte(text)
				bytes := b[:len(b)-6]
				u, i2 := strconv.ParseFloat(string(bytes), 32)
				if i2 != nil {
					return
				}
				if u > 0.0 {
					fmt.Println("pri: ", privateKeyECDSA.D.String(), "addr :", addr, "balance:", u)
					file.WriteString("pri: " + privateKeyECDSA.D.String() + "addr :" + addr + "\n")

				} else {
					doc.Find(".col-md-8").Last().Each(func(i int, selection *goquery.Selection) {

						text := selection.Text()
						b := []byte(text)
						bytes := b[:len(b)-6]
						u, i2 := strconv.ParseUint(string(bytes), 10, 64)
						if i2 != nil {
							fmt.Println(i2)
							return
						}
						if u > 0 {
							fmt.Println("pri: ", privateKeyECDSA.D.String(), "addr :", addr, "balance:", u, "txs:", selection.Text())
							file.WriteString("pri: " + privateKeyECDSA.D.String() + "addr :" + addr + "\n")
						}
					})
				}
			})
		}()

	}
}
