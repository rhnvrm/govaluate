package govaluate

import (
	"bytes"
)

// expressionOutputStream holds a series of "transactions" which represent each token as it is output by an outputter (such as ToSQLQuery()).
// Some outputs (such as SQL) require a function call or non-c-like syntax to represent an expression.
// To accomplish this, this struct keeps track of each translated token as it is output, and can return and rollback those transactions.
type expressionOutputStream struct {
	transactions []string
}

func (exprOS *expressionOutputStream) add(transaction string) {
	exprOS.transactions = append(exprOS.transactions, transaction)
}

func (exprOS *expressionOutputStream) rollback() string {

	index := len(exprOS.transactions) - 1
	ret := exprOS.transactions[index]

	exprOS.transactions = exprOS.transactions[:index]
	return ret
}

func (exprOS *expressionOutputStream) createString(delimiter string) string {

	var retBuffer bytes.Buffer
	var transaction string

	penultimate := len(exprOS.transactions) - 1

	for i := 0; i < penultimate; i++ {

		transaction = exprOS.transactions[i]

		retBuffer.WriteString(transaction)
		retBuffer.WriteString(delimiter)
	}
	retBuffer.WriteString(exprOS.transactions[penultimate])

	return retBuffer.String()
}
