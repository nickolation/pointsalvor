package main

import (
	"context"
	"fmt"

	"github.com/nickolation/pointsalvor-sdk/pointsalvor"
)

func main() {
	agent, err := pointsalvor.NewAgent("7100d79356a82a940ef358b9332b183390317a0d")
	if err != nil {
		fmt.Println(err.Error())
	}

	ctx := context.Background()

	res, err := agent.GetAllProjects(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
}
