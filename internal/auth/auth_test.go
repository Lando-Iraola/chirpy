package auth

import (
	"fmt"
	"testing"
)

func TestAuthHashCompare(t *testing.T) {
	const knownHash = "$argon2id$v=19$m=65536,t=1,p=24$gr5WSDgmedaR7OTmXqcBIA$NFwJ6qsYnnL9o+NPpATjJOEcMGKntrUN7HlnH8vKPmM"

	isSame, err := CheckPasswordHash("My beloved Go can't possible be this cute!!", knownHash)
	if err != nil {
		t.Error(err)
	}
	if !isSame {
		t.Error("Somehow, the passwords are different!")
	}
}

func TestAuthHashCreate(t *testing.T) {
	hash, err := HashPassword("My beloved Go can't possible be this cute!!")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("hash: %s", hash)
}

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}
