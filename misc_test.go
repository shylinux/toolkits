package kit

import "testing"

func TestArgDef(t *testing.T) {
	cases := []struct {
		arg []string
		def []string
		exp []string
	}{
		{
			arg: []string{"begin_time", ""},
			def: []string{"begin_time", "2023"},
			exp: []string{"begin_time", "2023"},
		},
		{
			arg: []string{"begin_time", "2022"},
			def: []string{"begin_time", "2023"},
			exp: []string{"begin_time", "2022"},
		},
	}
	for _, c := range cases {
		res := ArgDef(c.arg, c.def...)
		if res[1] != c.exp[1] {
			t.Errorf("failure %v != %v", res, c.exp)
		}
	}
}
