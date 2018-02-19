package pointer

import (
	"encoding/json"
	"errors"
)

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/pointer02.Person)
func (x Person) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name   *string  `json:"name"`
		Age    *int     `json:"age"`
		Father **Person `json:"father" required:"false"`
		Mother **Person `json:"mother" required:"false"`
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
	if p.Father != nil {
		x.Father = *p.Father
	}
	if p.Mother != nil {
		x.Mother = *p.Mother
	}

	return nil
}
