package model

// Item :
type Item struct {
	Name string
}

// Items :
type Items []Item

// Item2 : newtype
type Item2 Item

// Item3 : alias
type Item3 = Item

// Item4 : duplicated
type Item4 struct {
	Name string
}

// Item5 : difference only required/unrequired
type Item5 struct {
	Name string `required:"false"`
}

// Items2 :
type Items2 Items

// Items3 :
type Items3 = Items

// ItemMap :
type ItemMap map[string]Item

// ItemMap2 :
type ItemMap2 ItemMap

// ItemMap3 :
type ItemMap3 = ItemMap
