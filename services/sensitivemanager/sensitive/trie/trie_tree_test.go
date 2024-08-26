package trie

import (
	"testing"
)

func TestTrieTree(t *testing.T) {
	tree := NewTrie()

	tree.Add("特朗普", "拜登", "50w")
	r := tree.Replace("拜登文化水平高于特朗普", '*')

	rExp := "**文化水平高于***"
	if r != rExp {
		t.Errorf("got %s, expect %s", r, rExp)
	}

	f := tree.Filter("支持特朗普的都是50w，这就是美利坚，听拜主席说")
	fExp := "支持的都是，这就是美利坚，听拜主席说"
	if f != fExp {
		t.Errorf("got %s, expect %s", f, fExp)
	}
}
