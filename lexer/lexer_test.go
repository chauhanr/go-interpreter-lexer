package lexer

import (
	"testing"
	"go-interpreter-lexer/token"
)

type expectedTokens struct{
	expectedType token.TokenType
	expectedLiteral string
}

var tokenTests = [] struct{
	input string
	expectedTokens []expectedTokens
}{
	{
		`=+(){},;`,
		[]expectedTokens{
			{token.ASSIGN, "="},
			{token.PLUS, "+"},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{ token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		},
	},
	{
		`let five = 5;
let ten = 10;
let add = fn(x,y){
   x+y;
};
let result = add(five,ten);`,

		[]expectedTokens{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.PLUS, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{ token.LET, "let"},
			{token.IDENT, "result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		},
	},

}

func TestNextToken(t *testing.T){

	for _, tCase := range tokenTests{
		l := New(tCase.input)
       // t.Logf("Input : "+tCase.input)
		for i, tt := range tCase.expectedTokens{
			token := l.NextToken()
			//t.Logf("Token Type : %q and Token literal : %q", token.Type, token.Literal)
			if token.Type != tt.expectedType {
				t.Fatalf("Test [%d] - token type wrong. expected = %q, got = %q", i, tt.expectedType, token.Type)
			}
			if token.Literal != tt.expectedLiteral {
				t.Fatalf("Test [%d] - token literal wrong. expected = %q, got = %q", i, tt.expectedLiteral, token.Literal)
			}
		}
	}

}