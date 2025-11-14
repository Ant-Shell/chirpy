package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashAndCheckPassword_Success(t *testing.T) {
	pwd := "hunter2"
	hash, err := HashPassword(pwd)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if hash == "" || hash == pwd {
		t.Fatalf("hash should be non-empty and not equal to password")
	}
	ok, err := CheckPasswordHash(pwd, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash error: %v", err)
	}
	if !ok {
		t.Fatalf("expected password to match hash")
	}
}

func TestHashAndCheckPassword_Fail(t *testing.T) {
	hash, err := HashPassword("hunter-the-second-was-here")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	ok, err := CheckPasswordHash("wrong-password", hash)
	if err != nil{
		t.Fatalf("CheckPasswordHash error: %v", err)
	}
	if ok {
		t.Fatalf("expected mismatch for wrong password")
	}
}

func TestMakeAndValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}

	if gotID != userID {
		t.Fatalf("want %s, got %s", userID, gotID)
	}
}

func TestMakeAndValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	token, err := MakeJWT(userID, secret, -time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}
	if _, err := ValidateJWT(token, secret); err == nil {
		t.Fatal("expected erro for expired token, got nil")
	}
}

func TestMakeAndValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	token, err := MakeJWT(userID, "right-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT error: %v", err)
	}
	if _, err := ValidateJWT(token, "wrong-secret"); err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}