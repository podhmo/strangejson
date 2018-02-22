package pointer

import (
	"encoding/json"
	"testing"
)

func TestLoad(t *testing.T) {
	{
		var p Person
		err := json.Unmarshal([]byte(`{"name": "foo"}`), &p)
		if err == nil {
			t.Fatal("must error")
		}
		t.Logf("expected: %q", err)
	}
	{
		var p Person
		err := json.Unmarshal([]byte(`{"age": 10}`), &p)
		if err == nil {
			t.Fatal("must error")
		}
		t.Logf("expected: %q", err)
	}
	{
		var p Person
		err := json.Unmarshal([]byte(`{"name": "foo", "age": 10}`), &p)
		if err != nil {
			t.Fatalf("unexpected: %q", err)
		}
	}
	{
		var p Person
		err := json.Unmarshal([]byte(`{"name": "foo", "age": 10, "father": {}}`), &p)
		if err == nil {
			t.Fatal("must error")
		}
		t.Logf("expected: %q", err)
	}
	{
		var p Person
		err := json.Unmarshal([]byte(`{"name": "foo", "age": 10, "father": {"name": "boo", "age": 30}}`), &p)
		if err != nil {
			t.Fatalf("unexpected: %q", err)
		}
	}
}
