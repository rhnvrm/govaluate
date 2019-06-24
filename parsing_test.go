package govaluate

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"
	"unicode"
)

// Represents a test of parsing all tokens correctly from a string
type TokenParsingTest struct {
	Name      string
	Input     string
	Functions map[string]ExpressionFunction
	Expected  []ExpressionToken
}

func TestConstantParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{
		{
			Name:  "Single numeric",
			Input: "1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{
			Name:  "Single two-digit numeric",
			Input: "50",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 50.0,
				},
			},
		},
		{
			Name:  "Zero",
			Input: "0",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 0.0,
				},
			},
		},
		{
			Name:  "One digit hex",
			Input: "0x1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{
			Name:  "Two digit hex",
			Input: "0x10",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 16.0,
				},
			},
		},
		{
			Name:  "Hex with lowercase",
			Input: "0xabcdef",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 11259375.0,
				},
			},
		},
		{
			Name:  "Hex with uppercase",
			Input: "0xABCDEF",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 11259375.0,
				},
			},
		},
		{

			Name:  "Single string",
			Input: "'foo'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
			},
		},
		{

			Name:  "Single time, RFC3339, only date",
			Input: "'2014-01-02'",
			Expected: []ExpressionToken{
				{
					Kind:  timeToken,
					Value: time.Date(2014, time.January, 2, 0, 0, 0, 0, time.Local),
				},
			},
		},
		{

			Name:  "Single time, RFC3339, with hh:mm",
			Input: "'2014-01-02 14:12'",
			Expected: []ExpressionToken{
				{
					Kind:  timeToken,
					Value: time.Date(2014, time.January, 2, 14, 12, 0, 0, time.Local),
				},
			},
		},
		{

			Name:  "Single time, RFC3339, with hh:mm:ss",
			Input: "'2014-01-02 14:12:22'",
			Expected: []ExpressionToken{
				{
					Kind:  timeToken,
					Value: time.Date(2014, time.January, 2, 14, 12, 22, 0, time.Local),
				},
			},
		},
		{

			Name:  "Single boolean",
			Input: "true",
			Expected: []ExpressionToken{
				{
					Kind:  boolean,
					Value: true,
				},
			},
		},
		{

			Name:  "Single large numeric",
			Input: "1234567890",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1234567890.0,
				},
			},
		},
		{

			Name:  "Single floating-point",
			Input: "0.5",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 0.5,
				},
			},
		},
		{

			Name:  "Single large floating point",
			Input: "3.14567471",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 3.14567471,
				},
			},
		},
		{

			Name:  "Single false boolean",
			Input: "false",
			Expected: []ExpressionToken{
				{
					Kind:  boolean,
					Value: false,
				},
			},
		},
		{
			Name:  "Single internationalized string",
			Input: "'ÆŦǽഈᚥஇคٸ'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "ÆŦǽഈᚥஇคٸ",
				},
			},
		},
		{
			Name:  "Single internationalized parameter",
			Input: "ÆŦǽഈᚥஇคٸ",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "ÆŦǽഈᚥஇคٸ",
				},
			},
		},
		{
			Name:      "Parameterless function",
			Input:     "foo()",
			Functions: map[string]ExpressionFunction{"foo": noop},
			Expected: []ExpressionToken{
				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},
				{
					Kind: clauseClose,
				},
			},
		},
		{
			Name:      "Single parameter function",
			Input:     "foo('bar')",
			Functions: map[string]ExpressionFunction{"foo": noop},
			Expected: []ExpressionToken{
				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},
				{
					Kind:  stringToken,
					Value: "bar",
				},
				{
					Kind: clauseClose,
				},
			},
		},
		{
			Name:      "Multiple parameter function",
			Input:     "foo('bar', 1.0)",
			Functions: map[string]ExpressionFunction{"foo": noop},
			Expected: []ExpressionToken{
				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},
				{
					Kind:  stringToken,
					Value: "bar",
				},
				{
					Kind: separator,
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind: clauseClose,
				},
			},
		},
		{
			Name:      "Nested function",
			Input:     "foo(foo('bar'), 1.0, foo(2.0))",
			Functions: map[string]ExpressionFunction{"foo": noop},
			Expected: []ExpressionToken{
				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},

				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},
				{
					Kind:  stringToken,
					Value: "bar",
				},
				{
					Kind: clauseClose,
				},
				{
					Kind: separator,
				},

				{
					Kind:  numeric,
					Value: 1.0,
				},

				{
					Kind: separator,
				},

				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},
				{
					Kind:  numeric,
					Value: 2.0,
				},
				{
					Kind: clauseClose,
				},

				{
					Kind: clauseClose,
				},
			},
		},
		{
			Name:      "Function with modifier afterwards (#28)",
			Input:     "foo() + 1",
			Functions: map[string]ExpressionFunction{"foo": noop},
			Expected: []ExpressionToken{
				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},
				{
					Kind: clauseClose,
				},
				{
					Kind:  modifier,
					Value: "+",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{
			Name:      "Function with modifier afterwards and comparator",
			Input:     "(foo()-1) > 3",
			Functions: map[string]ExpressionFunction{"foo": noop},
			Expected: []ExpressionToken{
				{
					Kind: clause,
				},
				{
					Kind:  function,
					Value: noop,
				},
				{
					Kind: clause,
				},
				{
					Kind: clauseClose,
				},
				{
					Kind:  modifier,
					Value: "-",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind: clauseClose,
				},
				{
					Kind:  comparator,
					Value: ">",
				},
				{
					Kind:  numeric,
					Value: 3.0,
				},
			},
		},
		{
			Name:  "Double-quoted string added to square-brackted param (#59)",
			Input: "\"a\" + [foo]",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "a",
				},
				{
					Kind:  modifier,
					Value: "+",
				},
				{
					Kind:  variable,
					Value: "foo",
				},
			},
		},
		{
			Name:  "Accessor variable",
			Input: "foo.Var",
			Expected: []ExpressionToken{
				{
					Kind:  accessor,
					Value: []string{"foo", "Var"},
				},
			},
		},
		{
			Name:  "Accessor function",
			Input: "foo.Operation()",
			Expected: []ExpressionToken{
				{
					Kind:  accessor,
					Value: []string{"foo", "Operation"},
				},
				{
					Kind: clause,
				},
				{
					Kind: clauseClose,
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestLogicalOperatorParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{

		{

			Name:  "Boolean AND",
			Input: "true && false",
			Expected: []ExpressionToken{
				{
					Kind:  boolean,
					Value: true,
				},
				{
					Kind:  logicalop,
					Value: "&&",
				},
				{
					Kind:  boolean,
					Value: false,
				},
			},
		},
		{

			Name:  "Boolean OR",
			Input: "true || false",
			Expected: []ExpressionToken{
				{
					Kind:  boolean,
					Value: true,
				},
				{
					Kind:  logicalop,
					Value: "||",
				},
				{
					Kind:  boolean,
					Value: false,
				},
			},
		},
		{

			Name:  "Multiple logical operators",
			Input: "true || false && true",
			Expected: []ExpressionToken{
				{
					Kind:  boolean,
					Value: true,
				},
				{
					Kind:  logicalop,
					Value: "||",
				},
				{
					Kind:  boolean,
					Value: false,
				},
				{
					Kind:  logicalop,
					Value: "&&",
				},
				{
					Kind:  boolean,
					Value: true,
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestComparatorParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{

		{

			Name:  "Numeric EQ",
			Input: "1 == 2",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  numeric,
					Value: 2.0,
				},
			},
		},
		{

			Name:  "Numeric NEQ",
			Input: "1 != 2",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  comparator,
					Value: "!=",
				},
				{
					Kind:  numeric,
					Value: 2.0,
				},
			},
		},
		{

			Name:  "Numeric GT",
			Input: "1 > 0",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  comparator,
					Value: ">",
				},
				{
					Kind:  numeric,
					Value: 0.0,
				},
			},
		},
		{

			Name:  "Numeric LT",
			Input: "1 < 2",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  comparator,
					Value: "<",
				},
				{
					Kind:  numeric,
					Value: 2.0,
				},
			},
		},
		{

			Name:  "Numeric GTE",
			Input: "1 >= 2",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  comparator,
					Value: ">=",
				},
				{
					Kind:  numeric,
					Value: 2.0,
				},
			},
		},
		{

			Name:  "Numeric LTE",
			Input: "1 <= 2",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  comparator,
					Value: "<=",
				},
				{
					Kind:  numeric,
					Value: 2.0,
				},
			},
		},
		{

			Name:  "String LT",
			Input: "'ab.cd' < 'abc.def'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "ab.cd",
				},
				{
					Kind:  comparator,
					Value: "<",
				},
				{
					Kind:  stringToken,
					Value: "abc.def",
				},
			},
		},
		{

			Name:  "String LTE",
			Input: "'ab.cd' <= 'abc.def'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "ab.cd",
				},
				{
					Kind:  comparator,
					Value: "<=",
				},
				{
					Kind:  stringToken,
					Value: "abc.def",
				},
			},
		},
		{

			Name:  "String GT",
			Input: "'ab.cd' > 'abc.def'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "ab.cd",
				},
				{
					Kind:  comparator,
					Value: ">",
				},
				{
					Kind:  stringToken,
					Value: "abc.def",
				},
			},
		},
		{

			Name:  "String GTE",
			Input: "'ab.cd' >= 'abc.def'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "ab.cd",
				},
				{
					Kind:  comparator,
					Value: ">=",
				},
				{
					Kind:  stringToken,
					Value: "abc.def",
				},
			},
		},
		{

			Name:  "String REQ",
			Input: "'foobar' =~ 'bar'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foobar",
				},
				{
					Kind:  comparator,
					Value: "=~",
				},

				// it's not particularly clean to test for the contents of a pattern, (since it means modifying the harness below)
				// so pattern contents are left untested.
				{
					Kind: pattern,
				},
			},
		},
		{

			Name:  "String NREQ",
			Input: "'foobar' !~ 'bar'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foobar",
				},
				{
					Kind:  comparator,
					Value: "!~",
				},
				{
					Kind: pattern,
				},
			},
		},
		{

			Name:  "Comparator against modifier string additive (#22)",
			Input: "'foo' == '+'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  stringToken,
					Value: "+",
				},
			},
		},
		{

			Name:  "Comparator against modifier string multiplicative (#22)",
			Input: "'foo' == '/'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  stringToken,
					Value: "/",
				},
			},
		},
		{

			Name:  "Comparator against modifier string exponential (#22)",
			Input: "'foo' == '**'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  stringToken,
					Value: "**",
				},
			},
		},
		{

			Name:  "Comparator against modifier string bitwise (#22)",
			Input: "'foo' == '^'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  stringToken,
					Value: "^",
				},
			},
		},
		{

			Name:  "Comparator against modifier string shift (#22)",
			Input: "'foo' == '>>'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  stringToken,
					Value: ">>",
				},
			},
		},
		{

			Name:  "Comparator against modifier string ternary (#22)",
			Input: "'foo' == '?'",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  stringToken,
					Value: "?",
				},
			},
		},
		{

			Name:  "Array membership lowercase",
			Input: "'foo' in ('foo', 'bar')",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "in",
				},
				{
					Kind: clause,
				},
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind: separator,
				},
				{
					Kind:  stringToken,
					Value: "bar",
				},
				{
					Kind: clauseClose,
				},
			},
		},
		{

			Name:  "Array membership uppercase",
			Input: "'foo' IN ('foo', 'bar')",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: "in",
				},
				{
					Kind: clause,
				},
				{
					Kind:  stringToken,
					Value: "foo",
				},
				{
					Kind: separator,
				},
				{
					Kind:  stringToken,
					Value: "bar",
				},
				{
					Kind: clauseClose,
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestModifierParsing(test *testing.T) {

	tokenParsingTests := []TokenParsingTest{

		{

			Name:  "Numeric PLUS",
			Input: "1 + 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "+",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric MINUS",
			Input: "1 - 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "-",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric MULTIPLY",
			Input: "1 * 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "*",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric DIVIDE",
			Input: "1 / 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "/",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric MODULUS",
			Input: "1 % 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "%",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric BITWISE_AND",
			Input: "1 & 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "&",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric BITWISE_OR",
			Input: "1 | 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "|",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric BITWISE_XOR",
			Input: "1 ^ 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "^",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric BITWISE_LSHIFT",
			Input: "1 << 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: "<<",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Numeric BITWISE_RSHIFT",
			Input: "1 >> 1",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  modifier,
					Value: ">>",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
	}

	tokenParsingTests = combineWhitespaceExpressions(tokenParsingTests)
	runTokenParsingTest(tokenParsingTests, test)
}

func TestPrefixParsing(test *testing.T) {

	testCases := []TokenParsingTest{

		{

			Name:  "Sign prefix",
			Input: "-1",
			Expected: []ExpressionToken{
				{
					Kind:  prefix,
					Value: "-",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Sign prefix on variable",
			Input: "-foo",
			Expected: []ExpressionToken{
				{
					Kind:  prefix,
					Value: "-",
				},
				{
					Kind:  variable,
					Value: "foo",
				},
			},
		},
		{

			Name:  "Boolean prefix",
			Input: "!true",
			Expected: []ExpressionToken{
				{
					Kind:  prefix,
					Value: "!",
				},
				{
					Kind:  boolean,
					Value: true,
				},
			},
		},
		{

			Name:  "Boolean prefix on variable",
			Input: "!foo",
			Expected: []ExpressionToken{
				{
					Kind:  prefix,
					Value: "!",
				},
				{
					Kind:  variable,
					Value: "foo",
				},
			},
		},
		{

			Name:  "Bitwise not prefix",
			Input: "~1",
			Expected: []ExpressionToken{
				{
					Kind:  prefix,
					Value: "~",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Bitwise not prefix on variable",
			Input: "~foo",
			Expected: []ExpressionToken{
				{
					Kind:  prefix,
					Value: "~",
				},
				{
					Kind:  variable,
					Value: "foo",
				},
			},
		},
	}

	testCases = combineWhitespaceExpressions(testCases)
	runTokenParsingTest(testCases, test)
}

func TestEscapedParameters(test *testing.T) {

	testCases := []TokenParsingTest{

		{

			Name:  "Single escaped parameter",
			Input: "[foo]",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "foo",
				},
			},
		},
		{

			Name:  "Single escaped parameter with whitespace",
			Input: "[foo bar]",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "foo bar",
				},
			},
		},
		{

			Name:  "Single escaped parameter with escaped closing bracket",
			Input: "[foo[bar\\]]",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "foo[bar]",
				},
			},
		},
		{

			Name:  "Escaped parameters and unescaped parameters",
			Input: "[foo] > bar",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "foo",
				},
				{
					Kind:  comparator,
					Value: ">",
				},
				{
					Kind:  variable,
					Value: "bar",
				},
			},
		},
		{

			Name:  "Unescaped parameter with space",
			Input: "foo\\ bar > bar",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "foo bar",
				},
				{
					Kind:  comparator,
					Value: ">",
				},
				{
					Kind:  variable,
					Value: "bar",
				},
			},
		},
		{

			Name:  "Unescaped parameter with space",
			Input: "response\\-time > bar",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "response-time",
				},
				{
					Kind:  comparator,
					Value: ">",
				},
				{
					Kind:  variable,
					Value: "bar",
				},
			},
		},
		{

			Name:  "Parameters with snake_case",
			Input: "foo_bar > baz_quux",
			Expected: []ExpressionToken{
				{
					Kind:  variable,
					Value: "foo_bar",
				},
				{
					Kind:  comparator,
					Value: ">",
				},
				{
					Kind:  variable,
					Value: "baz_quux",
				},
			},
		},
		{

			Name:  "String literal uses backslash to escape",
			Input: "\"foo'bar\"",
			Expected: []ExpressionToken{
				{
					Kind:  stringToken,
					Value: "foo'bar",
				},
			},
		},
	}

	runTokenParsingTest(testCases, test)
}

func TestTernaryParsing(test *testing.T) {
	tokenParsingTests := []TokenParsingTest{

		{

			Name:  "Ternary after Boolean",
			Input: "true ? 1",
			Expected: []ExpressionToken{
				{
					Kind:  boolean,
					Value: true,
				},
				{
					Kind:  ternary,
					Value: "?",
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
		},
		{

			Name:  "Ternary after Comperator",
			Input: "1 == 0 ? true",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  comparator,
					Value: "==",
				},
				{
					Kind:  numeric,
					Value: 0.0,
				},
				{
					Kind:  ternary,
					Value: "?",
				},
				{
					Kind:  boolean,
					Value: true,
				},
			},
		},
		{

			Name:  "Null coalesce left",
			Input: "1 ?? 2",
			Expected: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind:  ternary,
					Value: "??",
				},
				{
					Kind:  numeric,
					Value: 2.0,
				},
			},
		},
	}

	runTokenParsingTest(tokenParsingTests, test)
}

// Tests to make sure that the String() reprsentation of an expression exactly matches what is given to the parse function.
func TestOriginalString(test *testing.T) {

	// include all the token types, to be sure there's no shenaniganery going on.
	expressionString := "2 > 1 &&" +
		"'something' != 'nothing' || " +
		"'2014-01-20' < 'Wed Jul  8 23:07:35 MDT 2015' && " +
		"[escapedVariable name with spaces] <= unescaped\\-variableName &&" +
		"modifierTest + 1000 / 2 > (80 * 100 % 2) && true ? true : false"

	expression, err := NewExpression(expressionString)
	if err != nil {

		test.Logf("failed to parse original string test: %v", err)
		test.Fail()
		return
	}

	if expression.String() != expressionString {
		test.Logf("String() did not give the same expression as given to parse")
		test.Fail()
	}
}

// Tests to make sure that the Vars() reprsentation of an expression identifies all variables contained within the expression.
func TestOriginalVars(test *testing.T) {

	// include all the token types, to be sure there's no shenaniganery going on.
	expressionString := "2 > 1 &&" +
		"'something' != 'nothing' || " +
		"'2014-01-20' < 'Wed Jul  8 23:07:35 MDT 2015' && " +
		"[escapedVariable name with spaces] <= unescaped\\-variableName &&" +
		"modifierTest + 1000 / 2 > (80 * 100 % 2) && true ? true : false"

	expectedVars := [3]string{"escapedVariable name with spaces",
		"modifierTest",
		"unescaped-variableName"}

	expression, err := NewExpression(expressionString)
	if err != nil {

		test.Logf("failed to parse original var test: %v", err)
		test.Fail()
		return
	}

	if len(expression.Vars()) == len(expectedVars) {
		variableMap := make(map[string]string)
		for _, v := range expression.Vars() {
			variableMap[v] = v
		}
		for _, v := range expectedVars {
			if _, ok := variableMap[v]; !ok {
				test.Logf("Vars() did not correctly identify all variables contained within the expression")
				test.Fail()
			}
		}
	} else {
		test.Logf("Vars() did not correctly identify all variables contained within the expression")
		test.Fail()
	}
}

func combineWhitespaceExpressions(testCases []TokenParsingTest) []TokenParsingTest {

	var currentCase, strippedCase TokenParsingTest
	var caseLength int

	caseLength = len(testCases)

	for i := 0; i < caseLength; i++ {

		currentCase = testCases[i]

		strippedCase = TokenParsingTest{

			Name:      (currentCase.Name + " (without whitespace)"),
			Input:     stripUnquotedWhitespace(currentCase.Input),
			Expected:  currentCase.Expected,
			Functions: currentCase.Functions,
		}

		testCases = append(testCases, strippedCase, currentCase)
	}

	return testCases
}

func stripUnquotedWhitespace(expression string) string {

	var expressionBuffer bytes.Buffer
	var quoted bool

	for _, character := range expression {

		if !quoted && unicode.IsSpace(character) {
			continue
		}

		if character == '\'' {
			quoted = !quoted
		}

		expressionBuffer.WriteString(string(character))
	}

	return expressionBuffer.String()
}

func runTokenParsingTest(tokenParsingTests []TokenParsingTest, test *testing.T) {

	var parsingTest TokenParsingTest
	var expression *Expression
	var actualTokens []ExpressionToken
	var actualToken ExpressionToken
	var expectedTokenKindString, actualTokenKindString string
	var expectedTokenLength, actualTokenLength int
	var err error

	fmt.Printf("Running %d parsing test cases...\n", len(tokenParsingTests))
	// defer func() {
	//     if r := recover(); r != nil {
	//         test.Logf("Panic in test '%s': %v", parsingTest.Name, r)
	// 		test.Fail()
	//     }
	// }()

	// Run the test cases.
	for _, parsingTest = range tokenParsingTests {

		if parsingTest.Functions != nil {
			expression, err = NewExpressionWithFunctions(parsingTest.Input, parsingTest.Functions)
		} else {
			expression, err = NewExpression(parsingTest.Input)
		}

		if err != nil {

			test.Logf("Test '%s' failed to parse: %s", parsingTest.Name, err)
			test.Logf("Expression: '%s'", parsingTest.Input)
			test.Fail()
			continue
		}

		actualTokens = expression.Tokens()

		expectedTokenLength = len(parsingTest.Expected)
		actualTokenLength = len(actualTokens)

		if actualTokenLength != expectedTokenLength {

			test.Logf("Test '%s' failed:", parsingTest.Name)
			test.Logf("Expected %d tokens, actually found %d", expectedTokenLength, actualTokenLength)
			test.Fail()
			continue
		}

		for i, expectedToken := range parsingTest.Expected {

			actualToken = actualTokens[i]
			if actualToken.Kind != expectedToken.Kind {

				actualTokenKindString = actualToken.Kind.String()
				expectedTokenKindString = expectedToken.Kind.String()

				test.Logf("Test '%s' failed:", parsingTest.Name)
				test.Logf("Expected token kind '%v' does not match '%v'", expectedTokenKindString, actualTokenKindString)
				test.Fail()
				continue
			}

			if expectedToken.Value == nil {
				continue
			}

			reflectedKind := reflect.TypeOf(expectedToken.Value).Kind()
			if reflectedKind == reflect.Func {
				continue
			}

			// gotta be an accessor
			if reflectedKind == reflect.Slice {

				if actualToken.Value == nil {
					test.Logf("Test '%s' failed:", parsingTest.Name)
					test.Logf("Expected token value '%v' does not match nil", expectedToken.Value)
					test.Fail()
				}

				for z, actual := range actualToken.Value.([]string) {

					if actual != expectedToken.Value.([]string)[z] {

						test.Logf("Test '%s' failed:", parsingTest.Name)
						test.Logf("Expected token value '%v' does not match '%v'", expectedToken.Value, actualToken.Value)
						test.Fail()
					}
				}
				continue
			}

			if actualToken.Value != expectedToken.Value {

				test.Logf("Test '%s' failed:", parsingTest.Name)
				test.Logf("Expected token value '%v' does not match '%v'", expectedToken.Value, actualToken.Value)
				test.Fail()
				continue
			}
		}
	}
}

func noop(arguments ...interface{}) (interface{}, error) {
	return nil, nil
}
