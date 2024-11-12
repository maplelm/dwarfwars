package client

type Account struct {
	Username string
	UID      string
	Email    string
	SSO      string // blank if SSO is not used
}
