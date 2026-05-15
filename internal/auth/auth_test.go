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
