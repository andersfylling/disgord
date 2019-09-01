package crs

// newLFUItem ...
func newLFUItem(content interface{}) *LFUItem {
	return &LFUItem{
		Val: content,
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
