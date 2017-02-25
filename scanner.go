package argv

import "unicode"

type Scanner struct {
	env map[string]string

	text      []rune
	rpos      int
	dollarBuf []rune
}

func NewScanner(text []rune, env map[string]string) *Scanner {
	return &Scanner{
		text: text,
		env:  env,
	}
}

func (s *Scanner) Env() map[string]string {

	return s.env
}

const _RUNE_EOF = 0

func (s *Scanner) nextRune() rune {
	if s.rpos >= len(s.text) {
		return _RUNE_EOF
	}

	r := s.text[s.rpos]
	s.rpos += 1
	return r
}

func (s *Scanner) unreadRune(r rune) {
	if r != _RUNE_EOF {
		s.rpos -= 1
	}
}

func (s *Scanner) isEscapeChars(r rune) (rune, bool) {
	switch r {
	case 'a':
		return '\a', true
	case 'b':
		return '\b', true
	case 'f':
		return '\f', true
	case 'n':
		return '\n', true
	case 'r':
		return '\r', true
	case 't':
		return '\t', true
	case 'v':
		return '\v', true
	case '\\':
		return '\\', true
	case '$':
		return '$', true
	}
	return r, false
}

func (s *Scanner) endEnv(r rune) bool {
	if r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
		return false
	}
	return true
}

type TokenType uint32

type Token struct {
	Type  TokenType
	Value []rune
}

const (
	TOK_STRING TokenType = iota + 1
	TOK_PIPE
	TOK_REVERSEQUOTE
	TOK_SPACE
	TOK_EOF
)

func (s *Scanner) getEnv(name string) string {
	return s.env[name]
}

func (s *Scanner) specialVar(r rune) (string, bool) {
	switch r {
	case '0', '*', '#', '@', '?', '$':
		v, has := s.env[string(r)]
		return v, has
	default:
		return "", false
	}
}

func (s *Scanner) checkDollarStart(tok *Token, r rune, from, switchTo uint8) uint8 {
	state := from
	nr := s.nextRune()
	if val, has := s.specialVar(nr); has {
		if val != "" {
			tok.Value = append(tok.Value, []rune(val)...)
		}
	} else if s.endEnv(nr) {
		tok.Value = append(tok.Value, r)
		s.unreadRune(nr)
	} else {
		state = switchTo
		s.dollarBuf = append(s.dollarBuf[:0], nr)
	}
	return state
}

func (s *Scanner) checkDollarEnd(tok *Token, r rune, from, switchTo uint8) uint8 {
	var state = from
	if s.endEnv(r) {
		tok.Value = append(tok.Value, []rune(s.getEnv(string(s.dollarBuf)))...)
		state = switchTo
		s.unreadRune(r)
	} else {
		s.dollarBuf = append(s.dollarBuf, r)
	}
	return state
}

func (s *Scanner) Next() (Token, error) {
	const (
		INITIAL = iota + 1
		SPACE
		REVERSE_QUOTE
		STRING
		STRING_DOLLAR
		STRING_QUOTE_SINGLE
		STRING_QUOTE_DOUBLE
		STRING_QUOTE_DOUBLE_DOLLAR
	)

	var (
		tok Token

		state uint8 = INITIAL
	)
	s.dollarBuf = s.dollarBuf[:0]
	for {
		r := s.nextRune()
		switch state {
		case INITIAL:
			switch {
			case r == _RUNE_EOF:
				tok.Type = TOK_EOF
				return tok, nil
			case r == '|':
				tok.Type = TOK_PIPE
				return tok, nil
			case r == '`':
				state = REVERSE_QUOTE
			case unicode.IsSpace(r):
				state = SPACE
				s.unreadRune(r)
			default:
				state = STRING
				s.unreadRune(r)
			}
		case SPACE:
			if r == _RUNE_EOF || !unicode.IsSpace(r) {
				s.unreadRune(r)
				tok.Type = TOK_SPACE
				return tok, nil
			}
		case REVERSE_QUOTE:
			switch r {
			case _RUNE_EOF:
				return tok, ErrInvalidSyntax
			case '`':
				tok.Type = TOK_REVERSEQUOTE
				return tok, nil
			default:
				tok.Value = append(tok.Value, r)
			}
		case STRING:
			switch {
			case r == _RUNE_EOF || r == '|' || r == '`' || unicode.IsSpace(r):
				tok.Type = TOK_STRING
				s.unreadRune(r)
				return tok, nil
			case r == '\'':
				state = STRING_QUOTE_SINGLE
			case r == '"':
				state = STRING_QUOTE_DOUBLE
			case r == '\\':
				nr := s.nextRune()
				if nr == _RUNE_EOF {
					return tok, ErrInvalidSyntax
				}
				tok.Value = append(tok.Value, nr)
			case r == '$':
				state = s.checkDollarStart(&tok, r, state, STRING_DOLLAR)
			default:
				tok.Value = append(tok.Value, r)
			}
		case STRING_DOLLAR:
			state = s.checkDollarEnd(&tok, r, state, STRING)
		case STRING_QUOTE_SINGLE:
			switch r {
			case _RUNE_EOF:
				return tok, ErrInvalidSyntax
			case '\'':
				state = STRING
			case '\\':
				nr := s.nextRune()
				if escape, ok := s.isEscapeChars(nr); ok {
					tok.Value = append(tok.Value, escape)
				} else {
					tok.Value = append(tok.Value, r)
					s.unreadRune(nr)
				}
			default:
				tok.Value = append(tok.Value, r)
			}
		case STRING_QUOTE_DOUBLE:
			switch r {
			case _RUNE_EOF:
				return tok, ErrInvalidSyntax
			case '"':
				state = STRING
			case '\\':
				nr := s.nextRune()
				if nr == _RUNE_EOF {
					return tok, ErrInvalidSyntax
				}
				if escape, ok := s.isEscapeChars(nr); ok {
					tok.Value = append(tok.Value, escape)
				} else {
					tok.Value = append(tok.Value, r)
					s.unreadRune(nr)
				}
			case '$':
				state = s.checkDollarStart(&tok, r, state, STRING_QUOTE_DOUBLE_DOLLAR)
			default:
				tok.Value = append(tok.Value, r)
			}
		case STRING_QUOTE_DOUBLE_DOLLAR:
			state = s.checkDollarEnd(&tok, r, state, STRING_QUOTE_DOUBLE)
		}
	}
}

func Scan(text []rune, env map[string]string) ([]Token, error) {
	s := NewScanner(text, env)
	var tokens []Token
	for {
		tok, err := s.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TOK_EOF {
			break
		}
	}
	return tokens, nil
}
