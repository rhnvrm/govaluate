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
		ExpressionTokenSyntaxTest{
			Name: "Nil numeric",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind: numeric,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil string",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind: stringToken,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil bool",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind: boolean,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil time",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind: timeToken,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil pattern",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind: pattern,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil variable",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind: variable,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil prefix",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind: prefix,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil comparator",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind:  numeric,
					Value: 1.0,
				},
				ExpressionToken{
					Kind: comparator,
				},
				ExpressionToken{
					Kind:  numeric,
					Value: 1.0,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil logicalop",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind:  boolean,
					Value: true,
				},
				ExpressionToken{
					Kind: logicalop,
				},
				ExpressionToken{
					Kind:  boolean,
					Value: true,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil modifer",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind:  numeric,
					Value: 1.0,
				},
				ExpressionToken{
					Kind: modifier,
				},
				ExpressionToken{
					Kind:  numeric,
					Value: 1.0,
				},
			},
			Expected: expErrNilValue,
		},
		ExpressionTokenSyntaxTest{
			Name: "Nil ternary",
			Input: []ExpressionToken{
				ExpressionToken{
					Kind:  boolean,
					Value: true,
				},
				ExpressionToken{
					Kind: ternary,
				},
				ExpressionToken{
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
