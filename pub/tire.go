package pub

type TireNode struct {
	Ch     byte
	End    bool
	Times  int
	Childs [27]*TireNode
}

func (t *TireNode) Insert(s string, pos int) {
	index := byte(0)
	if s[pos] <= 'Z' && s[pos] >= 'A' {
		index = s[pos] - 'A' + 1
	} else if s[pos] <= 'z' && s[pos] >= 'a' {
		index = s[pos] - 'a' + 1
	}

	if t.Childs[index] == nil {
		t.Childs[index] = &TireNode{Ch: s[pos]}
	}
	t.Childs[index].Times++
	if pos == len(s)-1 {
		t.Childs[index].End = true
	} else {
		t.Childs[index].Insert(s, pos+1)
	}

}

func (t *TireNode) Find(s string, pos int) bool {
	index := byte(0)
	if s[pos] <= 'Z' && s[pos] >= 'A' {
		index = s[pos] - 'A' + 1
	} else if s[pos] <= 'z' && s[pos] >= 'a' {
		index = s[pos] - 'a' + 1
	}
	if t.Childs[index] == nil {
		return false
	}
	if pos == len(s)-1 {
		if t.Childs[index].End {
			return true
		}
		return false
	} else {
		return t.Childs[index].Find(s, pos+1)
	}
}

func (t *TireNode) FindMax(s []byte, pos, max int) int {
	if len(s) == 0 {
		return 0
	}

	index := byte(0)
	if s[pos] <= 'Z' && s[pos] >= 'A' {
		index = s[pos] - 'A' + 1
	} else if s[pos] <= 'z' && s[pos] >= 'a' {
		index = s[pos] - 'a' + 1
	}
	if t.Childs[index] == nil {
		return max
	}

	if t.Childs[index].End {
		max = pos + 1
	}

	if pos == len(s)-1 {
		return max
	} else {
		return t.Childs[index].FindMax(s, pos+1, max)
	}
}

func (t *TireNode) ReplaceWord(s []byte) {
	max := t.FindMax(s, 0, 0)
	for i := 0; i < max; i++ {
		s[i] = '*'
	}
}

func (t *TireNode) ReplaceSentence(s string) string {
	if s == "" {
		return s
	}

	last := 0
	res := []byte(s)

	for i := 0; i < len(s); i++ {
		if res[i] == ' ' || res[i] == ',' || res[i] == '.' {
			if i != 0 {
				t.ReplaceWord(res[last:i])
				last = i + 1
			}
		}
	}

	return string(res)
}

var RootTire *TireNode
