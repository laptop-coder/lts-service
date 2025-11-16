package types

type ModeratorAccount struct {
	ModeratorId        int64
	Username           string
	Email              Email
	PasswordHash       string
	CredentialsVersion int
}
