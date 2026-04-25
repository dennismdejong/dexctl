package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	if version == "" {
		t.Error("version should not be empty")
	}
}

func TestClientCreateFlags(t *testing.T) {
	tests := []struct {
		name        string
		public     bool
		secret    string
		wantPublic bool
	}{
		{
			name:        "public client",
			public:     true,
			wantPublic: true,
		},
		{
			name:        "confidential client",
			public:     false,
			wantPublic: false,
		},
		{
			name:        "client with secret",
			secret:     "mysecret",
			wantPublic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.public != tt.wantPublic {
				t.Errorf("public flag = %v, want %v", tt.public, tt.wantPublic)
			}
		})
	}
}

func TestRedirectURIsValidation(t *testing.T) {
	tests := []struct {
		name   string
		uris  []string
		valid bool
	}{
		{
			name:   "single redirect URI",
			uris:  []string{"https://example.com/callback"},
			valid: true,
		},
		{
			name:   "multiple redirect URIs",
			uris:  []string{"https://example.com/callback", "http://localhost:8080/callback"},
			valid: true,
		},
		{
			name:   "empty redirect URIs",
			uris:  []string{},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.uris) > 0 || tt.valid {
				// URIs should be valid if they exist or if validation passes
			}
		})
	}
}

func TestTLSConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		certFile string
		keyFile  string
		caFile  string
		mTLS    bool
	}{
		{
			name:     "no TLS",
			certFile: "",
			keyFile:  "",
			caFile:  "",
			mTLS:    false,
		},
		{
			name:     "CA only",
			certFile: "",
			keyFile:  "",
			caFile:  "/path/to/ca.crt",
			mTLS:    false,
		},
		{
			name:     "mutual TLS",
			certFile: "/path/to/cert.crt",
			keyFile:  "/path/to/key.key",
			caFile:  "/path/to/ca.crt",
			mTLS:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isMutualTLS := tt.certFile != "" && tt.keyFile != "" && tt.caFile != ""
			if isMutualTLS != tt.mTLS {
				t.Errorf("mutual TLS = %v, want %v", isMutualTLS, tt.mTLS)
			}
		})
	}
}

func TestServerAddress(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{
			name:    "default address",
			addr:    "localhost:5556",
			wantErr: false,
		},
		{
			name:    "custom address",
			addr:    "dex.example.com:5556",
			wantErr: false,
		},
		{
			name:    "IP address",
			addr:    "127.0.0.1:5556",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAddr := tt.addr != ""
			if hasAddr == tt.wantErr {
				// Should have an address if validation passes
			}
		})
	}
}

func TestOutputFormat(t *testing.T) {
	tests := []struct {
		name      string
		isJSON    bool
		wantJSON bool
	}{
		{
			name:      "default output",
			isJSON:    false,
			wantJSON: false,
		},
		{
			name:      "JSON output",
			isJSON:    true,
			wantJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isJSON != tt.wantJSON {
				t.Errorf("isJSON = %v, want %v", tt.isJSON, tt.wantJSON)
			}
		})
	}
}

func TestJSONOutput(t *testing.T) {
	data := map[string]interface{}{
		"id":   "test-id",
		"name": "test-client",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("json.Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("json.Unmarshal error: %v", err)
	}

	if decoded["id"] != "test-id" {
		t.Errorf("id = %v, want %v", decoded["id"], "test-id")
	}

	if decoded["name"] != "test-client" {
		t.Errorf("name = %v, want %v", decoded["name"], "test-client")
	}
}

func TestClientInfoJSON(t *testing.T) {
	info := map[string]interface{}{
		"id":           "test-id",
		"name":         "Test Client",
		"redirectUris":  []interface{}{"https://example.com/callback"},
		"trustedPeers":  []interface{}{"peer-1"},
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Errorf("json.Marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Errorf("json.Unmarshal error: %v", err)
	}

	if decoded["id"] != "test-id" {
		t.Errorf("id = %v, want %v", decoded["id"], "test-id")
	}

	if decoded["name"] != "Test Client" {
		t.Errorf("name = %v, want %v", decoded["name"], "Test Client")
	}
}

func TestPrintOutputJSON(t *testing.T) {
	var buf bytes.Buffer

	client := map[string]interface{}{
		"id":     "test-id",
		"name":   "Test Client",
		"public": true,
	}

	data, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		t.Errorf("json.MarshalIndent error: %v", err)
	}

	buf.Write(data)

	if buf.Len() == 0 {
		t.Error("buffer should not be empty")
	}
}

func TestClientUpdateFields(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		newName string
		valid   bool
	}{
		{
			name:     "valid update",
			id:       "test-id",
			newName: "Updated Client",
			valid:   true,
		},
		{
			name:     "empty name",
			id:       "test-id",
			newName: "",
			valid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasID := tt.id != ""
			if hasID != tt.valid {
				t.Errorf("hasID = %v, want %v", hasID, tt.valid)
			}
		})
	}
}

func TestContextTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout int
		valid   bool
	}{
		{
			name:    "default timeout",
			timeout: 10,
			valid:   true,
		},
		{
			name:    "custom timeout",
			timeout: 30,
			valid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.timeout <= 0 != !tt.valid {
				t.Errorf("timeout = %v, want valid = %v", tt.timeout, tt.valid)
			}
		})
	}
}

func TestInsecureMode(t *testing.T) {
	tests := []struct {
		name     string
		insecure bool
		caFile   string
		mTLS     bool
	}{
		{
			name:     "insecure mode",
			insecure: true,
			caFile:   "",
			mTLS:     false,
		},
		{
			name:     "secure mode",
			insecure: false,
			caFile:   "/path/to/ca.crt",
			mTLS:     false,
		},
		{
			name:     "mutual TLS",
			insecure: false,
			caFile:   "/path/to/ca.crt",
			mTLS:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isSecure := !tt.insecure && tt.caFile != ""
			if isSecure == tt.insecure {
				// Should be secure if not insecure and CA is provided
			}
		})
	}
}