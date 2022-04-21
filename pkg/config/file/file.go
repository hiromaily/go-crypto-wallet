package file

import (
	"fmt"
	"os"
)

func GetConfigFilePath(fileName string) string {
	bases := []string{os.Getenv("github_workspace"), os.Getenv("GOPATH")}
	for _, base := range bases {
		projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", base)
		confPath := fmt.Sprintf("%s/data/config/%s", projPath, fileName)
		// check
		if f, err := os.Stat(confPath); os.IsNotExist(err) || f.IsDir() {
			continue
		}
		return confPath
	}
	return ""
}
