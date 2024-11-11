package lox

import "fmt"

type ScanError struct {
	line    int
	message string
}

func NewScanError(line int, message string) *ScanError {
	return &ScanError{line, message}
}

func (s *ScanError) Error() string {
	return fmt.Sprintf("[line %d]: Scan Error: %s", s.line, s.message)
}

type ParseError struct {
	token   *Token
	message string
}

func NewParseError(token *Token, message string) *ParseError {
	return &ParseError{token, message}
}

func (p *ParseError) Error() string {
	var where string
	if p.token.Type == EOF {
		where = "EOF"
	} else {
		where = p.token.Lexeme
	}
	return fmt.Sprintf("[line %d] at %s: Parse Error: %s",
		p.token.Line, where, p.message)
}

type RuntimeError struct {
	token   *Token
	message string
}

func NewRuntimeError(token *Token, message string) *RuntimeError {
	return &RuntimeError{token, message}
}

func (r *RuntimeError) Error() string {
	var where string
	if r.token.Type == EOF {
		where = "EOF"
	} else {
		where = r.token.Lexeme
	}
	return fmt.Sprintf("[line %d] at %s: Runtime Error: %s",
		r.token.Line, where, r.message)
}

type ResolveError struct {
	token   *Token
	message string
}

func NewResolveError(token *Token, message string) *ResolveError {
	return &ResolveError{token, message}
}

func (r *ResolveError) Error() string {
	var where string
	if r.token.Type == EOF {
		where = "EOF"
	} else {
		where = r.token.Lexeme
	}
	return fmt.Sprintf("[line %d] at %s: Resolve Error: %s",
		r.token.Line, where, r.message)
}
