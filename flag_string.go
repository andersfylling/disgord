// Code generated by "stringer -type=Flag"; DO NOT EDIT.

package disgord

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[IgnoreCache-1]
	_ = x[IgnoreEmptyParams-2]
	_ = x[SortByID-4]
	_ = x[SortByName-8]
	_ = x[SortByHoist-16]
	_ = x[SortByGuildID-32]
	_ = x[SortByChannelID-64]
	_ = x[OrderAscending-128]
	_ = x[OrderDescending-256]
}

const (
	_Flag_name_0 = "IgnoreCacheIgnoreEmptyParams"
	_Flag_name_1 = "SortByID"
	_Flag_name_2 = "SortByName"
	_Flag_name_3 = "SortByHoist"
	_Flag_name_4 = "SortByGuildID"
	_Flag_name_5 = "SortByChannelID"
	_Flag_name_6 = "OrderAscending"
	_Flag_name_7 = "OrderDescending"
)

var (
	_Flag_index_0 = [...]uint8{0, 11, 28}
)

// String
//
// Deprecated: schedule for removal
func (i Flag) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _Flag_name_0[_Flag_index_0[i]:_Flag_index_0[i+1]]
	case i == 4:
		return _Flag_name_1
	case i == 8:
		return _Flag_name_2
	case i == 16:
		return _Flag_name_3
	case i == 32:
		return _Flag_name_4
	case i == 64:
		return _Flag_name_5
	case i == 128:
		return _Flag_name_6
	case i == 256:
		return _Flag_name_7
	default:
		return "Flag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
