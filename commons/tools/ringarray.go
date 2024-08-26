package tools

type RingArray struct {
	data  []interface{}
	index int
	Count int
}

func NewRingArray(count int) *RingArray {
	if count <= 0 {
		return nil
	}
	return &RingArray{
		Count: count,
		index: -1,
		data:  make([]interface{}, count),
	}
}

func (arr *RingArray) Append(item interface{}) int {
	if item == nil {
		return -1
	}
	arr.index = (arr.index + 1) % arr.Count
	arr.data[arr.index] = item
	return arr.index
}

func (arr *RingArray) Foreach(f func(interface{}) bool) {
	if arr.index < 0 {
		return
	}
	start := (arr.index + 1) % arr.Count
	end := arr.index
	for {
		item := arr.data[start]
		if item != nil {
			isContinue := f(item)
			if !isContinue {
				break
			}
		}
		if start == end {
			break
		}
		start = (start + 1) % arr.Count
	}
}
