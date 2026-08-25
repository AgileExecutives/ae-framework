package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterBookingContractsRegistersFile(t *testing.T) {
	// change to repo root so paths used by the module resolve
	cwd, _ := os.Getwd()
	root := filepath.Clean(filepath.Join(cwd, "..", ".."))
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir repo root: %v", err)
	}

	contractsDir := filepath.Join(root, "shared-modules", "booking", "contracts")
	if err := os.MkdirAll(contractsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f := filepath.Join(contractsDir, "booking_confirmation-contract.json")
	content := `{"type":"object","properties":{"Booking":{"type":"object"}},"template_key":"booking_confirmation"}`
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	// Use fake registrar to assert booking module calls RegisterContractFromFile
	var called []string
	registrar := regImpl{calls: &called}

	if err := RegisterBookingContracts(registrar, 11); err != nil {
		t.Fatalf("RegisterBookingContracts: %v", err)
	}
	if len(called) == 0 {
		t.Fatalf("expected RegisterBookingContracts to call RegisterContractFromFile at least once")
	}
}

// regImpl is a minimal fake registrar for tests.
type regImpl struct{ calls *[]string }

func (r regImpl) RegisterContractFromFile(tenantID uint, module string, path string) error {
	*r.calls = append(*r.calls, module+":"+path)
	return nil
}
