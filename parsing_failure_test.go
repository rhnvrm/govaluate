package govaluate

import (
	"fmt"
	"regexp/syntax"
	"strings"
	"testing"
)

const (
	unexpectedEnd          string = "Unexpected end of expression"
	invalidTokenTransition        = "Cannot transition token types"
	invalidTokenKind              = "Invalid token"
	unclosedQuotes                = "Unclosed string literal"
	unclosedBrackets              = "Unclosed parameter bracket"
	unbalancedParenthesis         = "Unbalanced parenthesis"
	invalidNumeric                = "Unable to parse numeric value"
	undefinedFunction             = "Undefined function"
	hangingAccessor               = "Hanging accessor on token"
	unexportedAccessor            = "Unable to access unexported"
	invalidHex                    = "Unable to parse hex value"
)

// Represents a test for parsing failures
type ParsingFailureTest struct {
	Name     string
	Input    string
	Expected string
}

func TestParsingFailure(test *testing.T) {

	parsingTests := []ParsingFailureTest{
		{
			Name:     "Invalid equality comparator",
			Input:    "1 = 1",
			Expected: invalidTokenKind,
		},
		{
			Name:     "Invalid equality comparator",
			Input:    "1 === 1",
			Expected: invalidTokenKind,
		},
		{
			Name:     "Too many characters for logical operator",
			Input:    "true &&& false",
			Expected: invalidTokenKind,
		},
		{
			Name:     "Too many characters for logical operator",
			Input:    "true ||| false",
			Expected: invalidTokenKind,
		},
		{
			Name:     "Premature end to expression, via modifier",
			Input:    "10 > 5 +",
			Expected: unexpectedEnd,
		},
		{
			Name:     "Premature end to expression, via comparator",
			Input:    "10 + 5 >",
			Expected: unexpectedEnd,
		},
		{
			Name:     "Premature end to expression, via logical operator",
			Input:    "10 > 5 &&",
			Expected: unexpectedEnd,
		},
		{
			Name:     "Premature end to expression, via ternary operator",
			Input:    "true ?",
			Expected: unexpectedEnd,
		},
		{
			Name:     "Hanging REQ",
			Input:    "'wat' =~",
			Expected: unexpectedEnd,
		},
		{
			Name:     "Invalid operator change to REQ",
			Input:    " / =~",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid starting token, comparator",
			Input:    "> 10",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid starting token, modifier",
			Input:    "+ 5",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid starting token, logical operator",
			Input:    "&& 5 < 10",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid NUMERIC transition",
			Input:    "10 10",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid STRING transition",
			Input:    "'foo' 'foo'",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid operator transition",
			Input:    "10 > < 10",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Starting with unbalanced parens",
			Input:    " ) ( arg2",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Unclosed bracket",
			Input:    "[foo bar",
			Expected: unclosedBrackets,
		},
		{
			Name:     "Unclosed quote",
			Input:    "foo == 'responseTime",
			Expected: unclosedQuotes,
		},
		{
			Name:     "Constant regex pattern fail to compile",
			Input:    "foo =~ '[abc'",
			Expected: string(syntax.ErrMissingBracket),
		},
		{
			Name:     "Unbalanced parenthesis",
			Input:    "10 > (1 + 50",
			Expected: unbalancedParenthesis,
		},
		{
			Name:     "Multiple radix",
			Input:    "127.0.0.1",
			Expected: invalidNumeric,
		},
		{
			Name:     "Undefined function",
			Input:    "foobar()",
			Expected: undefinedFunction,
		},
		{
			Name:     "Hanging accessor",
			Input:    "foo.Bar.",
			Expected: hangingAccessor,
		},
		{
			// this is expected to change once there are structtags in place that allow aliasing of fields
			Name:     "Unexported parameter access",
			Input:    "foo.bar",
			Expected: unexportedAccessor,
		},
		{
			Name:     "Incomplete Hex",
			Input:    "0x",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid Hex literal",
			Input:    "0x > 0",
			Expected: invalidHex,
		},
		{
			Name:     "Hex float (Unsupported)",
			Input:    "0x1.1",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Hex invalid letter",
			Input:    "0x12g1",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid LOGICALOP transition",
			Input:    "(a > 100 &&) == false",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid MODIFIER transition",
			Input:    "(a + )",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid COMPARATOR transition",
			Input:    "(a > )",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid PREFIX transition",
			Input:    "(~)",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Invalid CLAUSE_CLOSE transition",
			Input:    "(a == b) c",
			Expected: invalidTokenTransition,
		},
		{
			Name:     "Hanging logical operation, followed by clause-close (#92)",
			Input:    "(amount > '100' &&) == false",
			Expected: invalidTokenTransition,
		},
	}

	runParsingFailureTests(parsingTests, test)
}

func runParsingFailureTests(parsingTests []ParsingFailureTest, test *testing.T) {

	var err error

	fmt.Printf("Running %d parsing test cases...\n", len(parsingTests))

	for _, testCase := range parsingTests {

		_, err = NewExpression(testCase.Input)

		if err == nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected a parsing error, found no error.")
			test.Fail()
			continue
		}

		if !strings.Contains(err.Error(), testCase.Expected) {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Got error: '%s', expected '%s'", err.Error(), testCase.Expected)
			test.Fail()
			continue
		}
	}
}
