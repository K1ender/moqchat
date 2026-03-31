package main

import (
	"context"

	"github.com/K1ender/moqchat/pkg/api"
)

func main() {
	ctx := context.Background()

	err := api.Run(ctx)
	if err != nil {
		panic(err)
	}
}
