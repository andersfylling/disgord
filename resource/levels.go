package resource

// ExplicitContentFilterLvl ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
type ExplicitContentFilterLvl uint

func (ecfl *ExplicitContentFilterLvl) Disabled() bool {
	return *ecfl == 0
}
func (ecfl *ExplicitContentFilterLvl) MembersWithoutRoles() bool {
	return *ecfl == 1
}
func (ecfl *ExplicitContentFilterLvl) AllMembers() bool {
	return *ecfl == 2
}

// MFA ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-mfa-level
type MFALvl uint

func (mfal *MFALvl) None() bool {
	return *mfal == 0
}
func (mfal *MFALvl) Elevated() bool {
	return *mfal == 1
}

// Verification ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-verification-level
type VerificationLvl uint

// None unrestricted
func (vl *VerificationLvl) None() bool {
	return *vl == 0
}

// Low must have verified email on account
func (vl *VerificationLvl) Low() bool {
	return *vl == 1
}

// Medium must be registered on Discord for longer than 5 minutes
func (vl *VerificationLvl) Medium() bool {
	return *vl == 2
}

// High (╯°□°）╯︵ ┻━┻ - must be a member of the server for longer than 10 minutes
func (vl *VerificationLvl) High() bool {
	return *vl == 3
}

// VeryHigh ┻━┻ミヽ(ಠ益ಠ)ﾉ彡┻━┻ - must have a verified phone number
func (vl *VerificationLvl) VeryHigh() bool {
	return *vl == 4
}

// DefaultMessageNotification ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type DefaultMessageNotificationLvl uint

func (dmnl *DefaultMessageNotificationLvl) AllMessages() bool {
	return *dmnl == 0
}
func (dmnl *DefaultMessageNotificationLvl) OnlyMentions() bool {
	return *dmnl == 1
}
func (dmnl *DefaultMessageNotificationLvl) Equals(v uint) bool {
	return uint(*dmnl) == v
}
