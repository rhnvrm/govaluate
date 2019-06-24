package govaluate

import (
	"errors"
	"fmt"
)

type lexerState struct {
	isEOF          bool
	isNullable     bool
	kind           TokenKind
	validNextKinds []TokenKind
}

// lexer states.
// Constant for all purposes except compiler.
var validLexerStates = []lexerState{

	lexerState{
		kind:       unknown,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []TokenKind{

			prefix,
			numeric,
			boolean,
			variable,
			pattern,
			function,
			accessor,
			stringToken,
			timeToken,
			clause,
		},
	},

	lexerState{

		kind:       clause,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []TokenKind{

			prefix,
			numeric,
			boolean,
			variable,
			pattern,
			function,
			accessor,
			stringToken,
			timeToken,
			clause,
			clauseClose,
		},
	},

	lexerState{

		kind:       clauseClose,
		isEOF:      true,
		isNullable: true,
		validNextKinds: []TokenKind{

			comparator,
			modifier,
			numeric,
			boolean,
			variable,
			stringToken,
			pattern,
			timeToken,
			clause,
			clauseClose,
			logicalop,
			ternary,
			separator,
		},
	},

	lexerState{

		kind:       numeric,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			modifier,
			comparator,
			logicalop,
			clauseClose,
			ternary,
			separator,
		},
	},
	lexerState{

		kind:       boolean,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			modifier,
			comparator,
			logicalop,
			clauseClose,
			ternary,
			separator,
		},
	},
	lexerState{

		kind:       stringToken,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			modifier,
			comparator,
			logicalop,
			clauseClose,
			ternary,
			separator,
		},
	},
	lexerState{

		kind:       timeToken,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			modifier,
			comparator,
			logicalop,
			clauseClose,
			separator,
		},
	},
	lexerState{

		kind:       pattern,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			modifier,
			comparator,
			logicalop,
			clauseClose,
			separator,
		},
	},
	lexerState{

		kind:       variable,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{

			modifier,
			comparator,
			logicalop,
			clauseClose,
			ternary,
			separator,
		},
	},
	lexerState{

		kind:       modifier,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			prefix,
			numeric,
			variable,
			function,
			accessor,
			stringToken,
			boolean,
			clause,
			clauseClose,
		},
	},
	lexerState{

		kind:       comparator,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			prefix,
			numeric,
			boolean,
			variable,
			function,
			accessor,
			stringToken,
			timeToken,
			clause,
			clauseClose,
			pattern,
		},
	},
	lexerState{

		kind:       logicalop,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			prefix,
			numeric,
			boolean,
			variable,
			function,
			accessor,
			stringToken,
			timeToken,
			clause,
			clauseClose,
		},
	},
	lexerState{

		kind:       prefix,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			numeric,
			boolean,
			variable,
			function,
			accessor,
			clause,
			clauseClose,
		},
	},

	lexerState{

		kind:       ternary,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{

			prefix,
			numeric,
			boolean,
			stringToken,
			timeToken,
			variable,
			function,
			accessor,
			clause,
			separator,
		},
	},
	lexerState{

		kind:       function,
		isEOF:      false,
		isNullable: false,
		validNextKinds: []TokenKind{
			clause,
		},
	},
	lexerState{

		kind:       accessor,
		isEOF:      true,
		isNullable: false,
		validNextKinds: []TokenKind{
			clause,
			modifier,
			comparator,
			logicalop,
			clauseClose,
			ternary,
			separator,
		},
	},
	lexerState{

		kind:       separator,
		isEOF:      false,
		isNullable: true,
		validNextKinds: []TokenKind{

			prefix,
			numeric,
			boolean,
			stringToken,
			timeToken,
			variable,
			function,
			accessor,
			clause,
		},
	},
}

func (ls lexerState) canTransitionTo(kind TokenKind) bool {

	for _, validKind := range ls.validNextKinds {

		if validKind == kind {
			return true
		}
	}

	return false
}

func checkExpressionSyntax(tokens []ExpressionToken) error {

	var state lexerState
	var lastToken ExpressionToken
	var err error

	state = validLexerStates[0]

	for _, token := range tokens {

		if !state.canTransitionTo(token.Kind) {

			// call out a specific error for tokens looking like they want to be functions.
			if lastToken.Kind == variable && token.Kind == clause {
				return errors.New("Undefined function " + lastToken.Value.(string))
			}

			firstStateName := fmt.Sprintf("%s [%v]", state.kind.String(), lastToken.Value)
			nextStateName := fmt.Sprintf("%s [%v]", token.Kind.String(), token.Value)

			return errors.New("Cannot transition token types from " + firstStateName + " to " + nextStateName)
		}

		state, err = getLexerStateForToken(token.Kind)
		if err != nil {
			return err
		}

		if !state.isNullable && token.Value == nil {

			errorMsg := fmt.Sprintf("Token kind '%v' cannot have a nil value", token.Kind.String())
			return errors.New(errorMsg)
		}

		lastToken = token
	}

	if !state.isEOF {
		return errors.New("Unexpected end of expression")
	}
	return nil
}

func getLexerStateForToken(kind TokenKind) (lexerState, error) {

	for _, possibleState := range validLexerStates {

		if possibleState.kind == kind {
			return possibleState, nil
		}
	}

	errorMsg := fmt.Sprintf("No lexer state found for token kind '%v'\n", kind.String())
	return validLexerStates[0], errors.New(errorMsg)
}
