package gee

import (
	"strings"
)

type trie struct {
	pattern  string  // 待匹配路由 /info/:id，只有注册的路由才有，其他节点只是用来查找的节点
	part     string  // 路由中的一部分 :id
	children []*trie // 子节点
	isWild   bool    // 是的精准匹配
}

func (t *trie) matchChild(part string) *trie {
	for _, child := range t.children {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

func (t *trie) matchChildren(part string) []*trie {
	tries := make([]*trie, 0)
	for _, child := range t.children {
		if child.part == part || child.isWild {
			tries = append(tries, child)
		}
	}
	return tries
}

func (t *trie) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		t.pattern = pattern
		return
	}

	part := parts[height]
	child := t.matchChild(part)
	if child == nil {
		child = &trie{part: part, isWild: part[0] == '*' || part[0] == ':'}
		t.children = append(t.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (t *trie) search(parts []string, height int) *trie {
	if len(parts) == height || strings.HasPrefix(t.part, "*") {
		if t.pattern == "" {
			return nil
		}
		return t
	}

	part := parts[height]
	children := t.matchChildren(part)

	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}

	return nil
}
