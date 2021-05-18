package dojo

const FlashOldKey = "_old_inputs"

type Old struct {
	context Context
}

func NewOld(ctx Context) *Old {
	return &Old{context: ctx}
}

func (o *Old) getData() map[string]interface{} {
	m := o.context.Session().GetFlash(FlashOldKey)
	return m[0].(map[string]interface{})
}

func (o *Old) Has(key string) bool {
	return o.getData()[key] != nil
}

func (o *Old) Get(key string) interface{} {
	v := o.getData()[key]
	_ = o.context.Session().Save()
	return v
}
