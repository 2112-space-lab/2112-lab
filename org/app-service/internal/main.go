package main

import (
	"context"

	"github.com/org/2112-space-lab/org/app-service/internal/cmd"
)

var VERSION string = "0.0.1"

func main() {
	mainCtx := context.Background()
	cmd.Version = VERSION
	cmd.Execute(mainCtx)
}
