// Package object contains the code for object representation during AST evaluation.
package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ncalibey/monkey/internal/ast"
)

// ObjectType details what the underlying Type of the object is.
type ObjectType string

// BuiltinFunction represents functions build into the language.
type BuiltinFunction func(args ...Object) Object

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
)

// Object represents a value that is being evaluated.
type Object interface {
	// Type returns the ObjectType of the current Object.
	Type() ObjectType
	// Inspect returns a string representation of the Object's value.
	Inspect() string
}

// Builtin implements the Object interface for built in functions.
type Builtin struct {
	// Fn is the associated BuiltinFunction
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// String implements the Object interface for string values.
type String struct {
	// Value is the value assocated with the String.
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Integer implements the Object interface for integer values.
type Integer struct {
	// Value is the value associated with the Integer.
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Boolean implements the Object interface for boolean values.
type Boolean struct {
	// Value is the value associated with the Boolean.
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// Null implements the Object interface for null values.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ReturnValue implements the Object interface for return values.
type ReturnValue struct {
	// Value is the value associated with the ReturnValue.
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

// Error implements the Object inerface for errors.
type Error struct {
	// Message is the error message of the Error.
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function implements the Object interface for functions.
type Function struct {
	// Parameters are the parameters of the function.
	Parameters []*ast.Identifier
	// Body is the function body.
	Body *ast.BlockStatement
	// Env is the environment the function carries.
	Env *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Array implements the Object interface for arrays.
type Array struct {
	// Elements are the elements of an array.
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Environment tracks assigned values during evaluation.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnclosedEnvironment returns a new *Environment instance that encloses the argument
// environment.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// NewEnvironment returns a new *Environment instance.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Get returns the object from the store associated with the passed in argument.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set inserts the object into the store using name as the key value.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
