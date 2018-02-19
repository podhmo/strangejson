package depends

import (
	"encoding/json"
	"errors"
)

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/depends01.User)
func (x User) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name      *string    `json:"name" required:"true"`
		Age       *int       `json:"age"`
		BloodType *BloodType `json:"bloodtype"`
		Skills    *[]Skill   `json:"skills"`
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
	if p.BloodType == nil {
		return errors.New("bloodtype is required")
	}
	x.BloodType = *p.BloodType
	if p.Skills == nil {
		return errors.New("skills is required")
	}
	x.Skills = *p.Skills

	return nil
}