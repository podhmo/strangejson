package depends

import "time"

// BloodType :
type BloodType string

const (
	// BloodTypeA : AA or AO
	BloodTypeA = BloodType("A")
	// BloodTypeB : BB or BO
	BloodTypeB = BloodType("B")
	// BloodTypeAB : AB
	BloodTypeAB = BloodType("AB")
	// BloodTypeO : O
	BloodTypeO = BloodType("O")
)

// User : user
type User struct {
	// Name : name of user
	Name      string    `json:"name" required:"true"`
	Age       int       `json:"age"` // no required option, treated as required
	NickName  string    `json:"nickname" required:"false"`
	Birth     time.Time `json:"birth" required:"false"`
	BloodType BloodType `json:"bloodtype"`
	Skills    []Skill   `json:"skill"`
}
