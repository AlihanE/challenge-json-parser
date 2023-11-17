package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/AlihanE/challenge-json-parser/ast"
	"github.com/AlihanE/challenge-json-parser/lexer"
	"github.com/AlihanE/challenge-json-parser/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	p.lexer.NextToken()
	p.lexer.NextToken()

	return p
}

func (p *Parser) currentTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) parseValue() ast.Value {
	switch p.currentToken.Type {
	case token.LeftBrace:
		return p.parseJSONObject()
	case token.LeftBracket:
		return p.parseJSONArray()
	default:
		return p.parseJSONLiteral()
	}
}

func (p *Parser) ParseProgram() (ast.RootNode, error) {
	var rootNode ast.RootNode
	if p.currentTokenIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
	}

	val := p.parseValue()
	if val == nil {
		p.parseError(fmt.Sprintf("Error parsing JSON expected a value, got %v:", p.currentToken.Literal))
		return ast.RootNode{}, errors.New(p.Errors())
	}
	rootNode.RootValue = &val

	return rootNode, nil
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseJSONObject() ast.Value {
	obj := ast.Object{Type: "Object"}
	objState := ast.ObjStart

	for !p.currentTokenIs(token.EOF) {
		switch objState {
		case ast.ObjStart:
			if p.currentTokenIs(token.LeftBrace) {
				objState = ast.ObjOpen
				obj.Start = p.currentToken.Start
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing JSON object Expected `{` token, got: %s", p.currentToken.Literal))
				return nil
			}
		case ast.ObjOpen:
			if p.currentTokenIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.End
				return obj
			}
			prop := p.parseProperty()
			obj.Children = append(obj.Children, prop)
			objState = ast.ObjProperty
		case ast.ObjProperty:
			if p.currentTokenIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.Start
				return obj
			} else if p.currentTokenIs(token.Comma) {
				objState = ast.ObjComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing JSON object Expected `}` or `,` token, got: %s", p.currentToken.Literal))
				return nil
			}
		case ast.ObjComma:
			prop := p.parseProperty()
			if prop.Value != nil {
				obj.Children = append(obj.Children, prop)
				objState = ast.ObjProperty
			}
		}
	}

	obj.End = p.currentToken.Start

	return obj
}

func (p *Parser) parseJSONArray() ast.Array {
	array := ast.Array{Type: "Array"}
	arrayState := ast.ArrayStart
	for !p.currentTokenIs(token.EOF) {
		switch arrayState {
		case ast.ArrayStart:
			if p.currentTokenIs(token.LeftBracket) {
				array.Start = p.currentToken.Start
				arrayState = ast.ArrayOpen
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing JSON object Expected `[` token, got: %s", p.currentToken.Literal))
			}
		case ast.ArrayOpen:
			if p.currentTokenIs(token.RightBracket) {
				array.End = p.currentToken.End
				p.nextToken()
				return array
			}
			val := p.parseValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue
			if p.peekTokenIs(token.RightBracket) {
				p.nextToken()
			}
		case ast.ArrayValue:
			if p.currentTokenIs(token.RightBracket) {
				array.End = p.currentToken.End
				p.nextToken()
				return array
			} else if p.currentTokenIs(token.Comma) {
				arrayState = ast.ArrayComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing JSON object Expected `]` or `,` token, got: %s", p.currentToken.Literal))
			}
		case ast.ArrayComma:
			val := p.parseValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue
		}
	}

	array.End = p.currentToken.Start
	return array
}

func (p *Parser) parseJSONLiteral() ast.Literal {
	val := ast.Literal{Type: "Literal"}

	defer p.nextToken()

	switch p.currentToken.Type {
	case token.String:
		val.Value = p.parseString()
		return val
	case token.Number:
		v, _ := strconv.Atoi(p.currentToken.Literal)
		val.Value = v
		return val
	case token.True:
		val.Value = true
		return val
	case token.False:
		val.Value = false
		return val
	default:
		val.Value = "null"
		return val
	}
}

func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{Type: "Property"}
	propertyState := ast.PropertyStart

	for !p.currentTokenIs(token.EOF) {
		switch propertyState {
		case ast.PropertyStart:
			if p.currentTokenIs(token.String) {
				key := ast.Identifier{
					Type:  "Identifier",
					Value: p.parseString(),
				}
				prop.Key = key
				propertyState = ast.PropertyKey
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing property start. Expected String token, got: %s", p.currentToken.Literal))
			}
		case ast.PropertyKey:
			if p.currentTokenIs(token.Colon) {
				propertyState = ast.PropertyColon
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf("Error parsing property. Expected Colon token, got: %s", p.currentToken.Literal))
			}
		case ast.PropertyColon:
			val := p.parseValue()
			prop.Value = val
			return prop
		}
	}

	return prop
}

func (p *Parser) parseString() string {
	return p.currentToken.Literal
}

func (p *Parser) parseError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() string {
	return strings.Join(p.errors, ", ")
}
