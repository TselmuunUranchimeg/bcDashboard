package services

import (
	"fmt"
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	hash, err := GenerateHash("password")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hash)
	result, err := VerifyPassword("password1", hash)
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Fatalf("they are not supposed to match\n")
	}
}

func TestGenerateJwtToken(t *testing.T) {
	token, err := GenerateJwtToken("Tselmuun")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
}

func TestVerifyToken(t *testing.T) {
	token, err := VerifyToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjQwMDAiLCJpYXQiOjE3MDQzMzEzMTYsImV4cCI6MTcwNDMzMTM3NiwiYXVkIjoiW2h0dHA6Ly9sb2NhbGhvc3Q6NTE3M10iLCJzdWIiOiIiLCJVc2VybmFtZSI6IlRzZWxtdXVuIn0.JGjYliSjyeF6He2qonDzgVJN0TmRbsjT6STWtynFwhc")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
}
