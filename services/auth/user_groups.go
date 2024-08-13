package auth

import "encoding/json"

type UserGroup struct {
	ID
	Name     string
	Members  []ID
	Profiles map[string]Permissions
}

type UserGroups []UserGroup

func (g UserGroups) Match(id ID) bool {
	for _, group := range g {
		if group.ID == id {
			return true
		}
	}

	return false
}

func (g UserGroups) MarshalIDs() ([]byte, error) {
	var ids []ID
	for _, group := range g {
		ids = append(ids, group.ID)
	}

	data, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshGroupIDs(data []byte) ([]ID, error) {
	var ids []ID
	err := json.Unmarshal(data, &ids)

	return ids, err
}

func (g UserGroup) UnmarshalIDs(data []byte) error {
	ids, err := UnmarshGroupIDs(data)
	if err != nil {
		return err
	}

	g.Members = ids
	return nil
}
