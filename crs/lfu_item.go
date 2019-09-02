package crs

// newLFUItem ...
func newLFUItem(v interface{}) LFUItem {
	return LFUItem{
		Val: v,
	}
}

// LFUItem ...
type LFUItem struct {
	ID      Snowflake
	Val     interface{}
	counter uint64
}

func (i *LFUItem) increment() {
	i.counter++
}
