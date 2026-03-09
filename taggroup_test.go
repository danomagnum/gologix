package gologix

import (
	"reflect"
	"testing"
)

func TestFormatName(t *testing.T) {
	tests := []struct {
		template string
		args     []any
		expected string
	}{
		{"Rate[{0}]", []any{1}, "Rate[1]"},
		{"STATS[{0},{1}]", []any{3, 0}, "STATS[3,0]"},
		{"NoPlaceholder", []any{1}, "NoPlaceholder"},
		{"String[{0},0]", []any{42}, "String[42,0]"},
		{"{0}_{1}", []any{"Tag", "Name"}, "Tag_Name"},
		{"Simple", nil, "Simple"},
		{"Multi_{0}_{0}", []any{5}, "Multi_5_5"},
	}

	for _, tt := range tests {
		result := formatName(tt.template, tt.args...)
		if result != tt.expected {
			t.Errorf("formatName(%q, %v) = %q, want %q", tt.template, tt.args, result, tt.expected)
		}
	}
}

func TestNewTagGroupNormalize(t *testing.T) {
	g := NewTagGroup(
		TagDef{Name: "A", Type: CIPTypeDINT, Elements: 0},
		TagDef{Name: "B", Type: CIPTypeREAL, Elements: -1},
		TagDef{Name: "C", Type: CIPTypeINT, Elements: 5},
	)

	defs := g.Defs()
	if defs[0].Elements != 1 {
		t.Errorf("Elements=0 should normalize to 1, got %d", defs[0].Elements)
	}
	if defs[1].Elements != 1 {
		t.Errorf("Elements=-1 should normalize to 1, got %d", defs[1].Elements)
	}
	if defs[2].Elements != 5 {
		t.Errorf("Elements=5 should stay 5, got %d", defs[2].Elements)
	}
}

func TestNewTagGroupDefsIsCopy(t *testing.T) {
	g := NewTagGroup(TagDef{Name: "A", Type: CIPTypeDINT})
	defs := g.Defs()
	defs[0].Name = "Modified"
	if g.defs[0].Name == "Modified" {
		t.Error("Defs() should return a copy, not a reference")
	}
}

func TestZeroValueForTagDefScalars(t *testing.T) {
	tests := []struct {
		def      TagDef
		expected any
	}{
		{TagDef{Type: CIPTypeBOOL, Elements: 1}, false},
		{TagDef{Type: CIPTypeSINT, Elements: 1}, int8(0)},
		{TagDef{Type: CIPTypeBYTE, Elements: 1}, byte(0)},
		{TagDef{Type: CIPTypeINT, Elements: 1}, int16(0)},
		{TagDef{Type: CIPTypeUINT, Elements: 1}, uint16(0)},
		{TagDef{Type: CIPTypeDINT, Elements: 1}, int32(0)},
		{TagDef{Type: CIPTypeUDINT, Elements: 1}, uint32(0)},
		{TagDef{Type: CIPTypeLINT, Elements: 1}, int64(0)},
		{TagDef{Type: CIPTypeLWORD, Elements: 1}, uint64(0)},
		{TagDef{Type: CIPTypeREAL, Elements: 1}, float32(0)},
		{TagDef{Type: CIPTypeLREAL, Elements: 1}, float64(0)},
		{TagDef{Type: CIPTypeSTRING, Elements: 1}, ""},
	}

	for _, tt := range tests {
		result := zeroValueForTagDef(tt.def)
		if reflect.TypeOf(result) != reflect.TypeOf(tt.expected) {
			t.Errorf("zeroValueForTagDef(%v) type = %T, want %T", tt.def.Type, result, tt.expected)
		}
	}
}

func TestZeroValueForTagDefArrays(t *testing.T) {
	tests := []struct {
		def         TagDef
		expectedLen int
		expectedTyp string
	}{
		{TagDef{Type: CIPTypeDINT, Elements: 5}, 5, "[]int32"},
		{TagDef{Type: CIPTypeREAL, Elements: 3}, 3, "[]float32"},
		{TagDef{Type: CIPTypeINT, Elements: 10}, 10, "[]int16"},
		{TagDef{Type: CIPTypeSTRING, Elements: 2}, 2, "[]string"},
		{TagDef{Type: CIPTypeBOOL, Elements: 8}, 8, "[]bool"},
	}

	for _, tt := range tests {
		result := zeroValueForTagDef(tt.def)
		rv := reflect.ValueOf(result)
		if rv.Kind() != reflect.Slice {
			t.Errorf("zeroValueForTagDef(%v, Elements=%d) should be a slice, got %T", tt.def.Type, tt.def.Elements, result)
			continue
		}
		if rv.Len() != tt.expectedLen {
			t.Errorf("zeroValueForTagDef(%v, Elements=%d) len = %d, want %d", tt.def.Type, tt.def.Elements, rv.Len(), tt.expectedLen)
		}
		if rv.Type().String() != tt.expectedTyp {
			t.Errorf("zeroValueForTagDef(%v, Elements=%d) type = %s, want %s", tt.def.Type, tt.def.Elements, rv.Type().String(), tt.expectedTyp)
		}
	}
}

func TestExpandArrayTag(t *testing.T) {
	tests := []struct {
		baseName string
		count    int
		expected []string
	}{
		{
			"STATS[1,0]", 5,
			[]string{"STATS[1,0]", "STATS[1,1]", "STATS[1,2]", "STATS[1,3]", "STATS[1,4]"},
		},
		{
			"MyArray[0]", 3,
			[]string{"MyArray[0]", "MyArray[1]", "MyArray[2]"},
		},
		{
			"Tag[2,3,0]", 2,
			[]string{"Tag[2,3,0]", "Tag[2,3,1]"},
		},
		{
			"Scalar", 1,
			[]string{"Scalar"},
		},
	}

	for _, tt := range tests {
		result := expandArrayTag(tt.baseName, tt.count)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("expandArrayTag(%q, %d) = %v, want %v", tt.baseName, tt.count, result, tt.expected)
		}
	}
}

func TestExpandArrayTagScalar(t *testing.T) {
	// Scalar tag (no brackets) — returns the name once
	result := expandArrayTag("MyTag", 1)
	if len(result) != 1 || result[0] != "MyTag" {
		t.Errorf("expandArrayTag(scalar, 1) = %v, want [MyTag]", result)
	}
}

func TestTagGroupResultString(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"tag": "hello"}}

	s, err := r.String("tag")
	if err != nil || s != "hello" {
		t.Errorf("String() = %q, %v; want hello, nil", s, err)
	}

	_, err = r.String("missing")
	if err == nil {
		t.Error("String(missing) should return error")
	}

	r2 := &TagGroupResult{values: map[string]any{"tag": int32(5)}}
	_, err = r2.String("tag")
	if err == nil {
		t.Error("String(int32) should return type mismatch error")
	}
}

func TestTagGroupResultBool(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"tag": true}}

	b, err := r.Bool("tag")
	if err != nil || !b {
		t.Errorf("Bool() = %v, %v; want true, nil", b, err)
	}
}

func TestTagGroupResultInt16(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"tag": int16(999)}}

	i, err := r.Int16("tag")
	if err != nil || i != 999 {
		t.Errorf("Int16() = %d, %v; want 999, nil", i, err)
	}
}

func TestTagGroupResultInt32(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"tag": int32(42)}}

	i, err := r.Int32("tag")
	if err != nil || i != 42 {
		t.Errorf("Int32() = %d, %v; want 42, nil", i, err)
	}
}

func TestTagGroupResultUint32(t *testing.T) {
	// Test with uint32
	r := &TagGroupResult{values: map[string]any{"a": uint32(0x80000001)}}
	u, err := r.Uint32("a")
	if err != nil || u != 0x80000001 {
		t.Errorf("Uint32(uint32) = %d, %v; want %d, nil", u, err, uint32(0x80000001))
	}

	// Test with int32→uint32 coercion (common: gologix returns int32 for DINT)
	r2 := &TagGroupResult{values: map[string]any{"b": int32(-1)}}
	u2, err := r2.Uint32("b")
	if err != nil || u2 != 0xFFFFFFFF {
		t.Errorf("Uint32(int32(-1)) = %d, %v; want %d, nil", u2, err, uint32(0xFFFFFFFF))
	}
}

func TestTagGroupResultFloat32(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"tag": float32(3.14)}}

	f, err := r.Float32("tag")
	if err != nil || f != 3.14 {
		t.Errorf("Float32() = %f, %v; want 3.14, nil", f, err)
	}
}

func TestTagGroupResultFloat64(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"tag": float64(3.14159)}}

	f, err := r.Float64("tag")
	if err != nil || f != 3.14159 {
		t.Errorf("Float64() = %f, %v; want 3.14159, nil", f, err)
	}
}

func TestTagGroupResultInt32Slice(t *testing.T) {
	// Test with []int32 (from DataTable)
	r := &TagGroupResult{values: map[string]any{"a": []int32{1, 2, 3}}}
	s, err := r.Int32Slice("a")
	if err != nil || !reflect.DeepEqual(s, []int32{1, 2, 3}) {
		t.Errorf("Int32Slice([]int32) = %v, %v; want [1,2,3], nil", s, err)
	}

	// Test with []any{int32...} (from ReadMap)
	r2 := &TagGroupResult{values: map[string]any{"b": []any{int32(10), int32(20)}}}
	s2, err := r2.Int32Slice("b")
	if err != nil || !reflect.DeepEqual(s2, []int32{10, 20}) {
		t.Errorf("Int32Slice([]any) = %v, %v; want [10,20], nil", s2, err)
	}

	// Test type mismatch within []any
	r3 := &TagGroupResult{values: map[string]any{"c": []any{int32(1), "bad"}}}
	_, err = r3.Int32Slice("c")
	if err == nil {
		t.Error("Int32Slice with mixed types should return error")
	}
}

func TestTagGroupResultFloat32Slice(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"a": []float32{1.0, 2.0}}}
	s, err := r.Float32Slice("a")
	if err != nil || !reflect.DeepEqual(s, []float32{1.0, 2.0}) {
		t.Errorf("Float32Slice([]float32) = %v, %v; want [1.0,2.0], nil", s, err)
	}

	r2 := &TagGroupResult{values: map[string]any{"b": []any{float32(3.0), float32(4.0)}}}
	s2, err := r2.Float32Slice("b")
	if err != nil || !reflect.DeepEqual(s2, []float32{3.0, 4.0}) {
		t.Errorf("Float32Slice([]any) = %v, %v; want [3.0,4.0], nil", s2, err)
	}
}

func TestTagGroupResultValue(t *testing.T) {
	r := &TagGroupResult{values: map[string]any{"tag": int32(42)}}
	v, err := r.Value("tag")
	if err != nil || v != int32(42) {
		t.Errorf("Value() = %v, %v; want 42, nil", v, err)
	}

	_, err = r.Value("missing")
	if err == nil {
		t.Error("Value(missing) should return error")
	}
}

func TestTagGroupResultRaw(t *testing.T) {
	m := map[string]any{"a": int32(1), "b": "hello"}
	r := &TagGroupResult{values: m}
	raw := r.Raw()
	if raw["a"] != int32(1) || raw["b"] != "hello" {
		t.Errorf("Raw() returned unexpected map: %v", raw)
	}
}
