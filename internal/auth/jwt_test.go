package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {

	id1 := uuid.New()
	id2 := uuid.New()
	secret1 := "secret1"
	secret2 := "secret2"
	exp1 := time.Duration(1000000000)
	exp2 := time.Duration(-2)
	ss1, _ := MakeJWT(id1, secret1, exp1)
	ss2, _ := MakeJWT(id2, secret2, exp2)

	tests := []struct {
		name    string
		secret  string
		ss      string
		wantErr bool
	}{
		{
			name:    "Correct jwt",
			secret:  secret1,
			ss:      ss1,
			wantErr: false,
		},
		{
			name:    "expired",
			secret:  secret2,
			ss:      ss2,
			wantErr: true,
		},
		{
			name:    "wrong secret",
			secret:  secret2,
			ss:      ss1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.ss, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
