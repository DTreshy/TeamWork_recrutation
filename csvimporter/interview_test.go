package csvimporter_test

import (
	"testing"

	"github.com/DTreshy/TeamWork_recrutation/csvimporter"
)

func TestIsValidHostname(t *testing.T) {
	var endpointTestCases = []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid domain",
			input:   "example.com",
			wantErr: false,
		},
		{
			name:    "invalid domain name",
			input:   "invalid..domain",
			wantErr: true,
		},
		{
			name:    "empty endpoint",
			input:   "",
			wantErr: true,
		},
		{
			name:    "maximum length domain label",
			input:   "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.com",
			wantErr: false,
		},
		{
			name:    "maximum length plus one domain label",
			input:   "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl.com",
			wantErr: true,
		},
		{
			name:    "maximum length domain name",
			input:   "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghi",
			wantErr: false,
		},
		{
			name:    "maximum length domain name plus one",
			input:   "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghij",
			wantErr: true,
		},
		{
			name:    "domain name with invalid character",
			input:   "example$domain.com",
			wantErr: true,
		},
		{
			name:    "label with invalid rune",
			input:   "example\uFFFDdomain.com",
			wantErr: true,
		},
		{
			name:    "label with hyphen on the beginning",
			input:   "example.-com.invalid",
			wantErr: true,
		},
		{
			name:    "label with hyphen on the end",
			input:   "example.com-.invalid",
			wantErr: true,
		},
		{
			name:    "domain ending with hyphen",
			input:   "example-",
			wantErr: true,
		},
		{
			name:    "domain ending with a dot",
			input:   "example.",
			wantErr: true,
		},
	}

	for _, tt := range endpointTestCases {
		t.Run(tt.name, func(t *testing.T) {
			err := csvimporter.IsValidHostname(tt.input)
			if err != nil != tt.wantErr {
				t.Errorf("Expected error other than: %v", err)
			}
		})
	}
}
