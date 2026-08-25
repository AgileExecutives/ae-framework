package services

import "path/filepath"

// Registrar is the minimal interface required for registering contracts.
type Registrar interface {
	RegisterContractFromFile(tenantID uint, module string, path string) error
}

// RegisterUserContracts registers contracts for the user module
func RegisterUserContracts(contractRegistrar Registrar, tenantID uint) error {
	contractsDir := "modules/user/contracts"
	contracts := []string{
		"welcome-contract.json",
		"password_reset-contract.json",
	}
	for _, contractFile := range contracts {
		contractPath := filepath.Join(contractsDir, contractFile)
		if err := contractRegistrar.RegisterContractFromFile(tenantID, "user", contractPath); err != nil {
			return err
		}
	}
	return nil
}
