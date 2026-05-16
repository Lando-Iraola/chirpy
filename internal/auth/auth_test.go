package auth

import (
	"fmt"
	"net/http"
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

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestValidateBearerToken(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)
	goodHeader := http.Header{}
	goodHeader.Add("Authorization", fmt.Sprintf("Bearer %s", validToken))

	badHeader := http.Header{}
	badHeader.Add("Authorization", "")

	badHeader2 := http.Header{}
	badHeader2.Add("", "")
	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		bearerToken http.Header
		wantErr     bool
	}{
		{
			name:        "Valid bearer token",
			tokenString: validToken,
			tokenSecret: "secret",
			bearerToken: goodHeader,
			wantErr:     false,
		},
		{
			name:        "Bad header",
			tokenString: "",
			tokenSecret: "secret",
			bearerToken: badHeader,
			wantErr:     true,
		},
		{
			name:        "Missing header",
			tokenString: "",
			tokenSecret: "secret",
			bearerToken: badHeader2,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bToken, err := GetBearerToken(tt.bearerToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bToken != tt.tokenString {
				t.Errorf("GetBearerToken() Bearer Token = %v, want %v", bToken, tt.tokenString)
			}

			if bToken != "" {
				user, err := ValidateJWT(bToken, "secret")
				if err != nil {
					t.Errorf("GetBearerToken() Bearer Token = %v, fails validation %v", bToken, err)
				}

				if user != userID {
					t.Errorf("GetBearerToken() Bearer Token = %v, gives user %v, expected: %v", bToken, user, userID)
				}
			}
		})
	}
}
