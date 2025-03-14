package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHasPassowrd(t *testing.T) {
	cases := []struct {
		val string
	}{
		{
			val: "testdata",
		},
		{
			val: "moretestdata",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			hash, err := HashPassword(c.val)
			if err != nil {
				t.Errorf("hash failed")
			}
			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(c.val))
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		val  string
		hash []byte
	}{
		{
			val: "testdata",
			hash: func() []byte {
				hash, _ := bcrypt.GenerateFromPassword([]byte("testdata"), 10)
				return hash
			}(),
		},
		{
			val: "moretestdata",
			hash: func() []byte {
				hash, _ := bcrypt.GenerateFromPassword([]byte("moretestdata"), 10)
				return hash
			}(),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			err := CheckPasswordHash(c.val, string(c.hash))
			if err != nil {
				t.Errorf("failed checking hash")
			}
		})
	}
}
