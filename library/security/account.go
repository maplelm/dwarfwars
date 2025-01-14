package security

type Account struct {
	Username string
	UID      string
	Email    string
	SSO      string // blank if SSO is not used
}

func NewAccount(user string, uid string, email string, sso string) *Account {
	return &Account{
		Username: user,
		UID:      uid,
		Email:    email,
		SSO:      sso,
	}
}
