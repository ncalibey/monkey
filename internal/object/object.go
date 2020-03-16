// Package object contains the code for object representation during AST evaluation.
package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/ncalibey/monkey/internal/ast"
)

// ObjectType details what the underlying Type of the object is.
type ObjectType string

// BuiltinFunction represents functions build into the language.
type BuiltinFunction func(args ...Object) Object

// HashKey is the internal value used for mapping literals as keys within a hash.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

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
	HASH_OBJ         = "HASH"
)

// Hashable is implemented by Objects that can be used as hash keys.
type Hashable interface {
	HashKey() HashKey
}

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

// String implements the Object and hashable interfaces for string values.
type String struct {
	// Value is the value assocated with the String.
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Integer implements the Object and Hash interfaces for integer values.
type Integer struct {
	// Value is the value associated with the Integer.
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: uint64(i.Value)} }

// Boolean implements the Object and Hashable interfaces for boolean values.
type Boolean struct {
	// Value is the value associated with the Boolean.
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

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

// HashPair is a key-value pair that is associated with a Hash.
type HashPair struct {
	Key   Object
	Value Object
}

// Array implements the Object interface for hashes.
type Hash struct {
	// Pairs are the key-value pairs associated with the Hash.
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf(":%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

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
