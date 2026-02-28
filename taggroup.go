package gologix

import (
	"fmt"
	"strings"
)

// TagDef describes a single PLC tag to read. Name may contain {0}, {1}, ...
// placeholders that are substituted at read time with runtime arguments.
type TagDef struct {
	Name     string  // tag path, e.g. "STATS[{0},0]"
	Type     CIPType // CIPTypeDINT, CIPTypeREAL, CIPTypeSTRING, etc.
	Elements int     // number of elements to read. 0 or 1 = scalar, >1 = array
}

// TagGroup is a reusable schema of tag definitions. Define once, read
// repeatedly with different parameter substitutions via ReadTagGroup
// or DataTableBuffer.AddTagGroup.
type TagGroup struct {
	defs []TagDef
}

// NewTagGroup creates a TagGroup from one or more TagDefs.
// Elements values of 0 are normalized to 1 (scalar).
func NewTagGroup(defs ...TagDef) *TagGroup {
	normalized := make([]TagDef, len(defs))
	copy(normalized, defs)
	for i := range normalized {
		if normalized[i].Elements <= 0 {
			normalized[i].Elements = 1
		}
	}
	return &TagGroup{defs: normalized}
}

// Defs returns a copy of the tag definitions.
func (g *TagGroup) Defs() []TagDef {
	out := make([]TagDef, len(g.defs))
	copy(out, g.defs)
	return out
}

// formatName replaces {0}, {1}, ... placeholders in a tag name template
// with the provided arguments.
func formatName(template string, args ...any) string {
	s := template
	for i, arg := range args {
		placeholder := fmt.Sprintf("{%d}", i)
		s = strings.ReplaceAll(s, placeholder, fmt.Sprintf("%v", arg))
	}
	return s
}

// zeroValueForTagDef returns a zero-valued Go variable matching the TagDef's
// CIP type and element count. The returned value is suitable for passing to
// ReadMap, which uses GoVarToCIPType() to infer the CIP type from the Go type.
func zeroValueForTagDef(def TagDef) any {
	if def.Elements > 1 {
		switch def.Type {
		case CIPTypeBOOL:
			return make([]bool, def.Elements)
		case CIPTypeSINT:
			return make([]int8, def.Elements)
		case CIPTypeBYTE, CIPTypeUSINT:
			return make([]byte, def.Elements)
		case CIPTypeINT:
			return make([]int16, def.Elements)
		case CIPTypeUINT:
			return make([]uint16, def.Elements)
		case CIPTypeDINT:
			return make([]int32, def.Elements)
		case CIPTypeUDINT, CIPTypeDWORD:
			return make([]uint32, def.Elements)
		case CIPTypeLINT:
			return make([]int64, def.Elements)
		case CIPTypeLWORD, CIPTypeULINT:
			return make([]uint64, def.Elements)
		case CIPTypeREAL:
			return make([]float32, def.Elements)
		case CIPTypeLREAL:
			return make([]float64, def.Elements)
		case CIPTypeSTRING:
			return make([]string, def.Elements)
		default:
			return make([]int32, def.Elements)
		}
	}
	switch def.Type {
	case CIPTypeBOOL:
		return false
	case CIPTypeSINT:
		return int8(0)
	case CIPTypeBYTE, CIPTypeUSINT:
		return byte(0)
	case CIPTypeINT:
		return int16(0)
	case CIPTypeUINT, CIPTypeWORD:
		return uint16(0)
	case CIPTypeDINT:
		return int32(0)
	case CIPTypeUDINT, CIPTypeDWORD:
		return uint32(0)
	case CIPTypeLINT:
		return int64(0)
	case CIPTypeLWORD, CIPTypeULINT:
		return uint64(0)
	case CIPTypeREAL:
		return float32(0)
	case CIPTypeLREAL:
		return float64(0)
	case CIPTypeSTRING:
		return ""
	default:
		return int32(0)
	}
}

// expandArrayTag takes a resolved tag name like "STATS[1,0]" and a count,
// and returns individual tag names with the last array dimension incremented:
// ["STATS[1,0]", "STATS[1,1]", ..., "STATS[1,4]"]
func expandArrayTag(baseName string, count int) []string {
	if count <= 1 {
		return []string{baseName}
	}

	tag, err := parse_tag_name(baseName)
	if err != nil || tag.Array_Order == nil || len(tag.Array_Order) == 0 {
		// No array index found — just return the base name repeated.
		// This shouldn't happen in normal use.
		result := make([]string, count)
		for i := range result {
			result[i] = baseName
		}
		return result
	}

	// Find the bracket position to reconstruct the name
	bracketPos := strings.LastIndex(baseName, "[")
	if bracketPos < 0 {
		result := make([]string, count)
		for i := range result {
			result[i] = baseName
		}
		return result
	}
	prefix := baseName[:bracketPos]

	lastDim := len(tag.Array_Order) - 1
	baseIndex := tag.Array_Order[lastDim]

	result := make([]string, count)
	for i := 0; i < count; i++ {
		indices := make([]string, len(tag.Array_Order))
		for d := 0; d < lastDim; d++ {
			indices[d] = fmt.Sprintf("%d", tag.Array_Order[d])
		}
		indices[lastDim] = fmt.Sprintf("%d", baseIndex+i)
		result[i] = fmt.Sprintf("%s[%s]", prefix, strings.Join(indices, ","))
	}
	return result
}

// ---------------------------------------------------------------------------
// TagGroupResult — typed accessors for read results
// ---------------------------------------------------------------------------

// TagGroupResult holds the results of a TagGroup read. Values are stored
// by resolved tag name (after parameter substitution).
type TagGroupResult struct {
	values map[string]any
}

// Raw returns the underlying map of tag name → value.
func (r *TagGroupResult) Raw() map[string]any {
	return r.values
}

// Value returns the raw any value for a tag name.
func (r *TagGroupResult) Value(name string) (any, error) {
	v, ok := r.values[name]
	if !ok {
		return nil, fmt.Errorf("tag %q not found in result", name)
	}
	return v, nil
}

// String returns the string value for the given resolved tag name.
func (r *TagGroupResult) String(name string) (string, error) {
	v, ok := r.values[name]
	if !ok {
		return "", fmt.Errorf("tag %q not found in result", name)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("tag %q is %T, not string", name, v)
	}
	return s, nil
}

// Bool returns the bool value for the given resolved tag name.
func (r *TagGroupResult) Bool(name string) (bool, error) {
	v, ok := r.values[name]
	if !ok {
		return false, fmt.Errorf("tag %q not found in result", name)
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("tag %q is %T, not bool", name, v)
	}
	return b, nil
}

// Int16 returns the int16 value for the given resolved tag name.
func (r *TagGroupResult) Int16(name string) (int16, error) {
	v, ok := r.values[name]
	if !ok {
		return 0, fmt.Errorf("tag %q not found in result", name)
	}
	i, ok := v.(int16)
	if !ok {
		return 0, fmt.Errorf("tag %q is %T, not int16", name, v)
	}
	return i, nil
}

// Int32 returns the int32 value for the given resolved tag name.
func (r *TagGroupResult) Int32(name string) (int32, error) {
	v, ok := r.values[name]
	if !ok {
		return 0, fmt.Errorf("tag %q not found in result", name)
	}
	i, ok := v.(int32)
	if !ok {
		return 0, fmt.Errorf("tag %q is %T, not int32", name, v)
	}
	return i, nil
}

// Uint32 returns the uint32 value for the given resolved tag name.
// Also handles int32→uint32 conversion since gologix returns int32 for
// DINT tags even when the PLC type is UDINT.
func (r *TagGroupResult) Uint32(name string) (uint32, error) {
	v, ok := r.values[name]
	if !ok {
		return 0, fmt.Errorf("tag %q not found in result", name)
	}
	switch x := v.(type) {
	case uint32:
		return x, nil
	case int32:
		return uint32(x), nil
	default:
		return 0, fmt.Errorf("tag %q is %T, not uint32 or int32", name, v)
	}
}

// Float32 returns the float32 value for the given resolved tag name.
func (r *TagGroupResult) Float32(name string) (float32, error) {
	v, ok := r.values[name]
	if !ok {
		return 0, fmt.Errorf("tag %q not found in result", name)
	}
	f, ok := v.(float32)
	if !ok {
		return 0, fmt.Errorf("tag %q is %T, not float32", name, v)
	}
	return f, nil
}

// Float64 returns the float64 value for the given resolved tag name.
func (r *TagGroupResult) Float64(name string) (float64, error) {
	v, ok := r.values[name]
	if !ok {
		return 0, fmt.Errorf("tag %q not found in result", name)
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("tag %q is %T, not float64", name, v)
	}
	return f, nil
}

// Int32Slice returns a []int32 for a multi-element tag read.
// Handles both []any (from ReadMap) and []int32 (from DataTable) return types.
func (r *TagGroupResult) Int32Slice(name string) ([]int32, error) {
	v, ok := r.values[name]
	if !ok {
		return nil, fmt.Errorf("tag %q not found in result", name)
	}
	switch x := v.(type) {
	case []int32:
		return x, nil
	case []any:
		result := make([]int32, len(x))
		for i, elem := range x {
			val, ok := elem.(int32)
			if !ok {
				return nil, fmt.Errorf("tag %q element %d is %T, not int32", name, i, elem)
			}
			result[i] = val
		}
		return result, nil
	default:
		return nil, fmt.Errorf("tag %q is %T, not []int32 or []any", name, v)
	}
}

// Float32Slice returns a []float32 for a multi-element tag read.
func (r *TagGroupResult) Float32Slice(name string) ([]float32, error) {
	v, ok := r.values[name]
	if !ok {
		return nil, fmt.Errorf("tag %q not found in result", name)
	}
	switch x := v.(type) {
	case []float32:
		return x, nil
	case []any:
		result := make([]float32, len(x))
		for i, elem := range x {
			val, ok := elem.(float32)
			if !ok {
				return nil, fmt.Errorf("tag %q element %d is %T, not float32", name, i, elem)
			}
			result[i] = val
		}
		return result, nil
	default:
		return nil, fmt.Errorf("tag %q is %T, not []float32 or []any", name, v)
	}
}
