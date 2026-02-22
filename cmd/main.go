package main

import (
	"fmt"
	"github.com/slidebolt/plugin-framework"
	 "github.com/slidebolt/plugin-system/pkg/bundle"
	"github.com/slidebolt/plugin-sdk"
)

func main() {
	fmt.Println("Starting System Plugin Sidecar...")
	framework.Init()

	b, err := sdk.RegisterBundle("plugin-system")
	if err != nil {
		fmt.Printf("Failed to register bundle: %v\n", err)
		return
	}

	p := bundle.NewPlugin()
	if err := p.Init(b); err != nil {
		fmt.Printf("Failed to init plugin: %v\n", err)
		return
	}

	fmt.Println("System Plugin is running.")
	select {}
}