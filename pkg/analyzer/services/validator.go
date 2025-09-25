package services

import (
	"fmt"

	"github.com/cherry-pick/pkg/analyzer/core"
)

type ValidatorService struct{}

func NewValidatorService() *ValidatorService {
	return &ValidatorService{}
}

func (vs *ValidatorService) ValidateRequest(request core.AnalysisRequest) error {
	if request.DatabaseType == "" {
		return fmt.Errorf("database type is required")
	}

	if err := vs.ValidateDatabaseType(request.DatabaseType); err != nil {
		return err
	}

	if err := vs.ValidateOptions(request.Options); err != nil {
		return err
	}

	return nil
}

func (vs *ValidatorService) ValidateDatabaseType(dbType core.DatabaseType) error {
	supportedTypes := []core.DatabaseType{
		core.DatabaseTypeMySQL,
		core.DatabaseTypePostgres,
		core.DatabaseTypeSQLite,
		core.DatabaseTypeMongoDB,
	}

	for _, supportedType := range supportedTypes {
		if dbType == supportedType {
			return nil
		}
	}

	return fmt.Errorf("unsupported database type: %s", dbType)
}

func (vs *ValidatorService) ValidateOptions(options core.AnalysisOptions) error {
	if options.SampleSize < 0 {
		return fmt.Errorf("sample size cannot be negative")
	}

	if options.SampleSize > 10000 {
		return fmt.Errorf("sample size cannot exceed 10000")
	}

	if options.MaxCollections < 0 {
		return fmt.Errorf("max collections cannot be negative")
	}

	if options.MaxCollections > 1000 {
		return fmt.Errorf("max collections cannot exceed 1000")
	}

	return nil
}
