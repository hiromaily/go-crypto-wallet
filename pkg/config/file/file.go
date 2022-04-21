package file

import (
	"fmt"
	"os"
)

func GetConfigFilePath(fileName string) string {
	// for github action
	if os.Getenv("GITHUB_WORKSPACE") != "" {
		confPath := fmt.Sprintf("%s/data/config/%s", os.Getenv("GITHUB_WORKSPACE"), fileName)
		if f, err := os.Stat(confPath); err == nil && !f.IsDir() {
			return confPath
		}
	}
	if os.Getenv("GOPATH") != "" {
		projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
		confPath := fmt.Sprintf("%s/data/config/%s", projPath, fileName)
		if f, err := os.Stat(confPath); err == nil && !f.IsDir() {
			return confPath
		}
	}
	return ""
}
