package simple00

import (
	"encoding/json"
	"errors"
)

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/simple00.User)
func (x *User) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name     *string `json:"name" required:"true"`
		Age      *int    `json:"age"`
		NickName *string `json:"nickname" required:"false"`
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.Name == nil {
		return errors.New("name is required")
	}
	x.Name = *p.Name
	if p.Age == nil {
		return errors.New("age is required")
	}
	x.Age = *p.Age
	if p.NickName != nil {
		x.NickName = *p.NickName
	}
	return nil
}
