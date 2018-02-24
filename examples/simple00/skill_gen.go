package simple00

import (
	"encoding/json"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/simple00.Skill)
func (x *Skill) FormatCheck() error {
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/simple00.Skill)
func (x *Skill) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name *string `json:"name"`
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
	return multierror.Append(merr, x.FormatCheck()).ErrorOrNil()
}
