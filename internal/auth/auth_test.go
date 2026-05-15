package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestJWT(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	userId1, err := uuid.NewUUID()

	if err != nil {
		t.Error(err)
	}
	token1, err := MakeJWT(userId1, "first secret", time.Duration(10*time.Minute))
	if err != nil {
		t.Error(err)
	}

	userId2, err := uuid.NewUUID()
	if err != nil {
		t.Error(err)
	}

	token2, err := MakeJWT(userId2, "secondo secret", time.Duration(-10*time.Minute))
	if err != nil {
		t.Error(err)
	}

	userId3, err := uuid.NewUUID()
	if err != nil {
		t.Error(err)
	}

	token3, err := MakeJWT(userId3, "right password", time.Duration(10*time.Minute))
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name        string
		userId      uuid.UUID
		tokenSecret string
		token       string
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			userId:      userId1,
			tokenSecret: "first secret",
			token:       token1,
			wantErr:     false,
		},
		{
			name:        "Stale Token",
			userId:      userId2,
			tokenSecret: "secondo secret",
			token:       token2,
			wantErr:     true,
		},
		{
			name:        "Wrong Secret",
			userId:      userId3,
			tokenSecret: "wrong password",
			token:       token3,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := ValidateJWT(tt.token, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.userId {
				t.Errorf("ValidateJWT() expects %v, got %v", tt.userId, match)
			}
		})
	}
}
