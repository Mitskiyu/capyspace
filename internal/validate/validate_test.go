package validate_test

import (
	"testing"

	"github.com/Mitskiyu/capyspace/internal/validate"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid simple", "test@example.com", false},
		{"Valid subdomain", "test@sub.example.co.uk", false},
		{"Valid plus alias", "test+alias@example.com", false},
		{"Valid dots in local", "first.last@example.com", false},
		{"Valid numbers", "test1234@example123.com", false},
		{"Valid hyphen domain", "test@example-site.com", false},
		{"Valid quoted local", `"test test"@example.com`, false},

		{"Empty string", "", true},

		{"Missing @", "testexample.com", true},
		{"Missing domain", "test@", true},
		{"Missing local part", "@example.com", true},
		{"Space in local", "test space@example.com", true},
		{"Space in domain", "test@exa mple.com", true},
		{"Multiple @", "test@@example.com", true},
		{"Leading dot local", ".test@example.com", true},
		{"Trailing dot local", "test.@example.com", true},
		{"Consecutive dots local", "test..test@example.com", true},
		{"Plain string", "juststring", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validate.Email(tt.email)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Email() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Email() succeeded unexpectedly")
			}
		})
	}
}

func TestVerificationToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Valid simple", "123456", false},
		{"Valid all zero", "000000", false},
		{"Valid all nine", "999999", false},

		{"Empty string", "", true},

		{"Too short", "12345", true},
		{"Too long", "1234567", true},
		{"Non-digit char", "12a456", true},
		{"Space in token", "12 456", true},
		{"Special char", "12#456", true},
		{"Hyphen in token", "12-456", true},
		{"Alpha only", "abcdef", true},
		{"Mixed alpha", "1a2b3c", true},
		{"Unicode digit", "１２３４５６", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := validate.VerificationToken(tt.token)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("VerificationToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("VerificationToken() succeeded unexpectedly")
			}
		})
	}
}
