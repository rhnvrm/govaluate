package govaluate

import (
	"testing"
)

// Tests to make sure that all the different token kinds have different string representations
// Gotta get that 95% code coverage yall. That's why tests like this get written; over-reliance on bad metrics.
func TestTokenKindStrings(test *testing.T) {

	var kindStrings []string
	var kindString string

	kinds := []TokenKind{
		unknown,
		prefix,
		numeric,
		boolean,
		stringToken,
		pattern,
		timeToken,
		variable,
		comparator,
		logicalop,
		modifier,
		clause,
		clauseClose,
		ternary,
	}

	for _, kind := range kinds {

		kindString = kind.String()

		for _, extantKind := range kindStrings {
			if extantKind == kindString {
				test.Logf("Token kind test found duplicate string for token kind %v ('%v')\n", kind, kindString)
				test.Fail()
				return
			}
		}

		kindStrings = append(kindStrings, kindString)
	}
}
