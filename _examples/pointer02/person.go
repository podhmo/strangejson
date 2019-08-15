package pointer

// Person :
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`

	Father *Person `json:"father" required:"false"`
	Mother *Person `json:"mother" required:"false"`
}
