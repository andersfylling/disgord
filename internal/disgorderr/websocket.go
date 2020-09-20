package disgorderr

type ClosedConnectionErr struct {
	info string
}

func (cce *ClosedConnectionErr) Error() string {
	return cce.info
}
