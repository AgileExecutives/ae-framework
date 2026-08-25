package services

import (
	"os"
	"path/filepath"
	"testing"
)

func repoRootFromPwd() string {
	cwd, _ := os.Getwd()
	// walk up until we find version.json or Makefile or reach filesystem root
	p := cwd
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(p, "version.json")); err == nil {
			return p
		}
		if _, err := os.Stat(filepath.Join(p, "Makefile")); err == nil {
			return p
		}
		p = filepath.Dir(p)
	}
	return cwd
}

// regImpl is a minimal fake registrar for tests.
type regImpl struct{ calls *[]string }

func (r regImpl) RegisterContractFromFile(tenantID uint, module string, path string) error {
	*r.calls = append(*r.calls, module+":"+path)
	return nil
}

func TestRegisterInvoiceContractsRegistersFile(t *testing.T) {
	root := repoRootFromPwd()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir repo root: %v", err)
	}

	// ensure contracts dir
	contractsDir := filepath.Join(root, "shared-modules", "invoice", "contracts")
	if err := os.MkdirAll(contractsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f := filepath.Join(contractsDir, "std_invoice-contract.json")
	content := `{"type":"object","properties":{"Invoice":{"type":"object"}},"template_key":"invoice"}`
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	// Use a fake registrar to assert the module attempts to register the expected file
	var called []string
	registrar := regImpl{calls: &called}

	if err := RegisterInvoiceContracts(registrar, 7); err != nil {
		t.Fatalf("RegisterInvoiceContracts: %v", err)
	}
	if len(called) == 0 {
		t.Fatalf("expected RegisterInvoiceContracts to call RegisterContractFromFile at least once")
	}
}
