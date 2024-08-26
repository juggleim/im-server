package trie

// Trie 短语组成的Trie树.
type Trie struct {
	Root *Node
}

// Node Trie树上的一个节点.
type Node struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node
}

// NewTrie 新建一棵Trie
func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(0),
	}
}

// Add 添加若干个词
func (tree *Trie) Add(words ...string) {
	for _, word := range words {
		tree.add(word)
	}
}

func (tree *Trie) add(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode(r)
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.isPathEnd = true
		}
	}
}

func (tree *Trie) Del(words ...string) {
	for _, word := range words {
		tree.del(word)
	}
}

func (tree *Trie) del(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; !ok {
			return
		} else {
			current = next
		}

		if position == len(runes)-1 {
			current.SoftDel()
		}
	}
}

// Replace 词语替换
func (tree *Trie) Replace(text string, character rune) string {
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		// println(string(current.Character), current.IsPathEnd(), left)
		if current.IsPathEnd() && left <= position {
			for i := left; i <= position; i++ {
				runes[i] = character
			}
		}

		parent = current
	}

	return string(runes)
}

// Filter 直接过滤掉字符串中的敏感词
func (tree *Trie) Filter(text string) string {
	var (
		parent      = tree.Root
		current     *Node
		left        = 0
		found       bool
		runes       = []rune(text)
		length      = len(runes)
		resultRunes = make([]rune, 0, length)
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			resultRunes = append(resultRunes, runes[left])
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() {
			left = position + 1
			parent = tree.Root
		} else {
			parent = current
		}

	}

	resultRunes = append(resultRunes, runes[left:]...)
	return string(resultRunes)
}

// Validate 验证字符串是否合法，如不合法则返回false和检测到
// 的第一个敏感词
func (tree *Trie) Validate(text string) (bool, string) {
	const (
		Empty = ""
	)
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			return false, string(runes[left : position+1])
		}

		parent = current
	}

	return true, Empty
}

// FindIn 判断text中是否含有词库中的词
func (tree *Trie) FindIn(text string) (bool, string) {
	validated, first := tree.Validate(text)
	return !validated, first
}

// FindAll 找有所有包含在词库中的词
func (tree *Trie) FindAll(text string) []string {
	var matches []string
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			matches = append(matches, string(runes[left:position+1]))
		}

		if position == length-1 {
			parent = tree.Root
			position = left
			left++
			continue
		}

		parent = current
	}

	var i = 0
	if count := len(matches); count > 0 {
		set := make(map[string]struct{})
		for i < count {
			_, ok := set[matches[i]]
			if !ok {
				set[matches[i]] = struct{}{}
				i++
				continue
			}
			count--
			copy(matches[i:], matches[i+1:])
		}
		return matches[:count]
	}

	return nil
}

// NewNode 新建子节点
func NewNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node),
	}
}

// NewRootNode 新建根节点
func NewRootNode(character rune) *Node {
	return &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node),
	}
}

// IsLeafNode 判断是否叶子节点
func (node *Node) IsLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断是否为根节点
func (node *Node) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断是否为某个路径的结束
func (node *Node) IsPathEnd() bool {
	return node.isPathEnd
}

// SoftDel 置软删除状态
func (node *Node) SoftDel() {
	node.isPathEnd = false
}
