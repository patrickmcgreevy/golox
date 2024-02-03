package scanner

import (
	"fmt"
	"strconv"
	"unicode"
)

type Scanner struct {
	source               string
	tokens               []Token
	start, current, line int
}

type ScannerError struct {
	line int
	seq  string
	err  string
}

func (e ScannerError) Error() string {
	return fmt.Sprintf("parsing error on line %d at \"%s\": %s", e.line, e.seq, e.err)
}

func NewScanner(source string) *Scanner {
	ret := Scanner{source: source, tokens: []Token{}, start: 0, current: 0, line: 1}

	return &ret
}

func (s *Scanner) ScanTokens() ([]Token, *ScannerError) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	s.tokens = append(
		s.tokens,
		Token{
			Token_type: EOF,
			Lexeme:     "",
			Literal:    nil,
			Line:       s.line,
		},
	)

	return s.tokens, nil
}

func (s *Scanner) scanToken() *ScannerError {
	var t TokenType
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)

	case ')':
		s.addToken(RIGHT_PAREN)

	case '{':
		s.addToken(LEFT_BRACE)

	case '}':
		s.addToken(RIGHT_BRACE)

	case ',':
		s.addToken(COMMA)

	case '.':
		s.addToken(DOT)

	case '-':
		s.addToken(MINUS)

	case '+':
		s.addToken(PLUS)

	case ';':
		s.addToken(SEMICOLON)

	case '*':
		s.addToken(STAR)

	case '!':
		if s.match('=') {
			t = BANG_EQUAL
		} else {
			t = BANG
		}
		s.addToken(t)

	case '=':
		if s.match('=') {
			t = EQUAL_EQUAL
		} else {
			t = EQUAL
		}
		s.addToken(t)

	case '<':
		if s.match('=') {
			t = LESS_EQUAL
		} else {
			t = LESS
		}
		s.addToken(t)

	case '>':
		if s.match('=') {
			t = GREATER_EQUAL
		} else {
			t = GREATER
		}
		s.addToken(t)

	case '/':
		if !s.match('/') {
			s.addToken(SLASH)
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
		err := s.tokenize_string()
		if err != nil {
			return err
		}

	default:
		if unicode.IsDigit(c) {
			s.tokenize_number()
		} else if unicode.IsLetter(c) || c == '_' {
			s.tokenize_identifier()
		} else {
			return &ScannerError{
				line: s.line,
				seq:  s.source[s.start:s.current],
				err:  "unexpected character",
			}
		}
	}

	return nil
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Token_type: t, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *Scanner) tokenize_identifier() {
	var identifier string

	for c := s.peek(); unicode.IsDigit(c) || unicode.IsLetter(c) || c == '_'; c = s.peek() {
		s.advance()
	}

	identifier = s.source[s.start:s.current]

	if KeywordMap[identifier] != 0 {
		s.addToken(KeywordMap[identifier])
	} else {
		s.addToken(IDENTIFIER)
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
	num, err := strconv.ParseFloat(new_string, 64)
	if err != nil {
		panic("Not a number!")
	}

	s.addTokenLiteral(NUMBER, num)
}

func (s *Scanner) tokenize_string() *ScannerError {
	var new_string string

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		// errorhandling.Report(s.line, s.source[s.start:s.current], "Unterminated string.")
		return &ScannerError{
			line: s.line,
			seq:  s.source[s.start:s.current],
			err:  "unterminated string",
		}
	}

	s.advance()
	new_string = s.source[s.start+1 : s.current-1]

	s.addTokenLiteral(STRING, &new_string)

	return nil
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
