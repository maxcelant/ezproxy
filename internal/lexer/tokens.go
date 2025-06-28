package lex

type TokenType int

const (
	VERB TokenType = iota
	ROUTE
	PROTOCOL
	PROTOCOL_VERSION
	HOST
	ACCEPT
	HEADER
	CONTENT_LENGTH
	CONTENT
)

type Token struct {
	TokenType
	value string
}
