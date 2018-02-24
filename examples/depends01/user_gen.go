package depends

import (
	"encoding/json"
	"fmt"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/depends01.User)
func (x *User) FormatCheck() error {
	var merr *multierror.Error

	for i, sub := range x.Skills {
		if err := sub.FormatCheck(); err != nil {
			merr = multierror.Append(merr, errors.WithMessage(err, fmt.Sprintf("Skills[%d]", i)))
		}
	}
	return merr.ErrorOrNil()
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/depends01.User)
func (x *User) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name      *string    `json:"name" required:"true"`
		Age       *int       `json:"age"`
		NickName  *string    `json:"nickname" required:"false"`
		Birth     *time.Time `json:"birth" required:"false"`
		BloodType *BloodType `json:"bloodtype"`
		Skills    *[]Skill   `json:"skills"`
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
	if p.Birth != nil {
		x.Birth = *p.Birth
	}
	if p.BloodType == nil {
		merr = multierror.Append(merr, errors.New("bloodtype is required"))
	} else {
		x.BloodType = *p.BloodType
	}
	if p.Skills == nil {
		merr = multierror.Append(merr, errors.New("skills is required"))
	} else {
		x.Skills = *p.Skills
	}
	return multierror.Append(merr, x.FormatCheck()).ErrorOrNil()
}
