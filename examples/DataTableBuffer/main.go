package main

import (
	"fmt"
	"log"

	"github.com/danomagnum/gologix"
)

// Demo program for reading multiple tags at once using a DataTable buffer.
//
// A DataTable buffer (CIP class 0xB2) is created on the PLC, tags are associated
// with it, and then all tag values are read in a single request. This is the
// same mechanism RSLinx uses for grouped/trend reads.
//
// Three approaches are demonstrated:
//  1. AddTag      — add tags one at a time (separate 0x4E calls)
//  2. AddTags     — batch-add a map of tags in a single 0x4E call
//  3. AddTagGroup — schema-based with parameter substitution and multi-element expansion
func main() {

	// Setup the client.
	client := gologix.NewClient("192.168.2.241")

	err := client.Connect()
	if err != nil {
		log.Printf("Error opening client. %v", err)
		return
	}
	defer client.Disconnect()

	// =========================================================================
	// 1. AddTag — add tags one at a time
	// =========================================================================
	// Each call sends a separate CIP 0x4E message to the PLC.

	fmt.Println("=== AddTag (one at a time) ===")

	buf, err := client.NewDataTableBuffer()
	if err != nil {
		log.Fatalf("Error creating datatable buffer: %v", err)
	}

	var x int32
	err = buf.AddTagRef("TestDint", &x)
	if err != nil {
		log.Fatalf("Error adding DINT tag: %v", err)
	}

	err = buf.AddTag("TestReal", gologix.CIPTypeREAL)
	if err != nil {
		log.Fatalf("Error adding REAL tag: %v", err)
	}

	err = buf.AddTag("TestInt", gologix.CIPTypeINT)
	if err != nil {
		log.Fatalf("Error adding INT tag: %v", err)
	}

	// Read all tag values in a single request.
	// The PLC returns all values concatenated; ReadAll() splits them by
	// the tracked type sizes and decodes into Go types.
	values, err := buf.ReadAll()
	if err != nil {
		log.Fatalf("Error reading all: %v", err)
	}

	for name, val := range values {
		fmt.Printf("  %s = %v\n", name, val)
	}

	fmt.Printf("  Direct variable access: TestDint = %v\n", x)

	buf.Close()

	// =========================================================================
	// 2. AddTags — batch-add multiple tags in a single 0x4E call
	// =========================================================================
	// All tags are sent to the PLC in one message, which is more efficient
	// than calling AddTag repeatedly. Mixed type sizes are fully supported.

	fmt.Println("\n=== AddTags (batch map) ===")

	buf, err = client.NewDataTableBuffer()
	if err != nil {
		log.Fatalf("Error creating datatable buffer: %v", err)
	}

	err = buf.AddTags(map[string]gologix.CIPType{
		"TestDint":   gologix.CIPTypeDINT,
		"TestReal":   gologix.CIPTypeREAL,
		"TestInt":    gologix.CIPTypeINT,
		"TestString": gologix.CIPTypeSTRING,
	})
	if err != nil {
		log.Fatalf("Error batch-adding tags: %v", err)
	}

	values, err = buf.ReadAll()
	if err != nil {
		log.Fatalf("Error reading all: %v", err)
	}

	for name, val := range values {
		fmt.Printf("  %s = %v\n", name, val)
	}

	buf.Close()

	// =========================================================================
	// 3. AddTagGroup — schema with parameter substitution and array expansion
	// =========================================================================
	// A TagGroup defines a reusable schema of tag definitions. Tag names can
	// contain {0}, {1}, ... placeholders that are substituted when calling
	// AddTagGroup. Tags with Elements > 1 are automatically expanded into
	// individual buffer entries (e.g., Elements=3 creates 3 sequential tags).
	//
	// ReadAllTyped() returns a TagGroupResult with typed accessors and
	// automatically re-collapses multi-element tags back into slices.

	fmt.Println("\n=== AddTagGroup (schema with expansion) ===")

	buf, err = client.NewDataTableBuffer()
	if err != nil {
		log.Fatalf("Error creating datatable buffer: %v", err)
	}

	group := gologix.NewTagGroup(
		// Elements=3 expands to: TestDintArray[0], TestDintArray[1], TestDintArray[2]
		gologix.TagDef{Name: "TestDintArr[{0}]", Type: gologix.CIPTypeDINT, Elements: 3},
		// Elements=1 (or omitted) adds a single tag as-is.
		gologix.TagDef{Name: "TestReal", Type: gologix.CIPTypeREAL},
	)

	// Substitute {0} → 0 and add all expanded tags in a single 0x4E call.
	err = buf.AddTagGroup(group, 0)
	if err != nil {
		log.Fatalf("Error adding tag group: %v", err)
	}

	type MyStruct struct {
		TestDint int32 `gologix:"program:gologix_tests.{0}.Field1"`
		TestInt  int16 `gologix:"TestInt"`
	}

	var s MyStruct
	err = buf.AddTaggedStruct(&s, "readudt")

	// ReadAllTyped returns typed accessors and auto-collapses multi-element
	// tags back into slices.
	result, err := buf.ReadAllTyped()
	if err != nil {
		log.Fatalf("Error reading typed: %v", err)
	}

	if s.TestDint != 85456 {
		log.Fatalf("Unexpected value for TestDint: %v", s.TestDint)
	}

	fmt.Printf("struct value: %+v\n", s)

	// Multi-element tag → slice accessor
	dints, err := result.Int32Slice("TestDintArr[0]")
	if err != nil {
		log.Fatalf("Error getting int32 slice: %v", err)
	}
	fmt.Printf("  TestDintArr[0..2] = %v\n", dints)

	// Single tag → scalar accessor
	realVal, err := result.Float32("TestReal")
	if err != nil {
		log.Fatalf("Error getting float32: %v", err)
	}
	fmt.Printf("  TestReal = %v\n", realVal)

	buf.Close()
}
