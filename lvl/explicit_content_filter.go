package lvl

// ExplicitContentFilter ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
type ExplicitContentFilter uint

func (ecfl *ExplicitContentFilter) Disabled() bool {
	return *ecfl == 0
}
func (ecfl *ExplicitContentFilter) MembersWithoutRoles() bool {
	return *ecfl == 1
}
func (ecfl *ExplicitContentFilter) AllMembers() bool {
	return *ecfl == 2
}
