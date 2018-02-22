package model

import (
	"encoding/json"
	"testing"
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
					var x Item
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
					var x Item2
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
					var x Item3
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
		t.Run("Item4", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x Item4
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
		t.Run("Item5", func(t *testing.T) {
			for _, c := range candidates {
				c := c
				t.Run(c.input, func(t *testing.T) {
					var x Item5
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
					var x Items
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
					var x Items2
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
					var x Items3
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
