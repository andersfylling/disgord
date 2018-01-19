package lvl

// Verification ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-verification-level
type Verification uint

// None unrestricted
func (vl *Verification) None() bool {
	return *vl == 0
}

// Low must have verified email on account
func (vl *Verification) Low() bool {
	return *vl == 1
}

// Medium must be registered on Discord for longer than 5 minutes
func (vl *Verification) Medium() bool {
	return *vl == 2
}

// High (╯°□°）╯︵ ┻━┻ - must be a member of the server for longer than 10 minutes
func (vl *Verification) High() bool {
	return *vl == 3
}

// VeryHigh ┻━┻ミヽ(ಠ益ಠ)ﾉ彡┻━┻ - must have a verified phone number
func (vl *Verification) VeryHigh() bool {
	return *vl == 4
}
