package pointer

// Person :
type Person struct {
	Name string `json:"name"`
	Age  string `json:"age"`

	Father *Person `json:"father" required:"false"`
	Mother *Person `json:"mother" required:"false"`
}
