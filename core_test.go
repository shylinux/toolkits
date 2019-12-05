package kit

import (
	"encoding/json"
	"testing"
)

func compareSplit(mis []string, res []string) ([]string, []string, bool) {
	if len(res) != len(mis) {
		return mis, res, false
	}
	for i := range res {
		if res[i] != mis[i] {
			return mis, res, false
		}
	}
	return mis, res, true
}
func TestSplit(t *testing.T) {
	seed := []string{
		"hi hello",
		"hi 'hello world' he",
		"hi 'hello \"world\"' he",
		"   hi 'hello \"world\"' he",
	}
	miss := [][]string{
		[]string{"hi", "hello"},
		[]string{"hi", "hello world", "he"},
		[]string{"hi", "hello \"world\"", "he"},
		[]string{"hi", "hello \"world\"", "he"},
	}
	for i := range seed {
		if mis, res, ok := compareSplit(miss[i], Split(seed[i], "")); ok {
			t.Logf("parse ok: %v : %d %v", mis, len(res), res)
		} else {
			t.Errorf("parse error: %v : %d %v", mis, len(res), res)
		}
	}
}
func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if mis, res, ok := compareSplit([]string{"hi", "hello \"world\"", "he"}, Split("hi 'hello \"world\"' he")); !ok {
			b.Errorf("%v parse error: %v : %d %v", b.Name(), mis, len(res), res)
		}
	}
}

func compareParse(mis interface{}, res interface{}) (interface{}, interface{}, bool) {
	switch mis := mis.(type) {
	case map[string]interface{}:
		if res, ok := res.(map[string]interface{}); !ok {
			return mis, res, false
		} else if len(mis) != len(res) {
			return mis, res, false
		} else {
			for k := range mis {
				if _, _, ok := compareParse(mis[k], res[k]); !ok {
					return mis, res, false
				}
			}
		}
	case []interface{}:
		if res, ok := res.([]interface{}); !ok {
			return mis, res, false
		} else if len(mis) != len(res) {
			return mis, res, false
		} else {
			for i := range mis {
				if _, _, ok := compareParse(mis[i], res[i]); !ok {
					return mis, res, false
				}
			}
		}
	case string:
		if res, ok := res.(string); !ok {
			return mis, res, false
		} else if len(mis) != len(res) {
			return mis, res, false
		} else {
			return mis, res, mis == res
		}
	default:
	}
	return mis, res, true
}
func TestParse(t *testing.T) {
	seed := []string{
		// `{ "hi" "hello" }`,
		// `{ "hi" : "hello" , "nice" : { "hash" "hi" } "he" "world" }`,
		// `{ "hi" : "hello" , "nice" : [ "hash" "hi" ] "he" "world" }`,
		// `[ "hi" "hello" "he" "world" ]`,
		// `[ "hi" "hello" [ "ha" "haha" ] "he" "world" ]`,
		// `[ "hi" "hello" [ "ha" "haha" { 1 2 } ] "he" "world" ]`,
		// `[ "hi" "hello" { "ha" "haha" } "he" "world" ]`,
		`{ meta { text miss } list [
			{ meta { text you } list [ ] }
			{ meta { text hi } list [ ] }
		] }`,
	}
	miss := []interface{}{
		// map[string]interface{}{"hi": "hello"},
		// map[string]interface{}{"hi": "hello", "nice": map[string]interface{}{"hash": "hi"}, "he": "world"},
		// map[string]interface{}{"hi": "hello", "nice": []interface{}{"hash", "hi"}, "he": "world"},
		// []interface{}{"hi", "hello", "he", "world"},
		// []interface{}{"hi", "hello", []interface{}{"ha", "haha"}, "he", "world"},
		// []interface{}{"hi", "hello", []interface{}{"ha", "haha", map[string]interface{}{"1": "2"}}, "he", "world"},
		// []interface{}{"hi", "hello", map[string]interface{}{"ha": "haha"}, "he", "world"},
		map[string]interface{}{"meta": map[string]interface{}{"text": "miss"}, "list": []interface{}{
			map[string]interface{}{"meta": map[string]interface{}{"text": "you"}, "list": []interface{}{}},
			map[string]interface{}{"meta": map[string]interface{}{"text": "hi"}, "list": []interface{}{}},
		}},
	}
	for i := range seed {
		if mis, res, ok := compareParse(miss[i], Parse(nil, "", Split(seed[i])...)); ok {
			t.Logf("parse ok: %v : %v", mis, res)
		} else {
			t.Errorf("parse error: %v : %v", Formats(mis), Formats(res))
		}
	}

}
func BenchmarkParse(b *testing.B) {
	seed := []string{
		`[ "hi" "hello" [ "ha" "haha" { 1 2 } ] "he" "world" ]`,
	}
	miss := []interface{}{
		[]interface{}{"hi", "hello", []interface{}{"ha", "haha", map[string]interface{}{"1": "2"}}, "he", "world"},
	}
	for i := 0; i < b.N; i++ {
		for i := range seed {
			if mis, res, ok := compareParse(miss[i], Parse(nil, "", Split(seed[i])...)); !ok {
				b.Errorf("parse error: %v : %v", mis, res)
			}
		}
	}
}
func BenchmarkParses(b *testing.B) {
	seed := []string{
		`[ "hi" "hello" [ "ha" "haha" { 1 2 } ] "he" "world" ]`,
	}
	miss := []interface{}{
		[]interface{}{"hi", "hello", []interface{}{"ha", "haha", map[string]interface{}{"1": "2"}}, "he", "world"},
	}
	for i := 0; i < b.N; i++ {
		for i := range seed {
			var data interface{}
			json.Unmarshal([]byte(seed[i]), &data)
			if mis, res, ok := compareParse(miss[i], Parse(nil, "", Split(seed[i])...)); !ok {
				b.Errorf("parse error: %v : %v", mis, res)
			}
		}
	}
}

func TestValue(t *testing.T) {
	seed := [][]interface{}{
		[]interface{}{"hi", "hello"},
		[]interface{}{"hi.0", "hello"},
		[]interface{}{"hi.0.nice", "hello"},
	}
	miss := []interface{}{
		map[string]interface{}{"hi": "hello"},
		map[string]interface{}{"hi": []interface{}{"hello"}},
		map[string]interface{}{"hi": []interface{}{map[string]interface{}{"nice": "hello"}}},
	}
	for i := range seed {
		if mis, res, ok := compareParse(miss[i], Value(nil, seed[i][0], seed[i][1])); ok {
			t.Logf("parse ok: %v : %v", mis, res)
		} else {
			t.Errorf("parse error: %v : %v", mis, res)
		}
	}
}

func BenchmarkValue(b *testing.B) {
	seed := [][]interface{}{
		[]interface{}{"hi", "hello"},
		[]interface{}{"hi.0", "hello"},
		[]interface{}{"hi.0.nice", "hello"},
	}
	miss := []interface{}{
		map[string]interface{}{"hi": "hello"},
		map[string]interface{}{"hi": []interface{}{"hello"}},
		map[string]interface{}{"hi": []interface{}{map[string]interface{}{"nice": "hello"}}},
	}
	for i := 0; i < b.N; i++ {
		for i := range seed {
			if mis, res, ok := compareParse(miss[i], Value(nil, seed[i][0], seed[i][1])); !ok {
				b.Errorf("parse error: %v : %v", mis, res)
			}
		}
	}
}
