package govaluate

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// ToSQLQuery returns a string representing this expression as if it were written in SQL.
// This function assumes that all parameters exist within the same table, and that the table essentially represents
// a serialized object of some sort (e.g., hibernate).
// If your data model is more normalized, you may need to consider iterating through each actual token given by `Tokens()`
// to create your query.
// Boolean values are considered to be "1" for true, "0" for false.
// Times are formatted according to this.QueryDateFormat.
func (expr Expression) ToSQLQuery() (string, error) {

	var stream *tokenStream
	var transactions *expressionOutputStream
	var transaction string
	var err error

	stream = newTokenStream(expr.tokens)
	transactions = new(expressionOutputStream)

	for stream.hasNext() {

		transaction, err = expr.findNextSQLString(stream, transactions)
		if err != nil {
			return "", err
		}

		transactions.add(transaction)
	}

	return transactions.createString(" "), nil
}

func (expr Expression) findNextSQLString(stream *tokenStream, transactions *expressionOutputStream) (string, error) {

	var token ExpressionToken
	var ret string

	token = stream.next()

	switch token.Kind {

	case stringToken:
		ret = fmt.Sprintf("'%v'", token.Value)
	case pattern:
		ret = fmt.Sprintf("'%s'", token.Value.(*regexp.Regexp).String())
	case timeToken:
		ret = fmt.Sprintf("'%s'", token.Value.(time.Time).Format(expr.QueryDateFormat))

	case logicalop:
		switch logicalSymbols[token.Value.(string)] {

		case and:
			ret = "AND"
		case or:
			ret = "OR"
		}

	case boolean:
		if token.Value.(bool) {
			ret = "1"
		} else {
			ret = "0"
		}

	case variable:
		ret = fmt.Sprintf("[%s]", token.Value.(string))

	case numeric:
		ret = fmt.Sprintf("%g", token.Value.(float64))

	case comparator:
		switch comparatorSymbols[token.Value.(string)] {

		case eq:
			ret = "="
		case neq:
			ret = "<>"
		case req:
			ret = "RLIKE"
		case nreq:
			ret = "NOT RLIKE"
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}

	case ternary:

		switch ternarySymbols[token.Value.(string)] {

		case coalesce:

			left := transactions.rollback()
			right, err := expr.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("COALESCE(%v, %v)", left, right)
		case ternaryTrue:
			fallthrough
		case ternaryFalse:
			return "", errors.New("Ternary operators are unsupported in SQL output")
		}
	case prefix:
		switch prefixSymbols[token.Value.(string)] {

		case invert:
			ret = fmt.Sprintf("NOT")
		default:

			right, err := expr.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("%s%s", token.Value.(string), right)
		}
	case modifier:

		switch modifierSymbols[token.Value.(string)] {

		case exponent:

			left := transactions.rollback()
			right, err := expr.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("POW(%s, %s)", left, right)
		case modulus:

			left := transactions.rollback()
			right, err := expr.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("MOD(%s, %s)", left, right)
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}
	case clause:
		ret = "("
	case clauseClose:
		ret = ")"
	case separator:
		ret = ","

	default:
		errorMsg := fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind)
		return "", errors.New(errorMsg)
	}

	return ret, nil
}
