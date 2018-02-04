package simple00

// User : user
type User struct {
	// Name : name of user
	Name     string `json:"name" required:"true"`
	Age      int    `json:"age"` // no required option, treated as required
	NickName string `json:"nickname" required:"false"`
}

// todo: interface
// todo: inline
// todo: json:"-"
// todo: external import
