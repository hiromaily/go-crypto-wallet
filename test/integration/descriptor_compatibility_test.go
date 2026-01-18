//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type descriptorRecord struct {
	Desc      string `json:"desc"`
	Timestamp any    `json:"timestamp,omitempty"`
	Active    bool   `json:"active,omitempty"`
	Range     []int  `json:"range,omitempty"`
}

func TestDescriptorCompatibility_ImportAndDerive(t *testing.T) {
	descriptorFile := os.Getenv("BTC_CORE_COMPAT_DESCRIPTOR_FILE")
	if descriptorFile == "" {
		t.Skip("set BTC_CORE_COMPAT_DESCRIPTOR_FILE to run Bitcoin Core compatibility tests")
	}

	cli := newBitcoinCLI()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	version, err := cli.getVersion(ctx)
	if err != nil {
		t.Fatalf("fetch Bitcoin Core version: %v", err)
	}
	t.Logf("Bitcoin Core version: %s", version)

	records, err := loadDescriptorRecords(descriptorFile)
	if err != nil {
		t.Fatalf("load descriptors: %v", err)
	}
	if len(records) == 0 {
		t.Fatalf("no descriptors found in %s", descriptorFile)
	}

	for i, rec := range records {
		t.Run(fmt.Sprintf("descriptor_%d", i), func(t *testing.T) {
			t.Helper()

			if err := cli.importDescriptors(ctx, []descriptorRecord{rec}); err != nil {
				t.Fatalf("import descriptor: %v", err)
			}

			addrs, err := cli.deriveAddresses(ctx, rec.Desc, 0, 2)
			if err != nil {
				t.Fatalf("derive addresses: %v", err)
			}
			if len(addrs) == 0 {
				t.Fatalf("no addresses derived from descriptor")
			}
			t.Logf("derived addresses: %v", addrs)
		})
	}
}

type bitcoinCLI struct {
	bin  string
	args []string
}

func newBitcoinCLI() bitcoinCLI {
	bin := os.Getenv("BITCOIN_CLI")
	if bin == "" {
		bin = "bitcoin-cli"
	}
	extra := os.Getenv("BITCOIN_CLI_ARGS")
	var args []string
	if extra != "" {
		args = strings.Fields(extra)
	}
	return bitcoinCLI{bin: bin, args: args}
}

func (c bitcoinCLI) run(ctx context.Context, subCmd string, extraArgs ...string) ([]byte, error) {
	cmdArgs := append([]string{}, c.args...)
	cmdArgs = append(cmdArgs, subCmd)
	cmdArgs = append(cmdArgs, extraArgs...)

	cmd := exec.CommandContext(ctx, c.bin, cmdArgs...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("bitcoin-cli %s failed: %w\n%s", subCmd, err, string(out))
	}
	return out, nil
}

func (c bitcoinCLI) getVersion(ctx context.Context) (string, error) {
	out, err := c.run(ctx, "-version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c bitcoinCLI) importDescriptors(ctx context.Context, recs []descriptorRecord) error {
	payload, err := json.Marshal(recs)
	if err != nil {
		return fmt.Errorf("marshal descriptors: %w", err)
	}
	_, err = c.run(ctx, "importdescriptors", string(payload))
	return err
}

func (c bitcoinCLI) deriveAddresses(ctx context.Context, descriptor string, start, count uint32) ([]string, error) {
	end := int(start) + int(count) - 1
	if count == 0 {
		end = int(start) - 1
	}
	rangeArg := fmt.Sprintf("[%d,%d]", start, end)
	out, err := c.run(ctx, "deriveaddresses", descriptor, rangeArg)
	if err != nil {
		return nil, err
	}

	var addrs []string
	if err := json.Unmarshal(out, &addrs); err != nil {
		return nil, fmt.Errorf("decode deriveaddresses response: %w", err)
	}

	return addrs, nil
}

func loadDescriptorRecords(path string) ([]descriptorRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	recs := make([]descriptorRecord, 0)
	if err := json.Unmarshal(data, &recs); err == nil && len(recs) > 0 {
		return recs, nil
	}

	// Fallback: treat as plain descriptors (one per line)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		recs = append(recs, descriptorRecord{
			Desc:      line,
			Timestamp: "now",
		})
	}

	return recs, nil
}
