package storage

import (
	"context"
	"time"
)

// Client represents an OAuth client
type Client struct {
	ID            string    `json:"id"`
	Secret        string    `json:"secret"`
	Name          string    `json:"name"`
	RedirectURIs  []string  `json:"redirect_uris"`
	GrantTypes    []string  `json:"grant_types"`
	ResponseTypes []string  `json:"response_types"`
	Scope         string    `json:"scope"`
	CreatedAt     time.Time `json:"created_at"`
}

// AuthorizationCode represents an authorization code
type AuthorizationCode struct {
	Code                string                 `json:"code"`
	ClientID            string                 `json:"client_id"`
	RedirectURI         string                 `json:"redirect_uri"`
	Scope               string                 `json:"scope"`
	CodeChallenge       string                 `json:"code_challenge"`
	CodeChallengeMethod string                 `json:"code_challenge_method"`
	SnowflakeToken      map[string]interface{} `json:"snowflake_token,omitempty"`
	ExpiresAt           time.Time              `json:"expires_at"`
	State               string                 `json:"state"`
}

// AccessToken represents an access token
type AccessToken struct {
	Token     string    `json:"token"`
	ClientID  string    `json:"client_id"`
	Scope     string    `json:"scope"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	Token       string    `json:"token"`
	ClientID    string    `json:"client_id"`
	AccessToken string    `json:"access_token"`
	Scope       string    `json:"scope"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// Store defines the interface for storage operations
type Store interface {
	// Connection management
	Connect(ctx context.Context, cfg Config) error
	Disconnect(ctx context.Context) error
	IsHealthy(ctx context.Context) bool

	// Client operations
	GetClientByNameAndRedirectURIs(ctx context.Context, name string, redirectURIs []string) (*Client, error)
	StoreClient(ctx context.Context, client *Client) error
	GetClient(ctx context.Context, clientID string) (*Client, error)

	// Authorization code operations
	StoreAuthorizationCode(ctx context.Context, code *AuthorizationCode) error
	GetAuthorizationCode(ctx context.Context, code string) (*AuthorizationCode, error)
	UpdateAuthorizationCodeToken(ctx context.Context, code string, token map[string]interface{}) error
	DeleteAuthorizationCode(ctx context.Context, code string) error

	// Access token operations
	StoreAccessToken(ctx context.Context, token *AccessToken) error
	GetAccessToken(ctx context.Context, token string) (*AccessToken, error)
	DeleteAccessToken(ctx context.Context, token string) error

	// Refresh token operations
	StoreRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error

	// Status
	GetStatus(ctx context.Context) map[string]interface{}
}

// Config holds storage configuration
type Config struct {
	Host           string
	Port           int
	Database       string
	Username       string
	Password       string
	PoolSize       int
	MaxConnections int
}
