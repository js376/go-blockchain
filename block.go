package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

//Block with four fields
type Block struct {
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	Hash      []byte
	Nonce     int
}

//Serialize a data
func (b *Block) Serialize() []byte {
	var result bytes.Buffer // a buffer will store a serialized data
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)

	return result.Bytes()
}

//DeserializeBlock is to deserialize
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	return &block
}

//NewBlock is a function to create a new block
func NewBlock(_data string, _prevHash []byte) *Block {
	_block := &Block{
		Timestamp: time.Now().Unix(),
		Data:      []byte(_data),
		PrevHash:  _prevHash,
		Hash:      []byte{},
	}

	pow := NewProofOfWork(_block)
	nonce, hash := pow.Run()

	_block.Nonce = nonce
	_block.Hash = hash[:]

	return _block
}

//NewGenesisBlock is a function to create a genesis block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis", []byte{})
}

// functions for BlockChain

//AddBlock adds block to the blockchain
func (bc *Blockchain) AddBlock(_data string) {

	// prevBlock := bc.blocks[len(bc.blocks)-1]
	// newBlock := NewBlock(_data, prevBlock.Hash)
	// bc.blocks = append(bc.blocks, newBlock)

	var lastHash []byte

	// View is read-only type of boltDb transaction
	err = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(_data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}
