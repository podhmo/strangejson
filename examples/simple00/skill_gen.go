package simple00

import (
	"encoding/json"
	"errors"
)

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/simple00.Skill)
func (x *Skill) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name *string `json:"name"`
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.Name == nil {
		return errors.New("name is required")
	}
	x.Name = *p.Name
	return nil
}
