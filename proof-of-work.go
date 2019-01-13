package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const targetBits = 24

var (
	maxNonce = math.MaxInt64
)

// ProofOfWork struct
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork creates new pow
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)

	//target.Lsh sets target = target << uint(256-targetBits), which is shifting bits left by 256-targetBits
	//hexadecimal representation of target is...
	// 0x10000000000000000000000000000000000000000000000000000000000
	target.Lsh(target, uint(256-targetBits))

	//target is a upper boundary of a range

	pow := &ProofOfWork{
		block:  b,
		target: target,
	}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

//IntToHex converts integer to hex
func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

// Run is a function where actual PoW algorithm runs
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)

	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

// Validate proof of work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
