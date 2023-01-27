package user

import (
	"testing"
)

func Test_PasswordHasher(t *testing.T) {
	tests := []struct {
		password string
		hash     string
	}{
		{password: "test", hash: "dce472b679aa4d3893d3166dee95725a"},
		{password: "demo", hash: "530902b46e256b2c066ddeec2273f400"},
		{password: "check", hash: "46235322ee994e9464bd84cbe04257b2"},
		{password: "password", hash: "1f0d3bce3d532a608138c96deadfc498"},
	}
	for _, tt := range tests {
		if got := PasswordHasher(tt.password); got != tt.hash {
			t.Errorf("PasswordHasher() = %v, want %v", got, tt.hash)
		}
	}
}

func Test_ValidateOrderNumber(t *testing.T) {
	tests := []struct {
		number int
		check  bool
	}{
		{number: 102543287, check: true},
		{number: 7813512, check: true},
		{number: 87643, check: true},
		{number: 814946751236548, check: true},
		{number: 74586128954123564, check: true},
		{number: 7500, check: true},
		{number: 54320, check: true},
		{number: 478654151, check: false},
		{number: 7891458, check: false},
	}
	for _, tt := range tests {
		if got := ValidateOrderNumber(tt.number); got != tt.check {
			t.Errorf("ValidateOrderNumber() = %v, want %v", got, tt.check)
		}
	}
}

func Test_luhn(t *testing.T) {
	tests := []struct {
		number int
		code   int
	}{
		{number: 45645645, code: 2},
		{number: 10860419, code: 0},
		{number: 851126540, code: 0},
		{number: 789453125, code: 5},
		{number: 74511320, code: 1},
		{number: 123475485, code: 1},
		{number: 4186862, code: 9},
		{number: 715469813, code: 7},
		{number: 50974104, code: 3},
	}
	for _, tt := range tests {
		if got := luhn(tt.number); got != tt.code {
			t.Errorf("luhn() = %v, want %v", got, tt.code)
		}
	}
}
