package validate

import "meta_harness/internal/schema"

// Validator defines the interface for validation
type Validator interface {
	Validate(files []string) schema.ValidationResult
}

// CompositeValidator combines multiple validators
type CompositeValidator struct {
	validators []Validator
}

// NewCompositeValidator creates a new CompositeValidator
func NewCompositeValidator(validators ...Validator) *CompositeValidator {
	return &CompositeValidator{validators: validators}
}

// Validate runs all validators and combines results
func (c *CompositeValidator) Validate(files []string) schema.ValidationResult {
	finalResult := schema.ValidationResult{Pass: true}

	for _, v := range c.validators {
		result := v.Validate(files)
		finalResult.FilesChecked = append(finalResult.FilesChecked, result.FilesChecked...)
		finalResult.Errors = append(finalResult.Errors, result.Errors...)
		finalResult.Warnings = append(finalResult.Warnings, result.Warnings...)
		if !result.Pass {
			finalResult.Pass = false
		}
	}

	return finalResult
}
