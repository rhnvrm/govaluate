package govaluate

import (
	"errors"
)

/*
	EvaluableExpression represents a set of ExpressionTokens which, taken together,
	represent an arbitrary expression that can be evaluated down into a single value.
*/
type EvaluableExpression struct {

	tokens []ExpressionToken
	inputExpression string
}

/*
	Creates a new EvaluableExpression from the given [expression] string.
	Returns an error if the given expression has invalid syntax.
*/
func NewEvaluableExpression(expression string) (*EvaluableExpression, error) {

	var ret *EvaluableExpression;
	var err error

	ret = new(EvaluableExpression)
	ret.inputExpression = expression;
	ret.tokens, err = parseTokens(expression)

	if(err != nil) {
		return nil, err
	}
	return ret, nil
}

/*
	Evaluate runs the entire expression using the given [parameters]. 
	Each parameter is mapped from a string to a value, such as "foo" = 1.0. 
	If the expression contains a reference to the variable "foo", it will be taken from parameters["foo"].

	This function returns errors if the combination of expression and parameters cannot be run,
	such as if a string parameter is given in an expression that expects it to be a boolean. 
	e.g., "foo == true", where foo is any string.
	These errors are almost exclusively returned for parameters not being present, or being of the wrong type.
	Structural problems with the expression (unexpected tokens, unexpected end of expression, etc) are discovered
	during parsing of the expression in NewEvaluableExpression.

	In all non-error circumstances, this returns the single value result of the expression and parameters given.
	e.g., if the expression is "1 + 1", Evaluate will return 2.0.
	e.g., if the expression is "foo + 1" and parameters contains "foo" = 2, Evaluate will return 3.0
*/
func (this EvaluableExpression) Evaluate(parameters map[string]interface{}) (interface{}, error) {

	var stream *tokenStream;

	stream = newTokenStream(this.tokens);
	return evaluateTokens(stream, parameters);
}

func evaluateTokens(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	if(stream.hasNext()) {
		return evaluateLogical(stream, parameters);
	}
	return nil, nil;
}

func evaluateLogical(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateComparator(stream, parameters);	

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {

		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = LOGICAL_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break;
		}

		switch(symbol) {

			case OR		:	if(value != nil) {
							return evaluateLogical(stream, parameters);
						} else {
							value, err = evaluateComparator(stream, parameters);
						}
			case AND	:	if(value == nil) {
							return evaluateLogical(stream, parameters);
						} else {
							value, err = evaluateComparator(stream, parameters);
						}
		}

		if(err != nil) {
			return nil, err;
		}
	}

	stream.rewind();
	return value, nil;
}

func evaluateComparator(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value, rightValue interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateAdditiveModifier(stream, parameters);

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {

		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = COMPARATOR_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break
		}

		rightValue, err = evaluateAdditiveModifier(stream, parameters);
		if(err != nil) {
			return nil, err;
		}

		switch(symbol) {

			case LT		:	return (value.(float64) < rightValue.(float64)), nil;
			case LTE	:	return (value.(float64) <= rightValue.(float64)), nil;
			case GT		:	return (value.(float64) > rightValue.(float64)), nil;
			case GTE	:	return (value.(float64) >= rightValue.(float64)), nil;
			case EQ		:	return (value == rightValue), nil;
			case NEQ	:	return (value != rightValue), nil;
		}
	}

	stream.rewind();
	return value, nil;
}

func evaluateAdditiveModifier(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value, rightValue interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateMultiplicativeModifier(stream, parameters);

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {
		
		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = MODIFIER_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break;
		}

		switch(symbol) {

			case PLUS	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						value = value.(float64) + rightValue.(float64);

			case MINUS	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}

						return value.(float64) - rightValue.(float64), nil;

			default		:	stream.rewind();
						return value, nil;
		}
	}

	stream.rewind();
	return value, nil;
}

func evaluateMultiplicativeModifier(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value, rightValue interface{};
	var symbol OperatorSymbol;
	var err error;
	var keyFound bool;

	value, err = evaluateValue(stream, parameters);

	if(err != nil) {
		return nil, err;
	}

	for stream.hasNext() {

		token = stream.next();

		if(!isString(token.Value)) {
			break;
		}

		symbol, keyFound = MODIFIER_SYMBOLS[token.Value.(string)];
		if(!keyFound) {
			break;
		}

		switch(symbol) {

			case MULTIPLY	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						return value.(float64) * rightValue.(float64), nil;

			case DIVIDE	:	rightValue, err = evaluateMultiplicativeModifier(stream, parameters);
						if(err != nil) {
							return nil, err;
						}
						return value.(float64) / rightValue.(float64), nil;

			default		:	stream.rewind();
						return value, nil;
		}
	}

	stream.rewind();	
	return value, nil;
}

func evaluateValue(stream *tokenStream, parameters map[string]interface{}) (interface{}, error) {

	var token ExpressionToken;
	var value interface{};
	var errorMessage, variableName string;
	var err error;

	token = stream.next();

	switch(token.Kind) {

		case CLAUSE	:	value, err = evaluateTokens(stream, parameters);
					if(err != nil) {
						return nil, err;
					}

					token = stream.next();
					if(token.Kind != CLAUSE_CLOSE) {

						return nil, errors.New("Unbalanced parenthesis");
					}

					return value, nil;

		case VARIABLE	:	variableName = token.Value.(string);
					value = parameters[variableName];

					if(value == nil) {
						errorMessage = "No parameter '"+ variableName +"' found."
						return nil, errors.New(errorMessage);
					}

					return value, nil;

		case NUMERIC	:	fallthrough
		case STRING	:	fallthrough
		case BOOLEAN	:	return token.Value, nil;
		default		:	break;
	}

	stream.rewind();
	return nil, errors.New("Unable to evaluate token kind: " + GetTokenKindString(token.Kind));
}

/*
	Returns an array representing the ExpressionTokens that make up this expression.
*/
func (this EvaluableExpression) Tokens() []ExpressionToken {

	return this.tokens;
}

/*
	Returns the original expression used to create this EvaluableExpression.
*/
func (this EvaluableExpression) String() string {

	return this.inputExpression;
}

func isString(value interface{}) bool {

	switch value.(type) {
		case string	:	return true;
		default		:	break;
	}
	return false;
}