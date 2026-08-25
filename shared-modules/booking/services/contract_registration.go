package services

import "path/filepath"

// Registrar is the minimal interface required for registering contracts.
type Registrar interface {
	RegisterContractFromFile(tenantID uint, module string, path string) error
}

// RegisterBookingContracts registers all booking template contracts with the template system
func RegisterBookingContracts(contractRegistrar Registrar, tenantID uint) error {
	// Get the module's contracts directory
	contractsDir := "shared-modules/booking/contracts"

	// Register all contract files
	contracts := []string{
		"booking_confirmation-contract.json",
	}

	for _, contractFile := range contracts {
		contractPath := filepath.Join(contractsDir, contractFile)
		if err := contractRegistrar.RegisterContractFromFile(tenantID, "booking", contractPath); err != nil {
			return err
		}
	}

	return nil
}
