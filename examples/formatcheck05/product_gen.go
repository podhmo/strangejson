package model

import (
	"encoding/json"
	"fmt"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Item)
func (x *Item) FormatCheck() error {
	var merr *multierror.Error

	if err := x.Product.FormatCheck(); err != nil {
		merr = multierror.Append(merr, errors.WithMessage(err, "Product"))
	}
	return merr.ErrorOrNil()
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

	var merr *multierror.Error
	if p.Product == nil {
		merr = multierror.Append(merr, errors.New("product is required"))
	} else {
		x.Product = *p.Product
	}
	if p.Count == nil {
		merr = multierror.Append(merr, errors.New("count is required"))
	} else {
		x.Count = *p.Count
	}
	if merr != nil {
		return merr.ErrorOrNil()
	}
	return x.FormatCheck()
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Order)
func (x *Order) FormatCheck() error {
	var merr *multierror.Error

	for i, sub := range x.Items {
		if err := sub.FormatCheck(); err != nil {
			merr = multierror.Append(merr, errors.WithMessage(err, fmt.Sprintf("Items[%v]", i)))
		}
	}
	return merr.ErrorOrNil()
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

	var merr *multierror.Error
	if p.OrderedAt == nil {
		merr = multierror.Append(merr, errors.New("orderedAt is required"))
	} else {
		x.OrderedAt = *p.OrderedAt
	}
	if p.Items == nil {
		merr = multierror.Append(merr, errors.New("items is required"))
	} else {
		x.Items = *p.Items
	}
	if merr != nil {
		return merr.ErrorOrNil()
	}
	return x.FormatCheck()
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Product)
func (x *Product) FormatCheck() error {
	var merr *multierror.Error

	if len(x.Name) > 255 {
		merr = multierror.Append(merr, errors.New("name maxLength"))
	}
	if len(x.Name) < 1 {
		merr = multierror.Append(merr, errors.New("name minLength"))
	}
	return merr.ErrorOrNil()
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

	var merr *multierror.Error
	if p.Name == nil {
		merr = multierror.Append(merr, errors.New("name is required"))
	} else {
		x.Name = *p.Name
	}
	if p.Price == nil {
		merr = multierror.Append(merr, errors.New("price is required"))
	} else {
		x.Price = *p.Price
	}
	if merr != nil {
		return merr.ErrorOrNil()
	}
	return x.FormatCheck()
}

// FormatCheck : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Setting)
func (x *Setting) FormatCheck() error {
	var merr *multierror.Error

	for i, sub := range x.Products {
		if err := sub.FormatCheck(); err != nil {
			merr = multierror.Append(merr, errors.WithMessage(err, fmt.Sprintf("Products[%v]", i)))
		}
	}
	return merr.ErrorOrNil()
}

// UnmarshalJSON : (generated from github.com/podhmo/strangejson/examples/formatcheck05.Setting)
func (x *Setting) UnmarshalJSON(b []byte) error {
	type internal struct {
		Products *map[string]*Product `json:"products"`
	}

	var p internal
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	var merr *multierror.Error
	if p.Products == nil {
		merr = multierror.Append(merr, errors.New("products is required"))
	} else {
		x.Products = *p.Products
	}
	if merr != nil {
		return merr.ErrorOrNil()
	}
	return x.FormatCheck()
}
