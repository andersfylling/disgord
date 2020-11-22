// Code generated by generate/interfaces; DO NOT EDIT.

package disgord

func (i *InviteMetadata) deepCopy() interface{} {
	cp := &InviteMetadata{}
	_ = DeepCopyOver(cp, i)
	return cp
}

func (e *Emoji) deepCopy() interface{} {
	cp := &Emoji{}
	_ = DeepCopyOver(cp, e)
	return cp
}

func (c *Channel) deepCopy() interface{} {
	cp := &Channel{}
	_ = DeepCopyOver(cp, c)
	return cp
}

func (u *User) deepCopy() interface{} {
	cp := &User{}
	_ = DeepCopyOver(cp, u)
	return cp
}

func (i *Invite) deepCopy() interface{} {
	cp := &Invite{}
	_ = DeepCopyOver(cp, i)
	return cp
}

func (m *Message) deepCopy() interface{} {
	cp := &Message{}
	_ = DeepCopyOver(cp, m)
	return cp
}

func (v *VoiceState) deepCopy() interface{} {
	cp := &VoiceState{}
	_ = DeepCopyOver(cp, v)
	return cp
}

func (g *Guild) deepCopy() interface{} {
	cp := &Guild{}
	_ = DeepCopyOver(cp, g)
	return cp
}

func (v *VoiceRegion) deepCopy() interface{} {
	cp := &VoiceRegion{}
	_ = DeepCopyOver(cp, v)
	return cp
}

func (r *Role) deepCopy() interface{} {
	cp := &Role{}
	_ = DeepCopyOver(cp, r)
	return cp
}
