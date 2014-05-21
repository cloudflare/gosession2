package main

import (
	"fmt"

	"code.google.com/p/go.crypto/bcrypt"
)

const cost = 13

func main() {
	password := []byte("password")

	hash, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		fmt.Printf("Bcrypt failed: %v\n", err)
		return
	}
	fmt.Printf("Hash: %x\n", hash)
	err = bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		fmt.Println("Hash and password don't match.")
		return
	}
	hash[8]++
	err = bcrypt.CompareHashAndPassword(hash, password)
	if err == nil {
		fmt.Println("Hash and password don't match.")
		return
	}
}
