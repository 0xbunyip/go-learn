package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

func main() {
	urls := []string{
		"https://www.myconstant.com/saving-api/saving-p2p/top-matches/",
		"https://ipfind.co/me?auth=4517756b-4da1-4a37-92e8-3b4eb9f78bfa",
		"https://www.myconstant.com/api/collateral-loan/collaterals",
		"https://www.myconstant.com/exchange-api/crypto/fiat-exchange-rate/?currency=USD",
		// "https://grafana-live.constant.money/api/dashboards/uid/HGVSSGMZk",
		"https://www.myconstant.com/api/system/countries",
		"https://www.myconstant.com/api/system/languages",
	}
	reqs := []*http.Request{}
	for _, url := range urls {
		payload := strings.NewReader("")
		req, _ := http.NewRequest("GET", url, payload)
		reqs = append(reqs, req)
	}

	r := 1000
	in := make(chan int, r)
	for k := 0; k < r; k++ {
		go func(j int, in chan int) {
			for {
				u := <-in
				res, err := http.DefaultClient.Do(reqs[u])
				if err != nil {
					fmt.Println("err:", err)
					continue
				}

				body, _ := ioutil.ReadAll(res.Body)
				res.Body.Close()

				fmt.Println(u, string(body[:10]))
			}
		}(k, in)
	}

	for i := 0; i < 1000000; i++ {
		in <- rand.Intn(len(reqs))
	}
	// return string(body)
}
