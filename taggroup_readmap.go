package gologix

import "fmt"

// ReadTagGroup reads all tags defined in the TagGroup in a single network
// call (or the minimum number of calls if tags exceed the connection size
// limit). The args are substituted into {0}, {1}, ... placeholders in tag names.
//
// Internally this builds a ReadMap call using the existing CIP multi-service
// (0x0A) batching and auto-splitting logic.
//
// Example:
//
//	group := gologix.NewTagGroup(
//	    gologix.TagDef{Name: "RATE[{0}]", Type: gologix.CIPTypeREAL},
//	    gologix.TagDef{Name: "STATS[{0},0]", Type: gologix.CIPTypeDINT, Elements: 5},
//	)
//	result, err := client.ReadTagGroup(group, 1)
//	rate, _ := result.Float32("RATE[1]")
//	stats, _ := result.Int32Slice("STATS[1,0]")
func (client *Client) ReadTagGroup(group *TagGroup, args ...any) (*TagGroupResult, error) {
	err := client.checkConnection()
	if err != nil {
		return nil, fmt.Errorf("could not start tag group read: %w", err)
	}

	m := make(map[string]any, len(group.defs))

	for _, def := range group.defs {
		name := formatName(def.Name, args...)
		m[name] = zeroValueForTagDef(def)
	}

	err = client.ReadMap(m)
	if err != nil {
		return nil, fmt.Errorf("tag group read failed: %w", err)
	}

	return &TagGroupResult{values: m}, nil
}
