package kit

import (
	"testing"
	"time"
)

func TestInt(t *testing.T) {
	mis, _ := time.ParseInLocation("20060102", "20191029", time.Local)
	seed := []interface{}{
		10,
		"10",
		[]interface{}{1, 2},
		map[string]interface{}{"hi": "hello"},
		mis,
	}
	miss := []int{
		10,
		10,
		2,
		1,
		1572278400,
	}
	for i := range seed {
		if res := Int(seed[i]); res == miss[i] {
			t.Logf("parse ok: %v = %v : %v", miss[i], res, seed[i])
		} else {
			t.Errorf("parse error: %v = %v : %v", miss[i], res, seed[i])
		}
	}
}
func BenchmarkInt(b *testing.B) {
	mis, _ := time.ParseInLocation("20060102", "20191029", time.Local)
	seed := []interface{}{
		10,
		"10",
		[]interface{}{1, 2},
		map[string]interface{}{"hi": "hello"},
		mis,
	}
	miss := []int{
		10,
		10,
		2,
		1,
		1572278400,
	}
	for i := 0; i < b.N; i++ {
		for i := range seed {
			if res := Int(seed[i]); res != miss[i] {
				b.Errorf("parse error: %v = %v : %v", miss[i], res, seed[i])
			}
		}
	}
}
func TestFormat(t *testing.T) {
	seed := []interface{}{
		nil,
		10,
		"hello",
		map[string]interface{}{"hi": "hello", "0": false},
	}
	miss := []string{
		"",
		"10",
		"hello",
		`{"0":false,"hi":"hello"}`,
	}
	for i := range seed {
		if res := Format(seed[i]); res == miss[i] {
			t.Logf("parse ok: %v : %v", miss[i], res)
		} else {
			t.Errorf("parse error: %v : %v", miss[i], res)
		}
	}
}
func TestFormats(t *testing.T) {
	seed := []interface{}{
		nil,
		10,
		"hello",
		map[string]interface{}{"hi": "hello", "0": false},
	}
	miss := []string{
		"",
		"10",
		"hello",
		`{
  "0": false,
  "hi": "hello"
}`,
	}
	for i := range seed {
		if res := Formats(seed[i]); res == miss[i] {
			t.Logf("parse ok: %v : %v", miss[i], res)
		} else {
			t.Errorf("parse error: %v : %v", miss[i], res)
		}
	}
}
func BenchmarkFormat(b *testing.B) {
	seed := []interface{}{
		map[string]interface{}{"hi": "hello", "0": false},
	}
	miss := []string{
		`{"0":false,"hi":"hello"}`,
	}
	for i := 0; i < b.N; i++ {
		for i := range seed {
			if res := Format(seed[i]); res != miss[i] {
				b.Errorf("parse error: %v : %v", miss[i], res)
			}
		}
	}
}
func compareSimple(mis []string, res []string) ([]string, []string, bool) {
	if len(mis) != len(res) {
		return mis, res, false
	}
	for i := range mis {
		if mis[i] != res[i] {
			return mis, res, false
		}
	}
	return mis, res, true
}
func TestSimple(t *testing.T) {
	seed := []interface{}{
		nil,
		123.2,
		[]string{"123", "234"},
		[]interface{}{123, 1.2, true, "false"},
		map[string]interface{}{"hi": "hello", "0": false},
	}
	miss := [][]string{
		[]string{},
		[]string{"123"},
		[]string{"123", "234"},
		[]string{"123", "1", "true", "false"},
		[]string{"0", "false", "hi", "hello"},
	}
	for i := range seed {
		if mis, res, ok := compareSimple(miss[i], Simple(seed[i])); ok {
			t.Logf("parse ok: %v : %v", mis, res)
		} else {
			t.Errorf("parse error: %v : %v", mis, res)
		}
	}
}
