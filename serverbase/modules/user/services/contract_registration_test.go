package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterUserContractsRegistersFiles(t *testing.T) {
	// work from repo root
	cwd, _ := os.Getwd()
	root := filepath.Clean(filepath.Join(cwd, "..", "..", ".."))
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir repo root: %v", err)
	}

	contractsDir := filepath.Join(root, "modules", "user", "contracts")
	if err := os.MkdirAll(contractsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	welcome := filepath.Join(contractsDir, "welcome-contract.json")
	if err := os.WriteFile(welcome, []byte(`{"type":"object","properties":{"X":{"type":"string"}},"template_key":"welcome"}`), 0o644); err != nil {
		t.Fatalf("write welcome: %v", err)
	}
	pwd := filepath.Join(contractsDir, "password_reset-contract.json")
	if err := os.WriteFile(pwd, []byte(`{"type":"object","properties":{"X":{"type":"string"}},"template_key":"password_reset"}`), 0o644); err != nil {
		t.Fatalf("write pwd: %v", err)
	}

	// Use fake registrar to assert both files attempted
	var called []string
	registrar := regImpl{calls: &called}

	if err := RegisterUserContracts(registrar, 3); err != nil {
		t.Fatalf("RegisterUserContracts: %v", err)
	}
	if len(called) < 2 {
		t.Fatalf("expected 2 calls to RegisterContractFromFile, got %d", len(called))
	}
}

// regImpl is a minimal fake registrar for tests.
type regImpl struct{ calls *[]string }

func (r regImpl) RegisterContractFromFile(tenantID uint, module string, path string) error {
	*r.calls = append(*r.calls, module+":"+path)
	return nil
}
