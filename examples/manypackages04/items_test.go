package manypackages

import (
	"encoding/json"
	"testing"

	"github.com/podhmo/strangejson/examples/manypackages04/item"
	"github.com/podhmo/strangejson/examples/manypackages04/item2"
	"github.com/podhmo/strangejson/examples/manypackages04/item3"
)

func TestItem(t *testing.T) {
	type C struct {
		input     string
		shouldErr bool
	}

	t.Run("one", func(t *testing.T) {
		candidates := []C{
			{
				input:     `{}`,
				shouldErr: true,
			},
			{
				input:     `{"name": "foo"}`,
				shouldErr: false,
			},
		}
		t.Run("Item", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x item.Item
					err := json.Unmarshal([]byte(c.input), &x)
					switch c.shouldErr {
					case true:
						if err == nil {
							t.Fatal("should error")
						}
						t.Logf("expected: %q", err)
					default:
						if err != nil {
							t.Fatalf("unexpected: %q", err)
						}
					}
				})
			}
		})

		t.Run("Item2", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x item2.Item2
					err := json.Unmarshal([]byte(c.input), &x)
					switch c.shouldErr {
					case true:
						if err == nil {
							t.Fatal("should error")
						}
						t.Logf("expected: %q", err)
					default:
						if err != nil {
							t.Fatalf("unexpected: %q", err)
						}
					}
				})
			}
		})

		t.Run("Item3", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x item3.Item3
					err := json.Unmarshal([]byte(c.input), &x)
					switch c.shouldErr {
					case true:
						if err == nil {
							t.Fatal("should error")
						}
						t.Logf("expected: %q", err)
					default:
						if err != nil {
							t.Fatalf("unexpected: %q", err)
						}
					}
				})
			}
		})
	})

	t.Run("many", func(t *testing.T) {
		candidates := []C{
			{
				input:     `[]`,
				shouldErr: false,
			},
			{
				input:     `{}`,
				shouldErr: true,
			},
			{
				input:     `[{}]`,
				shouldErr: true,
			},
			{
				input:     `[{"name": "foo"}]`,
				shouldErr: false,
			},
			{
				input:     `[{"name": "foo"}, {}]`,
				shouldErr: true,
			},
		}
		t.Run("Item", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x item.Items
					err := json.Unmarshal([]byte(c.input), &x)
					switch c.shouldErr {
					case true:
						if err == nil {
							t.Fatal("should error")
						}
						t.Logf("expected: %q", err)
					default:
						if err != nil {
							t.Fatalf("unexpected: %q", err)
						}
					}
				})
			}
		})

		t.Run("Item2", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x item2.Items2
					err := json.Unmarshal([]byte(c.input), &x)
					switch c.shouldErr {
					case true:
						if err == nil {
							t.Fatal("should error")
						}
						t.Logf("expected: %q", err)
					default:
						if err != nil {
							t.Fatalf("unexpected: %q", err)
						}
					}
				})
			}
		})

		t.Run("Item3", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x item3.Items3
					err := json.Unmarshal([]byte(c.input), &x)
					switch c.shouldErr {
					case true:
						if err == nil {
							t.Fatal("should error")
						}
						t.Logf("expected: %q", err)
					default:
						if err != nil {
							t.Fatalf("unexpected: %q", err)
						}
					}
				})
			}
		})
	})
}
