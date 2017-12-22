package token

const(

	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// Identifiers
	IDENT = "IDENT"  // identifiers add, foobar, x, y
 	INT = "INT"   // integers 23, 12343

	//OPERATORS
	ASSIGN = "="
	PLUS = "+"

	// Delimiters
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	COMMA = ","
	SEMICOLON = ";"

	// Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"

)

type TokenType string

type Token struct{
	Type TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType{
	if t, ok := keywords[ident]; ok{
		return t
	}
	return IDENT
}

