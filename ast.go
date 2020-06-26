package main

// AstValue is a JSON value
type AstValue interface{}

// AstObject is a JSON object
type AstObject []AstMember

// AstMember is JSON member
type AstMember struct {
	key string
	val AstValue
}

// AstArray is a JSON array
type AstArray []AstValue
