package govaluate

import (
	"fmt"
	"strings"
	"testing"
)

const (
	expErrNilValue string = "cannot have a nil value"
)

// Contains a single test case for the Expression.NewExpressionFromTokens() method.
//
// These tests, and the ones in `tokenExpressionFailure_test` will be fairly incomplete.
// Creating an expression from a string and from tokens _must_ both perform the same syntax checks.
// So all the checks in `parsing_test` will follow the same logic as the ones here.
//
// These tests check some corner cases - such as tokens having nil values when they must have something.
// Cases that cannot occur through the normal parser, but may occur in other parsers.
type ExpressionTokenSyntaxTest struct {
	Name     string
	Input    []ExpressionToken
	Expected string
}

func TestNilValues(test *testing.T) {

	cases := []ExpressionTokenSyntaxTest{
		{
			Name: "Nil numeric",
			Input: []ExpressionToken{
				{
					Kind: numeric,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil string",
			Input: []ExpressionToken{
				{
					Kind: stringToken,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil bool",
			Input: []ExpressionToken{
				{
					Kind: boolean,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil time",
			Input: []ExpressionToken{
				{
					Kind: timeToken,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil pattern",
			Input: []ExpressionToken{
				{
					Kind: pattern,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil variable",
			Input: []ExpressionToken{
				{
					Kind: variable,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil prefix",
			Input: []ExpressionToken{
				{
					Kind: prefix,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil comparator",
			Input: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind: comparator,
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil logicalop",
			Input: []ExpressionToken{
				{
					Kind:  boolean,
					Value: true,
				},
				{
					Kind: logicalop,
				},
				{
					Kind:  boolean,
					Value: true,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil modifer",
			Input: []ExpressionToken{
				{
					Kind:  numeric,
					Value: 1.0,
				},
				{
					Kind: modifier,
				},
				{
					Kind:  numeric,
					Value: 1.0,
				},
			},
			Expected: expErrNilValue,
		},
		{
			Name: "Nil ternary",
			Input: []ExpressionToken{
				{
					Kind:  boolean,
					Value: true,
				},
				{
					Kind: ternary,
				},
				{
					Kind:  boolean,
					Value: true,
				},
			},
			Expected: expErrNilValue,
		},
	}

	runExpressionFromTokenTests(cases, true, test)
}

func runExpressionFromTokenTests(cases []ExpressionTokenSyntaxTest, expectFail bool, test *testing.T) {

	var err error

	fmt.Printf("Running %d expression from expression token tests...\n", len(cases))

	for _, testCase := range cases {

		_, err = NewExpressionFromTokens(testCase.Input)

		if err != nil {
			if expectFail {

				if !strings.Contains(err.Error(), testCase.Expected) {

					test.Logf("Test '%s' failed", testCase.Name)
					test.Logf("Got error: '%s', expected '%s'", err.Error(), testCase.Expected)
					test.Fail()
				}
				continue
			}

			test.Logf("Test '%s' failed", testCase.Name)
			test.Logf("Got error: '%s'", err)
			test.Fail()
			continue
		} else {
			if expectFail {

				test.Logf("Test '%s' failed", testCase.Name)
				test.Logf("Expected error, found none\n")
				test.Fail()
				continue
			}
		}
	}
}
