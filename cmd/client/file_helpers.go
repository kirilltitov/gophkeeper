package main

import (
	"fmt"
	"os"
)

func getConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic("Could not get OS specific config dir: " + err.Error())
	}

	result := fmt.Sprintf("%s/%s", configDir, appDir)

	if err := os.Mkdir(result, 0o770); err != nil && !os.IsExist(err) {
		panic(fmt.Sprintf("Could not create directory '%s' for client: %s", result, err.Error()))
	}

	return result
}
