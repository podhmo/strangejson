package simple00

import (
	"encoding/json"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/_examples/simple00.User)
func (x *User) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/_examples/simple00.User)
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
	if p.NickName != nil {
		x.NickName = *p.NickName
	}
	return multierror.Append(merr, x.FormatCheck()).ErrorOrNil()
}
