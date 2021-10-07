package main

import (
	"context"
	"fmt"

	"github.com/nickolation/pointsalvor-sdk/pointsalvor"
)

func main() {
	//create new agent with token-api: api-token
	agent, err := pointsalvor.NewAgent("<api-token>")
	if err != nil {
		fmt.Println(err.Error())
	}

	ctx := context.Background()

	//get all project by agent
	res, err := agent.GetAllProjects(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
	fmt.Println(pointsalvor.BankIdProject)
}
