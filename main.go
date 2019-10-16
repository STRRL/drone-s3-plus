package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	plugin := &Plugin{}
	err := envconfig.Process("PLUGIN", plugin)
	if err != nil {
		fmt.Println("init plugin failed,", err)
		os.Exit(1)
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println("plugin exec failed,", err)
		os.Exit(1)
	}
}
