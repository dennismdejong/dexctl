package client

import (
	"testing"

	"github.com/dexidp/dex/api/v2"
)

func TestCreateClientRequest(t *testing.T) {
	tests := []struct {
		name        string
		clientName string
		public    bool
		wantName  string
		wantPublic bool
	}{
		{
			name:        "public client",
			clientName: "TestPublic",
			public:    true,
			wantName:  "TestPublic",
			wantPublic: true,
		},
		{
			name:        "confidential client",
			clientName: "TestConfidential",
			public:    false,
			wantName:  "TestConfidential",
			wantPublic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &api.CreateClientReq{
				Client: &api.Client{
					Name:   tt.clientName,
					Public: tt.public,
				},
			}
			if req.GetClient().GetName() != tt.wantName {
				t.Errorf("GetClient().GetName() = %v, want %v", req.GetClient().GetName(), tt.wantName)
			}
			if req.GetClient().GetPublic() != tt.wantPublic {
				t.Errorf("GetClient().GetPublic() = %v, want %v", req.GetClient().GetPublic(), tt.wantPublic)
			}
		})
	}
}

func TestGetClientRequest(t *testing.T) {
	tests := []struct {
		name    string
		id     string
		wantID string
	}{
		{
			name:    "valid ID",
			id:     "test-client-id",
			wantID: "test-client-id",
		},
		{
			name:    "empty ID",
			id:     "",
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &api.GetClientReq{
				Id: tt.id,
			}
			if req.GetId() != tt.wantID {
				t.Errorf("GetId() = %v, want %v", req.GetId(), tt.wantID)
			}
		})
	}
}

func TestUpdateClientRequest(t *testing.T) {
	tests := []struct {
		name         string
		id          string
		newName    string
		wantID     string
		wantName  string
	}{
		{
			name:     "update name",
			id:      "test-id",
			newName: "UpdatedName",
			wantID: "test-id",
			wantName: "UpdatedName",
		},
		{
			name:     "update with empty name",
			id:      "test-id",
			newName: "",
			wantID: "test-id",
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &api.UpdateClientReq{
				Id:   tt.id,
				Name: tt.newName,
			}
			if req.GetId() != tt.wantID {
				t.Errorf("GetId() = %v, want %v", req.GetId(), tt.wantID)
			}
			if req.GetName() != tt.wantName {
				t.Errorf("GetName() = %v, want %v", req.GetName(), tt.wantName)
			}
		})
	}
}

func TestDeleteClientRequest(t *testing.T) {
	tests := []struct {
		name    string
		id     string
		wantID string
	}{
		{
			name:    "valid ID",
			id:     "test-client-id",
			wantID: "test-client-id",
		},
		{
			name:    "empty ID",
			id:     "",
			wantID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &api.DeleteClientReq{
				Id: tt.id,
			}
			if req.GetId() != tt.wantID {
				t.Errorf("GetId() = %v, want %v", req.GetId(), tt.wantID)
			}
		})
	}
}

func TestListClientRequest(t *testing.T) {
	req := &api.ListClientReq{}
	if req == nil {
		t.Error("ListClientReq should not be nil")
	}
}

func TestClientRedirectURIs(t *testing.T) {
	tests := []struct {
		name        string
		redirectURIs []string
		wantCount  int
	}{
		{
			name:          "single URI",
			redirectURIs: []string{"https://example.com/callback"},
			wantCount:  1,
		},
		{
			name:          "multiple URIs",
			redirectURIs: []string{"https://example.com/callback", "http://localhost:8080/callback"},
			wantCount:  2,
		},
		{
			name:          "no URIs",
			redirectURIs: []string{},
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &api.Client{
				RedirectUris: tt.redirectURIs,
			}
			if len(client.GetRedirectUris()) != tt.wantCount {
				t.Errorf("len(GetRedirectUris()) = %v, want %v", len(client.GetRedirectUris()), tt.wantCount)
			}
		})
	}
}

func TestClientTrustedPeers(t *testing.T) {
	tests := []struct {
		name          string
		trustedPeers []string
		wantCount   int
	}{
		{
			name:          "single trusted peer",
			trustedPeers: []string{"peer-1"},
			wantCount:   1,
		},
		{
			name:          "multiple trusted peers",
			trustedPeers: []string{"peer-1", "peer-2"},
			wantCount:   2,
		},
		{
			name:          "no trusted peers",
			trustedPeers: []string{},
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &api.Client{
				TrustedPeers: tt.trustedPeers,
			}
			if len(client.GetTrustedPeers()) != tt.wantCount {
				t.Errorf("len(GetTrustedPeers()) = %v, want %v", len(client.GetTrustedPeers()), tt.wantCount)
			}
		})
	}
}

func TestClientLogoUrl(t *testing.T) {
	tests := []struct {
		name     string
		logoURL string
		wantURL string
	}{
		{
			name:     "with logo URL",
			logoURL: "https://example.com/logo.png",
			wantURL: "https://example.com/logo.png",
		},
		{
			name:     "empty logo URL",
			logoURL: "",
			wantURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &api.Client{
				LogoUrl: tt.logoURL,
			}
			if client.GetLogoUrl() != tt.wantURL {
				t.Errorf("GetLogoUrl() = %v, want %v", client.GetLogoUrl(), tt.wantURL)
			}
		})
	}
}

func TestClientSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret string
		want   string
	}{
		{
			name:    "with secret",
			secret: "my-secret",
			want:   "my-secret",
		},
		{
			name:    "empty secret",
			secret: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &api.Client{
				Secret: tt.secret,
			}
			if client.GetSecret() != tt.want {
				t.Errorf("GetSecret() = %v, want %v", client.GetSecret(), tt.want)
			}
		})
	}
}

func TestClientInfoFields(t *testing.T) {
	info := &api.ClientInfo{
		Id:           "test-id",
		Name:         "Test Client",
		RedirectUris: []string{"https://example.com/callback"},
		TrustedPeers: []string{"peer-1"},
	}

	if info.GetId() != "test-id" {
		t.Errorf("GetId() = %v, want %v", info.GetId(), "test-id")
	}

	if info.GetName() != "Test Client" {
		t.Errorf("GetName() = %v, want %v", info.GetName(), "Test Client")
	}

	if len(info.GetRedirectUris()) != 1 {
		t.Errorf("len(GetRedirectUris()) = %v, want %v", len(info.GetRedirectUris()), 1)
	}

	if len(info.GetTrustedPeers()) != 1 {
		t.Errorf("len(GetTrustedPeers()) = %v, want %v", len(info.GetTrustedPeers()), 1)
	}
}

func TestUpdateClientResponse(t *testing.T) {
	tests := []struct {
		name      string
		notFound bool
	}{
		{
			name:      "client found",
			notFound: false,
		},
		{
			name:      "client not found",
			notFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &api.UpdateClientResp{}
			if resp.GetNotFound() != false {
				t.Errorf("GetNotFound() = %v, want %v", resp.GetNotFound(), false)
			}
		})
	}
}

func TestDeleteClientResponse(t *testing.T) {
	tests := []struct {
		name      string
		notFound bool
	}{
		{
			name:      "client deleted",
			notFound: false,
		},
		{
			name:      "client not found",
			notFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &api.DeleteClientResp{}
			if resp.GetNotFound() != false {
				t.Errorf("GetNotFound() = %v, want %v", resp.GetNotFound(), false)
			}
		})
	}
}

func TestCreateClientResponse(t *testing.T) {
	resp := &api.CreateClientResp{}
	if resp.GetAlreadyExists() != false {
		t.Errorf("GetAlreadyExists() = %v, want %v", resp.GetAlreadyExists(), false)
	}
	if resp.GetClient() != nil {
		t.Error("GetClient() should be nil for empty response")
	}
}