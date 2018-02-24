package model

import (
	"encoding/json"
	"errors"
	"time"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Item)
func (x *Item) FormatCheck() error {
	if err := x.Product.FormatCheck(); err != nil {
		return err
	}
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Item)
func (x *Item) UnmarshalJSON(b []byte) error {
	type internal struct {
		Product *Product `json:"product"`
		Count   *int     `json:"count"`
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.Product == nil {
		return errors.New("product is required")
	}
	x.Product = *p.Product
	if p.Count == nil {
		return errors.New("count is required")
	}
	x.Count = *p.Count
	return x.FormatCheck()
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Order)
func (x *Order) FormatCheck() error {
	for _, sub := range x.Items {
		if err := sub.FormatCheck(); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Order)
func (x *Order) UnmarshalJSON(b []byte) error {
	type internal struct {
		OrderedAt *time.Time `json:"orderedAt"`
		Items     *[]Item    `json:"items"`
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.OrderedAt == nil {
		return errors.New("orderedAt is required")
	}
	x.OrderedAt = *p.OrderedAt
	if p.Items == nil {
		return errors.New("items is required")
	}
	x.Items = *p.Items
	return x.FormatCheck()
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Product)
func (x *Product) FormatCheck() error {
	if len(x.Name) > 255 {
		return errors.New("max")
	}
	if len(x.Name) < 1 {
		return errors.New("min")
	}
	return nil
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Product)
func (x *Product) UnmarshalJSON(b []byte) error {
	type internal struct {
		Name  *string `json:"name" required:"true" minLength:"1" maxLength:"255"`
		Price *int    `json:"price"`
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p.Name == nil {
		return errors.New("name is required")
	}
	x.Name = *p.Name
	if p.Price == nil {
		return errors.New("price is required")
	}
	x.Price = *p.Price
	return x.FormatCheck()
}
