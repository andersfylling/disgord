package disgord

import (
	"context"
	"testing"
)

var fakePermissionsRole = &Role{ID: 10, Permissions: 2048}

type permissionTestingSession struct {
	getFakeRole bool
}

func (p *permissionTestingSession) GetGuildRoles(_ context.Context, _ Snowflake, _ ...Flag) ([]*Role, error) {
	if p.getFakeRole {
		return []*Role{fakePermissionsRole}, nil
	}
	return []*Role{}, nil
}

func TestChannel_GetPermissions_Overwrite(t *testing.T) {
	unmarshal := createUnmarshalUpdater(defaultUnmarshaler)

	data := []byte(`{"permission_overwrites": [{"allow": 2048, "deny": 0, "id": "1", "type": "member"}]}`)
	var c Channel
	if err := unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	p, err := c.GetPermissions(context.TODO(), &permissionTestingSession{}, &Member{UserID: 1, Roles: []Snowflake{}})
	if err != nil {
		t.Fatal(err)
	}
	if p != 2048 {
		t.Fatal("permissions should be 2048, is", p)
	}
}

func TestMember_GetPermissions(t *testing.T) {
	m := &Member{UserID: 1, Roles: []Snowflake{}}
	s := &permissionTestingSession{}
	p, err := m.GetPermissions(context.TODO(), s)
	if err != nil {
		t.Fatal(err)
	}
	if p != 0 {
		t.Fatal("permissions should be 0, is", p)
	}
	s.getFakeRole = true
	m.Roles = append(m.Roles, fakePermissionsRole.ID)
	p, err = m.GetPermissions(context.TODO(), s)
	if err != nil {
		t.Fatal(err)
	}
	if p != 2048 {
		t.Fatal("permissions should be 2048, is", p)
	}
}
