package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gookit/color"
)

func sendNonce(nonce int) {
	str := baseURL + "block/" + strconv.Itoa(nonce)

	resp, err := http.Get(str)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	color.Green.Print(string(body))

	if string(body) == "approved" {
		fmt.Println(" - New Block Added to the Oracle, let's mine")
	}
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(target interface{}) error {
	r, err := myClient.Get(baseURL + "blocks.json")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
