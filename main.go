package main

import (
	"fmt"

	"github.com/nickolation/pointsalvor-sdk/pointsalvor"
)

func main() {
	//create new agent with token-api: api-token
	agent, err := pointsalvor.NewAgent("<token-api>")
	if err != nil {
		fmt.Println(err.Error())
	}
}
