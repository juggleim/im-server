package jpush

type Platform struct {
	IsAll bool
	Value []string
}

func NewPlatform() *Platform {
	return &Platform{}
}

func (p *Platform) Interface() interface{} {
	if p == nil || p.IsAll {
		return "all"
	}
	return p.Value
}

func (p *Platform) All() *Platform {
	p.IsAll = true
	return p
}

// 调用者负责去重
func (p *Platform) Add(oss ...string) *Platform {
	p.IsAll = false
	p.Value = append(p.Value, oss...)
	return p
}

func (p *Platform) Has(os string) bool {
	if p.IsAll {
		return true
	}
	for _, o := range p.Value {
		if o == os {
			return true
		}
	}
	return false
}
