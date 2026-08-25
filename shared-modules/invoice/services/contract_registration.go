package services

import "path/filepath"

// Registrar is the minimal interface required for registering contracts.
type Registrar interface {
	RegisterContractFromFile(tenantID uint, module string, path string) error
}

// RegisterInvoiceContracts registers all invoice module template contracts
func RegisterInvoiceContracts(contractRegistrar Registrar, tenantID uint) error {
	// Get the module's contracts directory
	contractsDir := "shared-modules/invoice/contracts"

	// Register all contract files
	contracts := []string{
		"std_invoice-contract.json",
	}

	for _, contractFile := range contracts {
		contractPath := filepath.Join(contractsDir, contractFile)
		if err := contractRegistrar.RegisterContractFromFile(tenantID, "invoice", contractPath); err != nil {
			return err
		}
	}

	return nil
}
