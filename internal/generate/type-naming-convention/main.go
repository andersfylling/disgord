package main

import (
	"fmt"
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
	"sort"
	"strings"
)

const (
	PKGName = "github.com/andersfylling/disgord"
)

var noStruct struct{}

func NewWhitelist(kinds ...types.Kind) Whitelist {
	whitelist := Whitelist{}
	for i := range kinds {
		whitelist[kinds[i]] = noStruct
	}
	return whitelist
}

type Whitelist map[types.Kind]struct{}

func (w Whitelist) ok(kind types.Kind) bool {
	if len(w) == 0 {
		return true
	}

	_, ok := w[kind]
	return ok
}

func DisgordTypes(whitelistedKinds ...types.Kind) (typesList []*types.Type, p *types.Package, err error) {
	return DisgordDefinitions(func(p *types.Package) map[string]*types.Type {
		return p.Types
	}, whitelistedKinds...)
}

func DisgordVars(whitelistedKinds ...types.Kind) (typesList []*types.Type, p *types.Package, err error) {
	return DisgordDefinitions(func(p *types.Package) map[string]*types.Type {
		return p.Variables
	}, whitelistedKinds...)
}

func DisgordConsts(whitelistedKinds ...types.Kind) (typesList []*types.Type, p *types.Package, err error) {
	return DisgordDefinitions(func(p *types.Package) map[string]*types.Type {
		return p.Constants
	}, whitelistedKinds...)
}

func DisgordDefinitions(target func(p *types.Package) map[string]*types.Type, whitelistedKinds ...types.Kind) (typesList []*types.Type, p *types.Package, err error) {
	builder := parser.New()
	if err := builder.AddDir(PKGName); err != nil {
		return nil, nil, fmt.Errorf("unable to add disgord package to gengo-parser builder. %w", err)
	}

	universe, err := builder.FindTypes()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to find types for disgord package. %w", err)
	}

	disgord := universe.Package(PKGName)
	whitelist := NewWhitelist(whitelistedKinds...)
	for _, typeData := range target(disgord) {
		if !whitelist.ok(typeData.Kind) {
			continue
		}

		typesList = append(typesList, typeData)
	}

	return typesList, disgord, nil
}



func Sort(tSlice []*types.Type) []*types.Type {
	sort.Slice(tSlice, func(i, j int) bool {
		name := func(t *types.Type) string {
			return strings.ToLower(t.Name.Name)
		}
		return name(tSlice[i]) < name(tSlice[j])
	})
	return tSlice
}

func FilterOutPrivateTypes(tSlice []*types.Type) []*types.Type {
	IsExported := func(name string) bool {
		firstChar := string(name[0])
		return firstChar == strings.ToUpper(firstChar)
	}

	filtered := make([]*types.Type, 0, len(tSlice))
	for i := range tSlice {
		if !IsExported(tSlice[i].Name.Name) {
			continue
		}
		filtered = append(filtered, tSlice[i])
	}

	return filtered
}

func main() {
	disgordTypes, _, err := DisgordTypes()
	if err != nil {
		panic(err)
	}

	disgordTypes = FilterOutPrivateTypes(disgordTypes)
	disgordTypes = Sort(disgordTypes)

	disgordVars, _, err := DisgordVars()
	if err != nil {
		panic(err)
	}

	disgordVars = FilterOutPrivateTypes(disgordVars)
	disgordVars = Sort(disgordVars)

	disgordConsts, _, err := DisgordConsts()
	if err != nil {
		panic(err)
	}

	disgordConsts = FilterOutPrivateTypes(disgordConsts)
	disgordConsts = Sort(disgordConsts)

	illegals := DirectionalNamesRule(disgordVars)
	illegals = append(illegals, DirectionalNamesRule(disgordVars)...)
	illegals = append(illegals, DirectionalNamesRule(disgordConsts)...)

	if len(illegals) > 0 {
		panic(fmt.Sprintf("%+v", illegals))
	}
}

func DirectionalNamesRule(typesList []*types.Type) (illegal []*types.Type) {
	HasCRUDIdentifier := func(name string) bool {
		for _, keyword := range []string{"Update", "Create", "Delete"} {
			if strings.HasPrefix(name, keyword) || strings.HasSuffix(name, keyword) {
				return true
			}
		}
		return false
	}

	for _, t := range typesList {
		if t.Kind == types.Struct {
			continue
		}

		if !HasCRUDIdentifier(TypeName(t)) {
			continue
		}

		// TODO: temporary edge case
		if t.Kind == types.Interface && strings.HasSuffix(TypeName(t), "Builder") {
			continue
		}

		illegal = append(illegal, t)
	}
	return illegal
}

func TypeName(t *types.Type) string {
	if strings.Contains(t.Name.Name, ".") {
		subs := strings.Split(t.Name.Name, ".")
		return subs[len(subs)-1]
	}
	return t.Name.Name
}