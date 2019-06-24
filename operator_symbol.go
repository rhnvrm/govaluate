package govaluate

// OperatorSymbol represents the valid symbols for operators.
type OperatorSymbol int

const (
	value OperatorSymbol = iota
	literal
	noopSymbol
	eq
	neq
	gt
	lt
	gte
	lte
	req
	nreq
	in

	and
	or

	plus
	minus
	bitwiseAnd
	bitwiseOr
	bitwiseXor
	bitwiseLshift
	bitwiseRshift
	multiply
	divide
	modulus
	exponent

	negate
	invert
	bitwiseNot

	ternaryTrue
	ternaryFalse
	coalesce

	functional
	access
	separate
)

type operatorPrecedence int

const (
	noopPrecedence operatorPrecedence = iota
	valuePrecedence
	functionalPrecedence
	prefixPrecedence
	exponentialPrecedence
	additivePrecedence
	bitwisePrecedence
	bitwiseShiftPrecedence
	multiplicativePrecedence
	comparatorPrecedence
	ternaryPrecedence
	logicalAndPrecedence
	logicalOrPrecedence
	separatePrecedence
)

func findOperatorPrecedenceForSymbol(symbol OperatorSymbol) operatorPrecedence {
	switch symbol {
	case noopSymbol:
		return noopPrecedence
	case value:
		return valuePrecedence
	case eq:
		fallthrough
	case neq:
		fallthrough
	case gt:
		fallthrough
	case lt:
		fallthrough
	case gte:
		fallthrough
	case lte:
		fallthrough
	case req:
		fallthrough
	case nreq:
		fallthrough
	case in:
		return comparatorPrecedence
	case and:
		return logicalAndPrecedence
	case or:
		return logicalOrPrecedence
	case bitwiseAnd:
		fallthrough
	case bitwiseOr:
		fallthrough
	case bitwiseXor:
		return bitwisePrecedence
	case bitwiseLshift:
		fallthrough
	case bitwiseRshift:
		return bitwiseShiftPrecedence
	case plus:
		fallthrough
	case minus:
		return additivePrecedence
	case multiply:
		fallthrough
	case divide:
		fallthrough
	case modulus:
		return multiplicativePrecedence
	case exponent:
		return exponentialPrecedence
	case bitwiseNot:
		fallthrough
	case negate:
		fallthrough
	case invert:
		return prefixPrecedence
	case coalesce:
		fallthrough
	case ternaryTrue:
		fallthrough
	case ternaryFalse:
		return ternaryPrecedence
	case access:
		fallthrough
	case functional:
		return functionalPrecedence
	case separate:
		return separatePrecedence
	default:
		return valuePrecedence
	}
}

// Map of all valid comparators, and their string equivalents.
// Used during parsing of expressions to determine if a symbol is, in fact, a comparator.
// Also used during evaluation to determine exactly which comparator is being used.
var comparatorSymbols = map[string]OperatorSymbol{
	"==": eq,
	"!=": neq,
	">":  gt,
	">=": gte,
	"<":  lt,
	"<=": lte,
	"=~": req,
	"!~": nreq,
	"in": in,
}

var logicalSymbols = map[string]OperatorSymbol{
	"&&": and,
	"||": or,
}

var bitwiseSymbols = map[string]OperatorSymbol{
	"^": bitwiseXor,
	"&": bitwiseAnd,
	"|": bitwiseOr,
}

var bitwiseShiftSymbols = map[string]OperatorSymbol{
	">>": bitwiseRshift,
	"<<": bitwiseLshift,
}

var additiveSymbols = map[string]OperatorSymbol{
	"+": plus,
	"-": minus,
}

var multiplicativeSymbols = map[string]OperatorSymbol{
	"*": multiply,
	"/": divide,
	"%": modulus,
}

var exponentialSymbolsS = map[string]OperatorSymbol{
	"**": exponent,
}

var prefixSymbols = map[string]OperatorSymbol{
	"-": negate,
	"!": invert,
	"~": bitwiseNot,
}

var ternarySymbols = map[string]OperatorSymbol{
	"?":  ternaryTrue,
	":":  ternaryFalse,
	"??": coalesce,
}

// this is defined separately from additiveSymbols et al because it's needed for parsing, not stage planning.
var modifierSymbols = map[string]OperatorSymbol{
	"+":  plus,
	"-":  minus,
	"*":  multiply,
	"/":  divide,
	"%":  modulus,
	"**": exponent,
	"&":  bitwiseAnd,
	"|":  bitwiseOr,
	"^":  bitwiseXor,
	">>": bitwiseRshift,
	"<<": bitwiseLshift,
}

var separatorSymbols = map[string]OperatorSymbol{
	",": separate,
}

// IsModifierType returns true if this operator is contained by the given array of candidate symbols.
// False otherwise.
func (os OperatorSymbol) IsModifierType(candidate []OperatorSymbol) bool {
	for _, symbolType := range candidate {
		if os == symbolType {
			return true
		}
	}
	return false
}

// String is generally used when formatting type check errors.
// We could store the stringified symbol somewhere else and not require a duplicated codeblock to translate
// OperatorSymbol to string, but that would require more memory, and another field somewhere.
// Adding operators is rare enough that we just stringify it here instead.
func (os OperatorSymbol) String() string {
	switch os {
	case noopSymbol:
		return "NOOP"
	case value:
		return "value"
	case eq:
		return "="
	case neq:
		return "!="
	case gt:
		return ">"
	case lt:
		return "<"
	case gte:
		return ">="
	case lte:
		return "<="
	case req:
		return "=~"
	case nreq:
		return "!~"
	case and:
		return "&&"
	case or:
		return "||"
	case in:
		return "in"
	case bitwiseAnd:
		return "&"
	case bitwiseOr:
		return "|"
	case bitwiseXor:
		return "^"
	case bitwiseLshift:
		return "<<"
	case bitwiseRshift:
		return ">>"
	case plus:
		return "+"
	case minus:
		return "-"
	case multiply:
		return "*"
	case divide:
		return "/"
	case modulus:
		return "%"
	case exponent:
		return "**"
	case negate:
		return "-"
	case invert:
		return "!"
	case bitwiseNot:
		return "~"
	case ternaryTrue:
		return "?"
	case ternaryFalse:
		return ":"
	case coalesce:
		return "??"
	default:
		return ""
	}
}
