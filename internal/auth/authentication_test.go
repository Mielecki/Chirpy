package auth

import "testing"


func TestHashPassword(t *testing.T) {
	password := "password1234"
	
	hash, err := HashPassword(password)
	
	if err != nil || hash == "" || hash == password {
		t.Fatalf("HashPassword failed")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password1234"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed")
	}

	if CheckPasswordHash(password, hash) != nil {
		t.Fatalf("CheckPasswordHash failed")
	}
}

