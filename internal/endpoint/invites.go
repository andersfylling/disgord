package endpoint

// Invites /invites
func Invites() string {
	return invites
}

// Invite /invites/{invite.code}
func Invite(code string) string {
	return Invites() + "/" + code
}
