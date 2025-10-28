package auth

import "testing"

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