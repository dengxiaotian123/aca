package aca

import (
"unicode/utf8"
)

type node struct {
	next       map[rune]*node
	fail       *node
	wordLength int
}

type Trie struct {
	root      *node
	nodeCount int
}

// New returns an empty aca.
func New() *Trie {
	return &Trie{root: &node{}, nodeCount: 1}
}
func (T *Trie) InitRootNode() {
	n := new(node)
	n.fail = nil
	n.next = make(map[rune]*node)
	n.wordLength = 0

}

// Add adds a new word to aca.
// After Add, and before Find,
// MUST Build.
func (T *Trie) Add(word string) {
	n := T.root
	for _, r := range word {
		//	fmt.Println("r",r)
		if n.next == nil {
			n.next = make(map[rune]*node)
		}
		if n.next[r] == nil {
			n.next[r] = &node{}
			T.nodeCount++
		}
		n = n.next[r]
	}
	n.wordLength = len(word)
}

// Del delete a word from aca.
// After Del, and before Find,
// MUST Build.
func (T *Trie) Del(word string) {
	rs := []rune(word)
	//rs:=word
	stack := make([]*node, len(rs))
	n := T.root

	for i, r := range rs {
		if n.next[r] == nil {
			return
		}
		stack[i] = n
		n = n.next[r]
	}

	// if it is NOT the leaf node
	if len(n.next) > 0 {
		n.wordLength = 0
		return
	}

	// if it is the leaf node
	for i := len(rs) - 1; i >= 0; i-- {
		stack[i].next[rs[i]].next = nil
		stack[i].next[rs[i]].fail = nil

		delete(stack[i].next, rs[i])
		T.nodeCount--
		if len(stack[i].next) > 0 ||
			stack[i].wordLength > 0 {
			return
		}
	}
}

// BuildTrie for Agentids, It MUST be called before Find.
func (T *Trie) BuildTrie(dictionary [][]byte) {
	for _, line := range dictionary {
		if len(line) <= 0 {
			continue
		}
		T.Add(string(line))
	}
	T.Build()
}

//Build BuildTrie
func (T *Trie) Build() {
	// allocate enough memory as a queue
	q := append(make([]*node, 0, T.nodeCount), T.root)
	for len(q) > 0 {
		n := q[0]
		q = q[1:]

		for r, c := range n.next {
			q = append(q, c)

			p := n.fail
			for p != nil {
				// ATTENTION: nil map cannot be writen
				// but CAN BE READ!!!
				if p.next[r] != nil {
					c.fail = p.next[r]
					break
				}
				p = p.fail
			}
			if p == nil {
				c.fail = T.root
			}
		}
	}
}

func (T *Trie) find(s string, cb func(start, end int)) {
	n := T.root
	for i, r := range s {
		for n.next[r] == nil && n != T.root {
			n = n.fail
		}
		n = n.next[r]
		if n == nil {
			n = T.root
			continue
		}

		end := i + utf8.RuneLen(r)
		for t := n; t != T.root; t = t.fail {
			if t.wordLength > 0 {
				cb(end-t.wordLength, end)
			}
		}
	}
}

// Find finds all the words contains in s.
// The results may duplicated.
// It is caller's responsibility to make results unique.
func (T *Trie) Find(s string) (words []string) {

	T.find(s, func(start, end int) {
		//charactor, length := utf8.DecodeRune(s[start:end])

		words = append(words, s[start:end])
	})
	return
}

//AddSensitivWordFromTrie Add word
func (T *Trie) AddSensitivWordFromTrie(texts []string) {
	for _, text := range texts {
		T.Add(text)
	}

	T.Build()

}

//DelSensitivWordFromTrie Delete word
func (T *Trie) DelSensitivWordFromTrie(texts []string) {
	for _, text := range texts {
		T.Del(text)
	}
	T.Build()

}

