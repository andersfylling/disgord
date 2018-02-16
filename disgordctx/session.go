package disgordctx

type Context interface {
	isDisgordCtx()
}

type Session struct{}

func (s *Session) isDisgordCtx() {}
