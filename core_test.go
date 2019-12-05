package kit

import (
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
	list := [][]string{
		[]string{"hi", "hello"},
		[]string{"hi", "hello world", "he"},
		[]string{"hi", "hello \"world\"", "he"},
		[]string{"hi", "hello \"world\"", "he"},
	}
	for i := range seed {
		if mis, res, ok := compareSplit(list[i], Split(seed[i])); !ok {
			t.Errorf("parse error: %v : %d %v", mis, len(res), res)
		}
	}
}
func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if mis, res, ok := compareSplit([]string{"hi", "hello \"world\"", "he"}, Split("hi 'hello \"world\"' he")); !ok {
			b.Errorf("parse error: %v : %d %v", mis, len(res), res)
		}
	}
}

func compareValue(mis interface{}, res interface{}) (interface{}, interface{}, bool) {
	switch mis := mis.(type) {
	case map[string]interface{}:
		if res, ok := res.(map[string]interface{}); !ok {
			return mis, res, false
		} else if len(mis) != len(res) {
			return mis, res, false
		} else {
			for k := range mis {
				if _, _, ok := compareValue(mis[k], res[k]); !ok {
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
				if _, _, ok := compareValue(mis[i], res[i]); !ok {
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
func TestValue(t *testing.T) {
	seed := []string{
		`{ "hi" "hello" }`,
		`{ "hi" "hello" "he" "world" }`,
		`[ "hi" "hello" ]`,
		`[ "hi" "hello" "he" "world" ]`,
	}
	list := []interface{}{
		map[string]interface{}{"hi": "hello"},
		map[string]interface{}{"hi": "hello", "he": "world"},
		[]interface{}{"hi", "hello"},
		[]interface{}{"hi", "hello", "he", "world"},
	}
	for i := range seed {
		if mis, res, ok := compareValue(list[i], Value(nil, "", Split(seed[i])...)); !ok {
			t.Errorf("parse error: %v : %v", mis, res)
		}
	}

}
func BenchmarkValue(b *testing.B) {
	seed := []string{
		`{ "hi" "hello" "he" "world" }`,
		`[ "hi" "hello" "he" "world" ]`,
	}
	list := []interface{}{
		map[string]interface{}{"hi": "hello", "he": "world"},
		[]interface{}{"hi", "hello", "he", "world"},
	}
	for i := 0; i < b.N; i++ {
		for i := range seed {
			if mis, res, ok := compareValue(list[i], Value(nil, "", Split(seed[i])...)); !ok {
				b.Errorf("parse error: %v : %v", mis, res)
			}
		}
	}
}
