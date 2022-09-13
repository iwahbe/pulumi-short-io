package main

import (
	"fmt"
	"os"

	p "github.com/pulumi/pulumi-go-provider"

	"github.com/pulumi/pulumi-short-io"
)

func main() {
	err := p.RunProvider("short-io", short.VERSION, short.Provider())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
