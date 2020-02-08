package ast

import (
	"bytes"

	"github.com/ncalibey/monkey/internal/token"
)

// Node represents a node within an AST.
type Node interface {
	// TokenLiteral is the string literal of the Node's token.
	TokenLiteral() string
	// String returns a string representation of the Node.
	String() string
}

// Statement represents a statement within an AST.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression within an AST.
type Expression interface {
	Node
	expressionNode()
}

// ExpressionStatement represents an expression statement within an AST.
type ExpressionStatement struct {
	// Token is the first token of the expression
	Token token.Token
	// Expression is the expression itself.
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// Program is the root node of an AST.
type Program struct {
	// Statements is a slice of all Statements within a monkey-produced AST.
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement represents a let statement within an AST.
type LetStatement struct {
	// Token is the token.LET token.
	Token token.Token
	// Name is the Identifier of the binding.
	Name *Identifier
	// Value is the expression bound to the Identifier.
	Value Expression
}

func (l *LetStatement) statementNode()       {}
func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String())
	out.WriteString(" = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// Identifier represents an identifier within an AST.
type Identifier struct {
	// Token is the token.IDENT token.
	Token token.Token
	// Value is the expression bound to the Identifier.
	Value string
}

func (i *Identifier) statementNode()       {}
func (i *Identifier) expressionNode() 	   {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// ReturnStatement represents a return statement wthin an AST.
type ReturnStatement struct {
	// Token is the token.RETURN token.
	Token token.Token
	// ReturnValue is the returned expression.
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
