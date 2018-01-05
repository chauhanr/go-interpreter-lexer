package lexer

import (
	"go-interpreter-lexer/token"
)

type Lexer struct{
	input string
	position int
	readPosition int
	ch byte
}

func New(input string) *Lexer{
	l := &Lexer{ input: input }
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token{
	var t  token.Token
	l.skipWhitespace()

	switch l.ch {
			case '=' :
				if l.peekChar() == '=' {
					ch := l.ch
					l.readChar()
					literal := string(ch) + string(l.ch)
					t = token.Token{token.EQ, literal}
				}else{
					t = newToken(token.ASSIGN, l.ch)
				}
			case ';' :
				t = newToken(token.SEMICOLON, l.ch)
			case '(' :
				t = newToken(token.LPAREN, l.ch)
			case ')' :
				t = newToken(token.RPAREN, l.ch)
			case '{' :
				t = newToken(token.LBRACE, l.ch)
			case '}' :
				t = newToken(token.RBRACE, l.ch)
			case ',' :
				t = newToken(token.COMMA, l.ch)
			case '+' :
				t = newToken(token.PLUS, l.ch)
			case '!' :
				if l.peekChar() == '=' {
					ch := l.ch
					l.readChar()
					literal := string(ch) + string(l.ch)
					t = token.Token{token.NOT_EQ, literal}
				}else {
					t = newToken(token.BANG, l.ch)
				}
			case '/' :
				t = newToken(token.SLASH, l.ch)
			case '*' :
				t = newToken(token.ASTERISK, l.ch)
			case '>' :
				t = newToken(token.GT, l.ch)
			case '<' :
				t = newToken(token.LT, l.ch)
			case '-' :
				t = newToken(token.MINUS, l.ch)
			case 0 :
				t = token.Token{token.EOF, ""}
			case '"':
				t.Type = token.STRING
				t.Literal = l.readString()
			default:
				if isLetter(l.ch){
					t.Literal = l.readIdentifier()
					t.Type = token.LookupIdent(t.Literal)
					return t
				}else if isDigit(l.ch) {
					t.Literal = l.readNumber()
					t.Type = token.INT
					return t
				} else{
					t = newToken(token.ILLEGAL, l.ch)
				}
				}

		l.readChar()
	return t
}

/** keep reading until we find the end of the string " or EOF */
func (l *Lexer) readString() string{
	pos := l.position+1
	for{
		//fmt.Printf("char %s", l.ch)
		l.readChar()
		if l.ch == '"' || l.ch == 0{
			break
		}
	}
	return l.input[pos:l.position]
}

func isDigit(ch byte) bool{
	return ch >= '0' && ch <= '9'
}


func isLetter(ch byte) bool{
	//log.Printf("character for letter %s\n", string(ch))
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_';
}


func (l *Lexer) readNumber() string{
	position := l.position
	for isDigit(l.ch){
		l.readChar()
	}
	return l.input[position:l.position]
}

/*
	method will read identifiers in the code base.
*/
func (l *Lexer) readIdentifier() string{
	position := l.position
	for isLetter(l.ch){
		l.readChar()
	}
	return l.input[position:l.position]
}

// avoid all the spaces, tabs, carriage returns etc.
func (l *Lexer) skipWhitespace(){
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r'{
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token{
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readChar(){
	if l.readPosition >= len(l.input){
		l.ch = 0
	}else{
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition +=1
}

func (l *Lexer) peekChar() byte{
	if l.readPosition >= len(l.input){
		return 0
	}else {
		return l.input[l.readPosition]
	}
}