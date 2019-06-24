package govaluate

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
	"time"
	"unicode"
)

func parseTokens(expression string, functions map[string]ExpressionFunction) ([]ExpressionToken, error) {

	var ret []ExpressionToken
	var token ExpressionToken
	var stream scanner.Scanner
	var state lexerState
	var err error
	var found bool

	reader := strings.NewReader(expression)
	stream.Init(reader)
	stream.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments
	stream.Error = func(*scanner.Scanner, string) {}
	state = validLexerStates[0]

	for stream.Peek() != scanner.EOF {

		token, err, found = readToken(&stream, state, functions)

		if err != nil {
			return ret, err
		}

		if !found {
			break
		}

		state, err = getLexerStateForToken(token.Kind)
		if err != nil {
			return ret, err
		}

		// append this valid token
		ret = append(ret, token)
	}

	err = checkBalance(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func readToken(stream *scanner.Scanner, state lexerState, functions map[string]ExpressionFunction) (ExpressionToken, error, bool) {

	var fnFunction ExpressionFunction
	var ret ExpressionToken
	var tokenValue interface{}
	var tokenTime time.Time
	var tokenString string
	var kind TokenKind
	var found bool
	var completed bool
	var err error

	// numeric is 0-9, or . or 0x followed by digits
	// string starts with '
	// variable is alphanumeric, always starts with a letter
	// bracket always means variable
	// symbols are anything non-alphanumeric
	// all others read into a buffer until they reach the end of the stream
	kind = unknown
	switch stream.Scan() {
	case scanner.EOF:
		break
	case scanner.Float:
		kind = numeric
		tokenValue, err = strconv.ParseFloat(stream.TokenText(), 64)
		if err != nil {
			errorMsg := fmt.Sprintf("Unable to parse numeric value '%v' to float64\n", stream.TokenText())
			return ExpressionToken{}, errors.New(errorMsg), false
		}
	case scanner.Int:
		kind = numeric
		i, err := strconv.ParseInt(stream.TokenText(), 0, 64)
		tokenValue = float64(i)
		if err != nil {
			errorMsg := fmt.Sprintf("Unable to parse numeric value '%v' to float64\n", stream.TokenText())
			return ExpressionToken{}, errors.New(errorMsg), false
		}
	case ',':
		tokenValue = ","
		kind = separator
	case '[':
		tokenValue, completed = readUntilFalse(stream, true, isNotClosingBracket)
		kind = variable

		if !completed {
			return ExpressionToken{}, errors.New("Unclosed parameter bracket"), false
		}

		// above method normally rewinds us to the closing bracket, which we want to skip.
		stream.Next()
		break

	case scanner.Ident:
		// regular variable - or function?

		tokenString = stream.TokenText()

		//Hack for crazy escapes in variable names
		if stream.Peek() == '\\' {
			s, _ := readUntilFalse(stream, true, isVariableName)
			tokenString = tokenString + s
		}

		tokenValue = tokenString
		kind = variable

		switch tokenValue {
		// boolean?
		case "true":

			kind = boolean
			tokenValue = true
		case "false":
			kind = boolean
			tokenValue = false
		// textual operator?
		case "in", "IN":

			// force lower case for consistency
			tokenValue = "in"
			kind = comparator

		default:

			// function?
			fnFunction, found = functions[tokenString]
			if found {
				kind = function
				tokenValue = fnFunction
				break
			}

			// accessor?
			if stream.Peek() == '.' {

				splits := []string{tokenString}
				for stream.Peek() == '.' {
					stream.Scan()
					// check that it doesn't end with a hanging period
					if stream.Scan() != scanner.Ident {
						errorMsg := fmt.Sprintf("Hanging accessor on token '%s'", tokenString)
						return ExpressionToken{}, errors.New(errorMsg), false
					}

					tokenString = stream.TokenText()
					//Hack for crazy escapes in variable names
					if stream.Peek() == '\'' {
						s, _ := readUntilFalse(stream, true, isVariableName)
						tokenString = tokenString + s
					}

					// check that none of them are unexported
					firstCharacter := getFirstRune(tokenString)
					if unicode.ToUpper(firstCharacter) != firstCharacter {
						errorMsg := fmt.Sprintf("Unable to access unexported field '%s' in token '%s'", tokenString, strings.Join(splits, "."))
						return ExpressionToken{}, errors.New(errorMsg), false
					}

					splits = append(splits, tokenString)
				}

				kind = accessor
				tokenValue = splits
			}
		}

	case scanner.String, scanner.RawString, '\'':
		tokenString = stream.TokenText()
		if tokenString == "'" {
			var tokenBuffer bytes.Buffer

			for c := stream.Next(); c != '\''; c = stream.Next() {
				if c == '\\' {
					c = stream.Next()
					if c != '\'' {
						tokenBuffer.WriteRune('\\')
					}
				}
				if c == scanner.EOF || c == '\n' {
					return ExpressionToken{}, fmt.Errorf("Unclosed string literal '%s", tokenBuffer.String()), false
				}
				tokenBuffer.WriteRune(c)

			}

			tokenString = "\"" + strings.Replace(tokenBuffer.String(), `"`, `\"`, -1) + "\""
		}

		tokenValue, err = strconv.Unquote(tokenString)

		if err != nil {
			return ExpressionToken{}, err, false
		}

		// check to see if this can be parsed as a time.
		tokenTime, found = tryParseTime(tokenValue.(string))
		if found {
			kind = timeToken
			tokenValue = tokenTime
		} else {
			kind = stringToken
		}
	case '(':
		tokenValue = '('
		kind = clause
	case ')':
		tokenValue = ')'
		kind = clauseClose

	default:

		// must be a known symbol
		tokenString = readTokenUntilFalse(stream, isOperation)
		tokenValue = tokenString

		// quick hack for the case where "-" can mean "prefixed negation" or "minus", which are used
		// very differently.
		if state.canTransitionTo(prefix) {
			_, found = prefixSymbols[tokenString]
			if found {

				kind = prefix
				break
			}
		}
		_, found = modifierSymbols[tokenString]
		if found {

			kind = modifier
			break
		}

		_, found = logicalSymbols[tokenString]
		if found {

			kind = logicalop
			break
		}

		_, found = comparatorSymbols[tokenString]
		if found {

			kind = comparator
			break
		}

		_, found = ternarySymbols[tokenString]
		if found {

			kind = ternary
			break
		}

		errorMessage := fmt.Sprintf("Invalid token: '%s'", tokenString)
		return ret, errors.New(errorMessage), false
	}

	ret.Kind = kind
	ret.Value = tokenValue

	return ret, nil, (kind != unknown)
}

func readTokenUntilFalse(stream *scanner.Scanner, condition func(rune) bool) string {

	var tokenBuffer bytes.Buffer
	var character rune

	tokenBuffer.WriteString(stream.TokenText())

	for {
		character = stream.Peek()
		if character == scanner.EOF || !condition(character) {
			break
		}
		stream.Next()
		tokenBuffer.WriteString(string(character))

	}

	return tokenBuffer.String()
}

/*
	Returns the string that was read until the given [condition] was false.
	Returns false if the stream ended before condition was met.
*/
func readUntilFalse(stream *scanner.Scanner, allowEscaping bool, condition func(rune) bool) (string, bool) {

	var tokenBuffer bytes.Buffer
	var character rune
	var conditioned bool

	conditioned = false

	for {
		character = stream.Peek()
		if character == scanner.EOF {
			break
		}
		if allowEscaping && character == '\\' {
			stream.Next()

			if character == stream.Peek() {
				break
			}
		} else if !condition(character) {
			conditioned = true
			break
		}
		character = stream.Next()

		// Use backslashes to escape anything

		tokenBuffer.WriteRune(character)

	}

	return tokenBuffer.String(), conditioned
}

/*
	Checks to see if any optimizations can be performed on the given [tokens], which form a complete, valid expression.
	The returns slice will represent the optimized (or unmodified) list of tokens to use.
*/
func optimizeTokens(tokens []ExpressionToken) ([]ExpressionToken, error) {

	var token ExpressionToken
	var symbol OperatorSymbol
	var err error
	var index int

	for index, token = range tokens {

		// if we find a regex operator, and the right-hand value is a constant, precompile and replace with a pattern.
		if token.Kind != comparator {
			continue
		}

		symbol = comparatorSymbols[token.Value.(string)]
		if symbol != req && symbol != nreq {
			continue
		}

		index++
		token = tokens[index]
		if token.Kind == stringToken {

			token.Kind = pattern
			token.Value, err = regexp.Compile(token.Value.(string))

			if err != nil {
				return tokens, err
			}

			tokens[index] = token
		}
	}
	return tokens, nil
}

/*
	Checks the balance of tokens which have multiple parts, such as parenthesis.
*/
func checkBalance(tokens []ExpressionToken) error {

	var stream *tokenStream
	var token ExpressionToken
	var parens int

	stream = newTokenStream(tokens)

	for stream.hasNext() {

		token = stream.next()
		if token.Kind == clause {
			parens++
			continue
		}
		if token.Kind == clauseClose {
			parens--
			continue
		}
	}

	if parens != 0 {
		return errors.New("Unbalanced parenthesis")
	}
	return nil
}

func isDigit(character rune) bool {
	return unicode.IsDigit(character)
}

func isOperation(character rune) bool {
	switch character {
	case '=', '!', '<', '>', '~', '&', '|', '+', '-', '*', '/', '^', '%', ':', '?':
		return true
	default:
		return false
	}
}

func isVariableName(character rune) bool {

	return unicode.IsLetter(character) ||
		unicode.IsDigit(character) ||
		character == '_'
}

func isNotClosingBracket(character rune) bool {

	return character != ']'
}

/*
	Attempts to parse the [candidate] as a Time.
	Tries a series of standardized date formats, returns the Time if one applies,
	otherwise returns false through the second return.
*/
func tryParseTime(candidate string) (time.Time, bool) {

	var ret time.Time
	var found bool

	timeFormats := [...]string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.Kitchen,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",                         // RFC 3339
		"2006-01-02 15:04",                   // RFC 3339 with minutes
		"2006-01-02 15:04:05",                // RFC 3339 with seconds
		"2006-01-02 15:04:05-07:00",          // RFC 3339 with seconds and timezone
		"2006-01-02T15Z0700",                 // ISO8601 with hour
		"2006-01-02T15:04Z0700",              // ISO8601 with minutes
		"2006-01-02T15:04:05Z0700",           // ISO8601 with seconds
		"2006-01-02T15:04:05.999999999Z0700", // ISO8601 with nanoseconds
	}

	for _, format := range timeFormats {

		ret, found = tryParseExactTime(candidate, format)
		if found {
			return ret, true
		}
	}

	return time.Now(), false
}

func tryParseExactTime(candidate string, format string) (time.Time, bool) {

	var ret time.Time
	var err error

	ret, err = time.ParseInLocation(format, candidate, time.Local)
	if err != nil {
		return time.Now(), false
	}

	return ret, true
}

func getFirstRune(candidate string) rune {

	for _, character := range candidate {
		return character
	}

	return 0
}
