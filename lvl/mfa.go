package lvl

// MFA ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-mfa-level
type MFA uint

func (mfal *MFA) None() bool {
	return *mfal == 0
}
func (mfal *MFA) Elevated() bool {
	return *mfal == 1
}
