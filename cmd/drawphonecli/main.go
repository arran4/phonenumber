package main

import (
	"log"
	"phonenumber"
)

func main() {
	s := phonenumber.Numbers("Hello how are you?")
	if err := phonenumber.DrawPhone(s, "out.png"); err != nil {
		log.Printf("Error: %s", err)
		return
	}
}
