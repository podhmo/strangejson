package codegen

import "go/types"

type cache struct {
	unmarshalStructCache map[*types.Struct]types.Object
}
