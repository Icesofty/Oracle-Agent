package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/sha3"
)

// Main mining function
func mining() {
	block := new(Block)

	for {
		getJSON(block)
		time.Sleep(5 * time.Second)
		color.Cyan.Printf("Mining...")
		mineBlock()
	}
}

func hexDecode(ht []byte) string {
	return hex.EncodeToString(ht)
}

// Function that calculate if a Hash is valid or not based on the Oracle difficulty
func isHashValide(hash []byte, difficulty string) bool {
	hashDecoded := hexDecode(hash)

	if strings.HasPrefix(hashDecoded, difficulty) {

		return true
	}
	return false

}

// Inifinte loop
func mineBlock() {

	clientBlockchain := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+"/blocks.json", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, getErr := clientBlockchain.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if getErr != nil {
		log.Fatal(getErr)
	}

	defer resp.Body.Close()

	/* Read json file on server */

	blockchain, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Take the actual length of the Json
	index := gjson.Get(string(blockchain), "#")
	// Get the actual Index of the current Block
	currentHash := strconv.Itoa(int(index.Int())-1) + ".Hash"
	previousHash := strconv.Itoa(int(index.Int())-1) + ".PreviousHash"
	curentTimestamp := strconv.Itoa(int(index.Int())-1) + ".Timestamp"
	curentDifficulty := strconv.Itoa(int(index.Int())-1) + ".Difficulty"

	currentHashData := gjson.Get(string(blockchain), currentHash)
	previousHashData := gjson.Get(string(blockchain), previousHash)
	previousTimestamp := gjson.Get(string(blockchain), curentTimestamp)
	previousDifficulty := gjson.Get(string(blockchain), curentDifficulty)

	// Define a new Block based on the previous block datas (index - 1) - PreviousHash
	block := []Block{}
	json.Unmarshal(blockchain, &block)

	// Creating the new block struct to be passed into the POW function
	newBlock := &Block{
		Index:        int(index.Int()) - 1,
		PreviousHash: previousHashData.String(),
		Hash:         currentHashData.String(),
		Timestamp:    previousTimestamp.Int(),
		Difficulty:   previousDifficulty.String(),
	}

	// Finding the nonce of the current block
	newBlock.Nonce, newBlock.Hash = findPOW(newBlock)

}

func findPOW(block *Block) (int, string) {
	start := time.Now()
	nonce := 0
	// Block datas to hash
	record := block.Hash + block.PreviousHash + strconv.Itoa(int(block.Timestamp)) + strconv.Itoa(block.Index) + strconv.Itoa(nonce)
	h := sha3.New256()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	//POW woking function - Nonce will be incremented by one each round until the Hash starts with number of 0
	// defined in the difficulty
	for {
		//if true, the POW has been found - Hash of Data + Nonce begin with n 0 (n is based on the difficulty)
		if isHashValide(hashed, block.Difficulty) {
			//Calculate total time to find the Nonce
			elapsed := time.Since(start)
			color.Green.Print("Founded!\n")
			log.Printf("Total time :\n %s", elapsed)

			hashDecoded := hexDecode(hashed)
			fmt.Println(nonce)

			sendNonce(nonce)

			return nonce, hashDecoded
		} else {
			nonce++
			// Reset the datas from origin, and add the new nonce
			record = block.Hash + block.PreviousHash + strconv.Itoa(int(block.Timestamp)) + strconv.Itoa(block.Index) + strconv.Itoa(nonce)
			h = sha3.New256()
			h.Write([]byte(record))
			hashed = h.Sum(nil)

		}
	}

}
