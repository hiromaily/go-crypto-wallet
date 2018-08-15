package api_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromaily/go-bitcoin/pkg/service"
)

var (
	wlt      *service.Wallet
	confPath = flag.String("conf", "../../data/toml/config.toml", "Path for configuration toml file")
)

func setup() {
	if wlt != nil {
		return
	}

	flag.Parse()

	var err error
	wlt, err = service.InitialSettings(*confPath)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	wlt.Done()
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

// TestValidateAddress
func TestValidateAddress(t *testing.T) {
	var tests = []struct {
		addr  string
		isErr bool
	}{
		{"2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr", false},
		{"2NDGkbQTwg2v1zP6yHZw3UJhmsBh9igsSos", false},
		{"4VHGkbQTGg2vN5P6yHZw3UJhmsBh9igsSos", true},
	}

	for _, val := range tests {
		//t.Logf("check address: %s", val.addr)
		fmt.Printf("check address: %s\n", val.addr)

		err := wlt.Btc.ValidateAddress(val.addr)
		if err != nil && !val.isErr {
			t.Errorf("Unexpectedly error occorred. %v", err)
		}
		if err == nil && val.isErr {
			t.Error("Error is expected. However nothing happened.")
		}
	}
}
