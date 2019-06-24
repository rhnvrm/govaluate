package govaluate

// TokenKind represents all valid types of tokens that a token can be.
type TokenKind int

const (
	unknown TokenKind = iota

	prefix
	numeric
	boolean
	stringToken
	pattern
	timeToken
	variable
	function
	separator
	accessor

	comparator
	logicalop
	modifier

	clause
	clauseClose

	ternary
)

// GetTokenKindString returns a string that describes the given TokenKind.
// e.g., when passed the NUMERIC TokenKind, this returns the string "NUMERIC".
func (kind TokenKind) String() string {

	switch kind {

	case prefix:
		return "PREFIX"
	case numeric:
		return "NUMERIC"
	case boolean:
		return "BOOLEAN"
	case stringToken:
		return "STRING"
	case pattern:
		return "PATTERN"
	case timeToken:
		return "TIME"
	case variable:
		return "VARIABLE"
	case function:
		return "FUNCTION"
	case separator:
		return "SEPARATOR"
	case comparator:
		return "COMPARATOR"
	case logicalop:
		return "LOGICALOP"
	case modifier:
		return "MODIFIER"
	case clause:
		return "CLAUSE"
	case clauseClose:
		return "CLAUSE_CLOSE"
	case ternary:
		return "TERNARY"
	case accessor:
		return "ACCESSOR"
	}

	return "UNKNOWN"
}
