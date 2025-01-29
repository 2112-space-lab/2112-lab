package main

import (
	"context"

	"github.com/Elbujito/2112/src/app-service/internal/cmd"
)

// VERSION application version
var VERSION string = "0.0.1"

func main() {

	mainCtx := context.Background()
	cmd.Version = VERSION
	cmd.Execute(mainCtx)
}
