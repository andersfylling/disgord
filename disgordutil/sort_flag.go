package disgordutil

type sortAction uint32

type SortFieldType sortAction

// sort by field options
const (
	SortByID SortFieldType = iota
	SortByName
	SortByHoist
	SortByGuildID
	SortByChannelID
)

type SortOrderType sortAction

// elements order
const (
	// OrderAscending default when sorting
	OrderAscending  SortOrderType = iota
	OrderDescending
)