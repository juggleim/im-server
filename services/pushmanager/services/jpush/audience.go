package jpush

type Audience struct {
	IsAll bool
	Value map[string][]string
}

func NewAudience() *Audience {
	return &Audience{
		Value: make(map[string][]string),
	}
}

func (a *Audience) Interface() interface{} {
	if a.IsAll {
		return "all"
	}
	return a.Value
}

func (a *Audience) All() *Audience {
	a.IsAll = true
	return a
}

func (a *Audience) SetTag(tags ...string) *Audience {
	a.set("tag", tags)
	return a
}

func (a *Audience) SetTagAnd(tagAnds ...string) *Audience {
	a.set("tag_and", tagAnds)
	return a
}

func (a *Audience) SetTagNot(tagNots ...string) *Audience {
	a.set("tag_not", tagNots)
	return a
}

func (a *Audience) SetRegistrationId(regIds ...string) *Audience {
	a.set("registration_id", regIds)
	return a
}

func (a *Audience) SetSegment(segments ...string) *Audience {
	a.set("segment", segments)
	return a
}

func (a *Audience) SetAbtest(abtests ...string) *Audience {
	a.set("abtest", abtests)
	return a
}

func (a *Audience) SetAlias(alias ...string) *Audience {
	a.set("alias", alias)
	return a
}

func (a *Audience) set(key string, v []string) {
	if len(v) > 0 {
		a.IsAll = false
		a.Value[key] = v
	}
}
