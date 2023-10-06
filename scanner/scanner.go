package scanner

import (
	"golox/errorhandling"
	"golox/token"
	"unicode"
)

type Scanner struct {
	source               string
	tokens               []token.Token
	start, current, line int
}

func NewScanner(source string) *Scanner {
	ret := Scanner{source: source, tokens: []token.Token{}, start: 0, current: 0, line: 1}

	return &ret
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(
		s.tokens,
		token.Token{
			Token_type: token.EOF,
			Lexeme:     "",
			Literal:    nil,
			Line:       s.line,
		},
	)

	return s.tokens
}

func (s *Scanner) scanToken() {
	var t token.TokenType
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)

	case ')':
		s.addToken(token.RIGHT_PAREN)

	case '{':
		s.addToken(token.LEFT_BRACE)

	case '}':
		s.addToken(token.RIGHT_BRACE)

	case ',':
		s.addToken(token.COMMA)

	case '.':
		s.addToken(token.DOT)

	case '-':
		s.addToken(token.MINUS)

	case '+':
		s.addToken(token.PLUS)

	case ';':
		s.addToken(token.SEMICOLON)

	case '*':
		s.addToken(token.STAR)

	case '!':
		if s.match('=') {
			t = token.BANG_EQUAL
		} else {
			t = token.BANG
		}
		s.addToken(t)

	case '=':
		if s.match('=') {
			t = token.EQUAL_EQUAL
		} else {
			t = token.EQUAL
		}
		s.addToken(t)

	case '<':
		if s.match('=') {
			t = token.LESS_EQUAL
		} else {
			t = token.LESS
		}
		s.addToken(t)

	case '>':
		if s.match('=') {
			t = token.GREATER_EQUAL
		} else {
			t = token.GREATER
		}
		s.addToken(t)

	case '/':
		if !s.match('/') {
			s.addToken(token.SLASH)
		} else {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		}

	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line += 1
	case '"':
		s.tokenize_string()

	default:
		if unicode.IsDigit(c) {
			s.tokenize_number()
		} else if unicode.IsLetter(c) || c == '_' {
            s.tokenize_identifier()
		} else {
			errorhandling.Report(s.line, s.source[s.start:s.current], "Unexpected character.")
		}
	}
}

func (s *Scanner) addToken(t token.TokenType) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t token.TokenType, literal *string) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.Token{Token_type: t, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *Scanner) tokenize_identifier() {
	var identifier string

	for c := s.peek(); unicode.IsDigit(c) || unicode.IsLetter(c) || c == '_'; c = s.advance() {
	}

	identifier = s.source[s.start:s.current]

	if token.KeywordMap[identifier] != 0 {
		s.addToken(token.KeywordMap[identifier])
	} else {
		s.addToken(token.IDENTIFIER)
	}

}

func (s *Scanner) tokenize_number() {
	var new_string string
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()
		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	new_string = s.source[s.start:s.current]

	s.addTokenLiteral(token.NUMBER, &new_string)

	// var decimal bool = false
	// var c rune
	//
	// for c = s.peek(); !s.isAtEnd(); c = s.peek() {
	// 	if unicode.IsDigit(c) {
	// 		s.advance()
	// 	} else if c == '.' {
	// 		if decimal {
	// 			errorhandling.Report(s.line, s.source[s.start:s.current+1], "Unexpected character.")
	// 			return
	// 		}
	// 		decimal = true
	//            s.advance()
	// 	} else {
	//            break
	// 	}
	// }
	//
	// new_string = s.source[s.start:s.current]
	//
	// s.addTokenLiteral(token.NUMBER, &new_string)
}

func (s *Scanner) tokenize_string() {
	var new_string string

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		errorhandling.Report(s.line, s.source[s.start:s.current], "Unterminated string.")
		return
	}

	s.advance()
	new_string = s.source[s.start+1 : s.current-1]

	s.addTokenLiteral(token.STRING, &new_string)

}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || rune(s.source[s.current]) != expected {
		return false
	}

	s.current += 1

	return true
}

func (s *Scanner) advance() rune {
    if s.isAtEnd() {
        return rune(0)
    }
	r := s.source[s.current]
	s.current += 1

	return rune(r)
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}

	return rune(s.source[s.current])
}

func (s Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return rune(0)
	}

	return rune(s.source[s.current+1])
}
