package hakucho

import (
	"fmt"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		credPath  string
		tokenPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no cred and token", args{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.credPath, tt.args.tokenPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func fileExists(name string) bool {
	if f, err := os.Stat(name); os.IsNotExist(err) || f.IsDir() {
		return false
	}
	return true
}

func newTestClient(t *testing.T) (*Client, error) {
	t.Helper()
	if !fileExists("credentials/credentials.json") || !fileExists("credentials/token.json") {
		t.Skip("no credential files")
		return nil, fmt.Errorf("no credentials")
	}
	c, err := NewClient("credentials/credentials.json", "credentials/token.json")
	if err != nil {
		t.Errorf("client init error: %s", err)
		return nil, err
	}
	return c, nil
}
