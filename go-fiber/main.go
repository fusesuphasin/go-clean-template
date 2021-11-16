package main

import (
	"context"
	"fmt"

	"github.com/fusesuphasin/go-fiber/app/bootstrap"
	"github.com/fusesuphasin/go-fiber/app/infrastructure"
)

func main() {
	ctx := context.Background()
	logger := infrastructure.NewLogger()
	
	infrastructure.Load(logger)
	infrastructure.Open()
	enforcer, err := infrastructure.NewMongoHandler(ctx)
	if err != nil {
		/* logger.LogError("%s", err) */
		fmt.Println("Error: ", err)
	}
	bootstrap.Dispatch(ctx, logger, enforcer)
}