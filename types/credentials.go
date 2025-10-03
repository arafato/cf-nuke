package types

type Credentials struct {
	APIKey         string
	AccountID      string
	User           string
	Mode           Mode
	S3AccessKeyID  string
	S3AccessSecret string
}

type Mode string

const (
	Account Mode = "account"
	Token   Mode = "token"
)
