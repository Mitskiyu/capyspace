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
