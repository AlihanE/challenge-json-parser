package token

import "fmt"

type Type string

const (
	Illegal      Type = "ILLEGAL"
	EOF          Type = "EOF"
	String       Type = "STRING"
	Number       Type = "NUMBER"
	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Comma        Type = ","
	Colon        Type = ":"
	True         Type = "TRUE"
	False        Type = "FALSE"
	Null         Type = "NULL"
)

type Token struct {
	Type    Type
	Literal string
	Start   int
	End     int
	Line    int
}

var validJSONIdentifiers = map[string]Type{
	"true":  True,
	"false": False,
	"null":  Null,
}

func LookupIdentifier(identifier string) (Type, error) {
	if token, ok := validJSONIdentifiers[identifier]; ok {
		return token, nil
	}
	return "", fmt.Errorf("Expected a valid JSON identifier. Found %s", identifier)
}
