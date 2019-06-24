package govaluate

// Tests to make sure evaluation fails in the expected ways.
import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type DebugStruct struct {
	x int
}

// Represents a test for parsing failures
type EvaluationFailureTest struct {
	Name       string
	Input      string
	Functions  map[string]ExpressionFunction
	Parameters map[string]interface{}
	Expected   string
}

const (
	invalidModifierTypes   string = "cannot be used with the modifier"
	invalidComparatorTypes        = "cannot be used with the comparator"
	invalidLogicalopTypes         = "cannot be used with the logical operator"
	invalidTernaryTypes           = "cannot be used with the ternary operator"
	invalidRegex                  = "Unable to compile regexp pattern"
	invalidParameterCall          = "No method or field"
	tooFewArgs                    = "Too few arguments to parameter call"
	tooManyArgs                   = "Too many arguments to parameter call"
	mismatchedParameters          = "Argument type conversion failed"
)

// preset parameter map of types that can be used in an evaluation failure test to check typing.
var evaluationFailureParameters = map[string]interface{}{
	"number": 1,
	"string": "foo",
	"bool":   true,
}

func TestComplexParameter(test *testing.T) {

	var expression *Expression
	var err error
	var v interface{}

	parameters := map[string]interface{}{
		"complex64":  complex64(0),
		"complex128": complex128(0),
	}

	expression, _ = NewExpression("complex64")
	v, err = expression.Evaluate(parameters)
	if err != nil {
		test.Errorf("Expected no error, but have %s", err)
	}
	if v.(complex64) != complex64(0) {
		test.Errorf("Expected %v == %v", v, complex64(0))
	}

	expression, _ = NewExpression("complex128")
	v, err = expression.Evaluate(parameters)
	if err != nil {
		test.Errorf("Expected no error, but have %s", err)
	}
	if v.(complex128) != complex128(0) {
		test.Errorf("Expected %v == %v", v, complex128(0))
	}
}

func TestStructParameter(t *testing.T) {
	expected := DebugStruct{}
	expression, _ := NewExpression("foo")
	parameters := map[string]interface{}{"foo": expected}
	v, err := expression.Evaluate(parameters)
	if err != nil {
		t.Errorf("Expected no error, but have %s", err)
	} else if v.(DebugStruct) != expected {
		t.Errorf("Values mismatch: %v != %v", expected, v)
	}
}

func TestNilParameterUsage(test *testing.T) {

	expression, err := NewExpression("2 > 1")
	_, err = expression.Evaluate(nil)

	if err != nil {
		test.Errorf("Expected no error from nil parameter evaluation, got %v\n", err)
		return
	}
}

func TestModifierTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:     "PLUS literal number to literal bool",
			Input:    "1 + true",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "PLUS number to bool",
			Input:    "number + bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "MINUS number to bool",
			Input:    "number - bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "MINUS number to bool",
			Input:    "number - bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "MULTIPLY number to bool",
			Input:    "number * bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "DIVIDE number to bool",
			Input:    "number / bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "EXPONENT number to bool",
			Input:    "number ** bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "MODULUS number to bool",
			Input:    "number % bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "XOR number to bool",
			Input:    "number % bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "BITWISE_OR number to bool",
			Input:    "number | bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "BITWISE_AND number to bool",
			Input:    "number & bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "BITWISE_XOR number to bool",
			Input:    "number ^ bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "BITWISE_LSHIFT number to bool",
			Input:    "number << bool",
			Expected: invalidModifierTypes,
		},
		EvaluationFailureTest{
			Name:     "BITWISE_RSHIFT number to bool",
			Input:    "number >> bool",
			Expected: invalidModifierTypes,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestLogicalOperatorTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:     "AND number to number",
			Input:    "number && number",
			Expected: invalidLogicalopTypes,
		},
		EvaluationFailureTest{
			Name:     "OR number to number",
			Input:    "number || number",
			Expected: invalidLogicalopTypes,
		},
		EvaluationFailureTest{
			Name:     "AND string to string",
			Input:    "string && string",
			Expected: invalidLogicalopTypes,
		},
		EvaluationFailureTest{
			Name:     "OR string to string",
			Input:    "string || string",
			Expected: invalidLogicalopTypes,
		},
		EvaluationFailureTest{
			Name:     "AND number to string",
			Input:    "number && string",
			Expected: invalidLogicalopTypes,
		},
		EvaluationFailureTest{
			Name:     "OR number to string",
			Input:    "number || string",
			Expected: invalidLogicalopTypes,
		},
		EvaluationFailureTest{
			Name:     "AND bool to string",
			Input:    "bool && string",
			Expected: invalidLogicalopTypes,
		},
		EvaluationFailureTest{
			Name:     "OR string to bool",
			Input:    "string || bool",
			Expected: invalidLogicalopTypes,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

// While there is type-safe transitions checked at parse-time, tested in the "parsing_test" and "parsingFailure_test" files,
// we also need to make sure that we receive type mismatch errors during evaluation.
func TestComparatorTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:     "GT literal bool to literal bool",
			Input:    "true > true",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "GT bool to bool",
			Input:    "bool > bool",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "GTE bool to bool",
			Input:    "bool >= bool",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "LT bool to bool",
			Input:    "bool < bool",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "LTE bool to bool",
			Input:    "bool <= bool",
			Expected: invalidComparatorTypes,
		},

		EvaluationFailureTest{
			Name:     "GT number to string",
			Input:    "number > string",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "GTE number to string",
			Input:    "number >= string",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "LT number to string",
			Input:    "number < string",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "REQ number to string",
			Input:    "number =~ string",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "REQ number to bool",
			Input:    "number =~ bool",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "REQ bool to number",
			Input:    "bool =~ number",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "REQ bool to string",
			Input:    "bool =~ string",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "NREQ number to string",
			Input:    "number !~ string",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "NREQ number to bool",
			Input:    "number !~ bool",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "NREQ bool to number",
			Input:    "bool !~ number",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "NREQ bool to string",
			Input:    "bool !~ string",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "IN non-array numeric",
			Input:    "1 in 2",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "IN non-array string",
			Input:    "1 in 'foo'",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "IN non-array boolean",
			Input:    "1 in true",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "EQ number to boolean",
			Input:    "1 == true",
			Expected: invalidComparatorTypes,
		},
		EvaluationFailureTest{
			Name:     "NEQ string to number",
			Input:    "'hello' != 10",
			Expected: invalidComparatorTypes,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestTernaryTyping(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:     "Ternary with number",
			Input:    "10 ? true",
			Expected: invalidTernaryTypes,
		},
		EvaluationFailureTest{
			Name:     "Ternary with string",
			Input:    "'foo' ? true",
			Expected: invalidTernaryTypes,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestRegexParameterCompilation(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:  "Regex equality runtime parsing",
			Input: "'foo' =~ foo",
			Parameters: map[string]interface{}{
				"foo": "[foo",
			},
			Expected: invalidRegex,
		},
		EvaluationFailureTest{
			Name:  "Regex inequality runtime parsing",
			Input: "'foo' =~ foo",
			Parameters: map[string]interface{}{
				"foo": "[foo",
			},
			Expected: invalidRegex,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestFunctionExecution(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:  "Function error bubbling",
			Input: "error()",
			Functions: map[string]ExpressionFunction{
				"error": func(arguments ...interface{}) (interface{}, error) {
					return nil, errors.New("Huge problems")
				},
			},
			Expected: "Huge problems",
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func TestInvalidParameterCalls(test *testing.T) {

	evaluationTests := []EvaluationFailureTest{
		EvaluationFailureTest{
			Name:       "Missing parameter field reference",
			Input:      "foo.NotExists",
			Parameters: fooFailureParameters,
			Expected:   invalidParameterCall,
		},
		EvaluationFailureTest{
			Name:       "Parameter method call on missing function",
			Input:      "foo.NotExist()",
			Parameters: fooFailureParameters,
			Expected:   invalidParameterCall,
		},
		EvaluationFailureTest{
			Name:       "Nested missing parameter field reference",
			Input:      "foo.Nested.NotExists",
			Parameters: fooFailureParameters,
			Expected:   invalidParameterCall,
		},
		EvaluationFailureTest{
			Name:       "Parameter method call returns error",
			Input:      "foo.AlwaysFail()",
			Parameters: fooFailureParameters,
			Expected:   "function should always fail",
		},
		EvaluationFailureTest{
			Name:       "Too few arguments to parameter call",
			Input:      "foo.FuncArgStr()",
			Parameters: fooFailureParameters,
			Expected:   tooFewArgs,
		},
		EvaluationFailureTest{
			Name:       "Too many arguments to parameter call",
			Input:      "foo.FuncArgStr('foo', 'bar', 15)",
			Parameters: fooFailureParameters,
			Expected:   tooManyArgs,
		},
		EvaluationFailureTest{
			Name:       "Mismatched parameters",
			Input:      "foo.FuncArgStr(5)",
			Parameters: fooFailureParameters,
			Expected:   mismatchedParameters,
		},
	}

	runEvaluationFailureTests(evaluationTests, test)
}

func runEvaluationFailureTests(evaluationTests []EvaluationFailureTest, test *testing.T) {

	var expression *Expression
	var err error

	fmt.Printf("Running %d negative parsing test cases...\n", len(evaluationTests))

	for _, testCase := range evaluationTests {

		if len(testCase.Functions) > 0 {
			expression, err = NewExpressionWithFunctions(testCase.Input, testCase.Functions)
		} else {
			expression, err = NewExpression(testCase.Input)
		}

		if err != nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected evaluation error, but got parsing error: '%s'", err)
			test.Fail()
			continue
		}

		if testCase.Parameters == nil {
			testCase.Parameters = evaluationFailureParameters
		}

		_, err = expression.Evaluate(testCase.Parameters)

		if err == nil {

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Expected error, received none.")
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
