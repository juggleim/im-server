package sensitive

import (
	"bufio"
	"fmt"
	"im-server/services/sensitivemanager/sensitive/trie"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Filter 敏感词过滤器
type Filter struct {
	trie  *trie.Trie
	noise *regexp.Regexp
}

// New 返回一个敏感词过滤器
func New() *Filter {
	return &Filter{
		trie:  trie.NewTrie(),
		noise: regexp.MustCompile(`[|\s&%$@*！!#^~_—｜'";.。，,?<>《》：:]+`),
	}
}

// AddWord 添加敏感词
func (filter *Filter) AddWord(words ...string) {
	filter.trie.Add(words...)
}

// DelWord 删除敏感词
func (filter *Filter) DelWord(words ...string) {
	filter.trie.Del(words...)
}

// Filter 过滤敏感词
func (filter *Filter) Filter(text string) string {
	return filter.trie.Filter(text)
}

// Replace 和谐敏感词
func (filter *Filter) Replace(text string, repl rune) string {
	return filter.trie.Replace(text, repl)
}

// FindIn 检测敏感词
func (filter *Filter) FindIn(text string) (bool, string) {
	text = filter.RemoveNoise(text)
	return filter.trie.FindIn(text)
}

// FindAll 找到所有匹配词
func (filter *Filter) FindAll(text string) []string {
	return filter.trie.FindAll(text)
}

// Validate 检测字符串是否合法
func (filter *Filter) Validate(text string) (bool, string) {
	text = filter.RemoveNoise(text)
	return filter.trie.Validate(text)
}

// RemoveNoise 去除空格等噪音，噪音可以使用 UpdateNoisePattern 更新
func (filter *Filter) RemoveNoise(text string) string {
	return filter.noise.ReplaceAllString(text, "")
}

// UpdateNoisePattern 更新去噪模式
func (filter *Filter) UpdateNoisePattern(pattern string) {
	filter.noise = regexp.MustCompile(pattern)
}

// LoadWordDict 加载本地敏感词字典
func (filter *Filter) LoadWordDict(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(f)

	return filter.Load(f)
}

// LoadNetWordDict 加载网络敏感词字典
func (filter *Filter) LoadNetWordDict(url string) error {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	rsp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(rsp.Body)

	return filter.Load(rsp.Body)
}

// Load common method to add words
func (filter *Filter) Load(rd io.Reader) error {
	buf := bufio.NewReader(rd)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		filter.trie.Add(string(line))
	}

	return nil
}
