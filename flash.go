package dojo

type Flash struct {
	context Context
}

func NewFlash(ctx Context) *Flash {
	return &Flash{context: ctx}
}

func (f *Flash) Has(key string) bool {
	return f.context.Session().GetFlash(key) != nil
}

func (f *Flash) Get(key string) interface{} {
	m := f.context.Session().GetFlash(key)
	f.context.Session().Save()
	if m == nil {
		return nil
	}
	return m[0]
}
