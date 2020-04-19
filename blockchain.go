package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type myJSON struct {
	Array []Block
}

// Block define the struct of a default Block
type Block struct {
	Index        int    `json:"Index"`
	Nonce        int    `json:"Nonce"`
	PreviousHash string `json:"PreviousHash"`
	Hash         string `json:"Hash"`
	Timestamp    int64  `json:"Timestamp"`
	Difficulty   string `json:"Difficulty"`
}

//If there is no Blockchain or no Previous Block, the genesisBlock function generates the first Block of the Blockchain
func genesisBlock(gHash string, gTimestamp int64, gNonce int, gDifficulty string) {
	block := []Block{}
	newBlock := &Block{
		Index:        1,
		PreviousHash: " ",
		Hash:         gHash,
		Timestamp:    gTimestamp,
		Nonce:        gNonce,
		Difficulty:   gDifficulty,
	}
	block = append(block, *newBlock)
	file, err := json.MarshalIndent(block, "", " ")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	_ = ioutil.WriteFile("Blockchain.json", file, 0644)

}
