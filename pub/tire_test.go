package pub

import "testing"

func TestTireNode_ReplaceSentence(t *testing.T) {
	root := &TireNode{}

	root.Insert("a", 0)
	root.Insert("ab", 0)
	root.Insert("abc", 0)

	root.Insert("de", 0)
	root.Insert("def", 0)

	root.Insert("fuck", 0)
	root.Insert("god", 0)
	root.Insert("what", 0)

	type table struct {
		str     string
		wantStr string
	}
	tables := []table{
		{"a.", "*."},
		{"ab.", "**."},
		{"abc.", "***."},
		{"abcd.", "***d."},

		{"bc.", "bc."},
		{"ac.", "*c."},

		{"dfe.", "dfe."},
		{"def.", "***."},

		{"abcdef.", "***def."},

		{"fuck,world.", "****,world."},
		{"fuck,god.", "****,***."},
		{"hello,god.", "hello,***."},

		{"fuckit,world.", "****it,world."},
		{"fuckit,god.", "****it,***."},

		{"hello,good,godoo.", "hello,good,***oo."},
	}
	for _, v := range tables {
		if res := root.ReplaceSentence(v.str); res != v.wantStr {
			t.Logf("failed test,str=%s,wantStr=%s,res=%s", v.str, v.wantStr, res)
		}
	}
}
