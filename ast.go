package main

// AstValue docs here
type AstValue interface{}

// AstObject docs here
type AstObject []AstMember

// AstMember docs here
type AstMember struct {
	key string
	val AstValue
}

// AstArray docs here
type AstArray []AstValue
