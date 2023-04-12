package user

import (
	"testing"
)

func Test_isValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			"positive case",
			"testemail@mail.com",
			true,
		},
		{
			"without @",
			"testemailmail.com",
			false,
		},
		{
			"without domain",
			"testemail@",
			false,
		},
		{
			"wrong domain",
			"testemail@asd",
			false,
		},
		{
			"small domain",
			"testemail@a.a",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.want {
				t.Errorf("isValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
