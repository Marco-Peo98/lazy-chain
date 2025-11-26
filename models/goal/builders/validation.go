package builders

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// =============================================================================
// VALIDATION FUNCTIONS
// =============================================================================

// ValidationResult holds the result of a validation check
type ValidationResult struct {
	Valid   bool
	Message string
}

// Valid returns a successful validation result
func Valid() ValidationResult {
	return ValidationResult{Valid: true, Message: ""}
}

// Invalid returns a failed validation result with a message
func Invalid(message string) ValidationResult {
	return ValidationResult{Valid: false, Message: message}
}

// =============================================================================
// ADDRESS VALIDATION
// =============================================================================

// Base32 alphabet used by Algorand (RFC 4648 without padding)
const base32Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// ValidateAlgorandAddress validates an Algorand address
// Algorand addresses are 58 characters, base32 encoded (uppercase + 2-7)
func ValidateAlgorandAddress(address string) ValidationResult {
	address = strings.TrimSpace(address)

	// Empty is handled separately (required check)
	if address == "" {
		return Valid() // Let required check handle this
	}

	// Check length
	if len(address) != 58 {
		return Invalid(fmt.Sprintf("address must be 58 characters (got %d)", len(address)))
	}

	// Check characters (base32: A-Z, 2-7)
	for i, c := range address {
		if !strings.ContainsRune(base32Alphabet, c) {
			return Invalid(fmt.Sprintf("invalid character '%c' at position %d (must be A-Z or 2-7)", c, i+1))
		}
	}

	// Basic checksum validation (last 4 bytes are checksum)
	// Full validation would require decoding and SHA512/256, but this catches most errors
	// For now, we just validate the format

	return Valid()
}

// =============================================================================
// AMOUNT VALIDATION
// =============================================================================

// ValidateAmount validates an amount (positive integer, microAlgos)
func ValidateAmount(amount string) ValidationResult {
	amount = strings.TrimSpace(amount)

	if amount == "" {
		return Valid() // Let required check handle this
	}

	// Check if it's a valid positive integer
	val, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		// Check for common mistakes
		if strings.Contains(amount, ".") {
			return Invalid("amount must be in microAlgos (whole number, no decimals)")
		}
		if strings.Contains(amount, "-") {
			return Invalid("amount must be positive")
		}
		if strings.Contains(amount, ",") {
			return Invalid("use plain number without commas")
		}
		return Invalid("amount must be a positive whole number")
	}

	// Check for zero (usually not valid for transfers)
	if val == 0 {
		return Invalid("amount must be greater than 0")
	}

	// Warn about very large amounts (> 10 billion Algos worth)
	// 10B Algos = 10_000_000_000_000_000 microAlgos
	if val > 10_000_000_000_000_000 {
		// Not invalid, but suspicious - could add a warning field
	}

	return Valid()
}

// ValidateAmountAllowZero validates an amount that can be zero (for total supply, etc.)
func ValidateAmountAllowZero(amount string) ValidationResult {
	amount = strings.TrimSpace(amount)

	if amount == "" {
		return Valid()
	}

	_, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		if strings.Contains(amount, ".") {
			return Invalid("must be a whole number (no decimals)")
		}
		if strings.Contains(amount, "-") {
			return Invalid("must be positive or zero")
		}
		return Invalid("must be a positive whole number")
	}

	return Valid()
}

// =============================================================================
// INTEGER VALIDATION
// =============================================================================

// ValidateInteger validates a positive integer
func ValidateInteger(value string) ValidationResult {
	value = strings.TrimSpace(value)

	if value == "" {
		return Valid()
	}

	_, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		if strings.Contains(value, "-") {
			return Invalid("must be positive")
		}
		return Invalid("must be a valid number")
	}

	return Valid()
}

// ValidateDecimals validates decimals (0-19 for Algorand ASAs)
func ValidateDecimals(value string) ValidationResult {
	value = strings.TrimSpace(value)

	if value == "" {
		return Valid()
	}

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return Invalid("must be a number between 0 and 19")
	}

	if val > 19 {
		return Invalid("decimals cannot exceed 19")
	}

	return Valid()
}

// ValidateRoundNumber validates a round number
func ValidateRoundNumber(value string) ValidationResult {
	value = strings.TrimSpace(value)

	if value == "" {
		return Valid()
	}

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return Invalid("must be a valid round number")
	}

	// Round 0 is technically valid but suspicious
	if val == 0 {
		return Invalid("round number should be greater than 0")
	}

	return Valid()
}

// =============================================================================
// APPLICATION-SPECIFIC VALIDATION
// =============================================================================

// ValidateAppID validates an application ID
func ValidateAppID(value string) ValidationResult {
	value = strings.TrimSpace(value)

	if value == "" {
		return Valid()
	}

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return Invalid("must be a valid application ID (number)")
	}

	if val == 0 {
		return Invalid("application ID must be greater than 0")
	}

	return Valid()
}

// ValidateAssetID validates an asset ID
func ValidateAssetID(value string) ValidationResult {
	value = strings.TrimSpace(value)

	if value == "" {
		return Valid()
	}

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return Invalid("must be a valid asset ID (number)")
	}

	if val == 0 {
		return Invalid("asset ID must be greater than 0")
	}

	return Valid()
}

// =============================================================================
// BASE64 VALIDATION (for keys)
// =============================================================================

// base64Regex matches valid base64 strings
var base64Regex = regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)

// ValidateBase64 validates a base64 encoded string
func ValidateBase64(value string) ValidationResult {
	value = strings.TrimSpace(value)

	if value == "" {
		return Valid()
	}

	// Check length is multiple of 4 or can be padded
	if len(value)%4 != 0 {
		// Try to check if it's valid without padding
		padded := value + strings.Repeat("=", (4-len(value)%4)%4)
		if !base64Regex.MatchString(padded) {
			return Invalid("must be valid base64 encoded")
		}
	} else if !base64Regex.MatchString(value) {
		return Invalid("must be valid base64 encoded")
	}

	return Valid()
}

// ValidateVotingKey validates an Algorand voting key (specific length)
func ValidateVotingKey(value string) ValidationResult {
	value = strings.TrimSpace(value)

	if value == "" {
		return Valid()
	}

	// Voting keys are 32 bytes = 44 base64 chars (with padding) or 43 without
	baseResult := ValidateBase64(value)
	if !baseResult.Valid {
		return baseResult
	}

	// Approximate length check (32 bytes base64 = ~43-44 chars)
	if len(value) < 40 || len(value) > 50 {
		return Invalid("voting key should be ~44 characters (32 bytes base64)")
	}

	return Valid()
}

// =============================================================================
// FIELD VALIDATION DISPATCHER
// =============================================================================

// ValidateField validates a field based on its type
func ValidateField(field *FieldDef) ValidationResult {
	if field == nil {
		return Valid()
	}

	// Skip validation for bool and select types
	if field.Type == FieldTypeBool || field.Type == FieldTypeSelect {
		return Valid()
	}

	// Get the value to validate
	value := ""
	if field.Value != nil {
		value = *field.Value
	}

	// Empty values pass validation (required check is separate)
	if strings.TrimSpace(value) == "" {
		return Valid()
	}

	// Dispatch based on field type
	switch field.Type {
	case FieldTypeAddress:
		return ValidateAlgorandAddress(value)
	case FieldTypeAmount:
		return ValidateAmount(value)
	case FieldTypeInteger:
		// Check for specific field keys that need special validation
		switch field.Key {
		case "appID":
			return ValidateAppID(value)
		case "assetID":
			return ValidateAssetID(value)
		case "decimals":
			return ValidateDecimals(value)
		case "firstVotingRound", "lastVotingRound":
			return ValidateRoundNumber(value)
		default:
			return ValidateInteger(value)
		}
	case FieldTypeText:
		// Check for specific field keys that need validation
		switch field.Key {
		case "votingKey", "selectionKey":
			return ValidateVotingKey(value)
		case "stateProofKey":
			return ValidateBase64(value)
		}
		return Valid() // Text fields generally don't need validation
	default:
		return Valid()
	}
}

// =============================================================================
// BATCH VALIDATION
// =============================================================================

// ValidateAllFields validates all fields and returns the first error
func ValidateAllFields(fs *FieldsState) ValidationResult {
	if fs == nil {
		return Valid()
	}

	for i := range fs.Fields {
		field := &fs.Fields[i]
		result := ValidateField(field)
		if !result.Valid {
			return Invalid(fmt.Sprintf("%s: %s", field.Label, result.Message))
		}
	}

	return Valid()
}

// HasValidationErrors checks if any field has a validation error
func HasValidationErrors(fs *FieldsState) bool {
	if fs == nil {
		return false
	}

	for i := range fs.Fields {
		if fs.Fields[i].ValidationError != "" {
			return true
		}
	}

	return false
}
