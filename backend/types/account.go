package types

type UserAccountAuthorizationData struct {
	Username     string
	PasswordHash string
}

type ModeratorAccountAuthorizationData struct {
	Username     string
	PasswordHash string
}
