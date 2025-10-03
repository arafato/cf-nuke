package types

type Credentials struct {
	APIKey    string
	AccountID string
	User      string
	Mode      Mode
}

type Mode string

const (
	Account Mode = "account"
	Token   Mode = "token"
)
