package types

type UserAccountAuthorizationData struct {
	Username     string
	Email        string
	PasswordHash string
}

type ModeratorAccountAuthorizationData struct {
	Username     string
	PasswordHash string
}
