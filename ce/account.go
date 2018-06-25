package ce

var (
	AccountsURL = "/accounts"
	AccountIDURLFormat = "/accounts/%s"
	SignupURL = "/signup"
)

// GetAccounts lists all the accounts
func GetAccounts(auth string) error {
	// GET AccountsURL
	return nil
}

// CreateAccount given first, last, email, create a new account
// which sends a verification e-mail
func CreateAccount(first, last, email string) error {
	// POST SignupURL
	return nil
}

// DisableAccount given an accountID, disable it
func DisableAccount(accountID string) error {
	// PATCH AccountIDURLFormat
	return nil
}
