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
    MINUS = "-"
    BANG = "!"
    ASTERISK = "*"
    SLASH = "/"

    LT = "<"
    GT = ">"

	// Delimiters
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	COMMA = ","
	SEMICOLON = ";"
	LBRACKET = "["
	RBRACKET = "]"
	COLON = ":"

	// Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
	IF = "IF"
	ELSE = "ELSE"

	TRUE = "TRUE"
	FALSE  = "FALSE"
	RETURN = "RETURN"
	EQ = "=="
	NOT_EQ = "!="
	STRING = "STRING"

)

type TokenType string

type Token struct{
	Type TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
	"if": IF,
	"else": ELSE,
	"true": TRUE,
	"false": FALSE,
	"return": RETURN,
 }

func LookupIdent(ident string) TokenType{
	if t, ok := keywords[ident]; ok{
		return t
	}
	return IDENT
}

