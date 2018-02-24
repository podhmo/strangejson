package pointer

import (
	"encoding/json"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/pointer02.Person)
func (x *Person) FormatCheck() error {
	var merr *multierror.Error

	if x.Father != nil {
		if err := x.Father.FormatCheck(); err != nil {
			merr = multierror.Append(merr, errors.WithMessage(err, "Father"))
		}
	}
	if x.Mother != nil {
		if err := x.Mother.FormatCheck(); err != nil {
			merr = multierror.Append(merr, errors.WithMessage(err, "Mother"))
		}
	}
	return merr.ErrorOrNil()
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/pointer02.Person)
func (x *Person) UnmarshalJSON(b []byte) error {
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

	var merr *multierror.Error
	if p.Name == nil {
		merr = multierror.Append(merr, errors.New("name is required"))
	} else {
		x.Name = *p.Name
	}
	if p.Age == nil {
		merr = multierror.Append(merr, errors.New("age is required"))
	} else {
		x.Age = *p.Age
	}
	if p.Father != nil {
		x.Father = *p.Father
	}
	if p.Mother != nil {
		x.Mother = *p.Mother
	}
	return multierror.Append(merr, x.FormatCheck()).ErrorOrNil()
}
