package gologix_tests

import (
	"testing"

	"github.com/danomagnum/gologix"
)

func TestDataTableBuffer(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
	if err := client.Connect(); err != nil {
		t.Fatalf("Error opening client: %v", err)
	}
	defer client.Disconnect()

	// 1. AddTagRef, AddTag
	buf, err := client.NewDataTableBuffer()
	if err != nil {
		t.Fatalf("Error creating datatable buffer: %v", err)
	}

	var x int32
	if err := buf.AddTagRef("TestDint", &x); err != nil {
		t.Fatalf("Error adding DINT tag: %v", err)
	}
	if err := buf.AddTag("TestReal", gologix.CIPTypeREAL); err != nil {
		t.Fatalf("Error adding REAL tag: %v", err)
	}
	if err := buf.AddTag("TestInt", gologix.CIPTypeINT); err != nil {
		t.Fatalf("Error adding INT tag: %v", err)
	}

	values, err := buf.ReadAll()
	if err != nil {
		t.Fatalf("Error reading all: %v", err)
	}
	if v, ok := values["TestDint"]; !ok || v != int32(36) {
		t.Errorf("TestDint = %v, want 36", v)
	}
	if v, ok := values["TestReal"]; !ok || v != float32(93.45) {
		t.Errorf("TestReal = %v, want 93.45", v)
	}
	if v, ok := values["TestInt"]; !ok || v != int16(999) {
		t.Errorf("TestInt = %v, want 999", v)
	}
	if x != 36 {
		t.Errorf("Direct variable access: TestDint = %v, want 36", x)
	}
	buf.Close()

	// 2. AddTags
	buf, err = client.NewDataTableBuffer()
	if err != nil {
		t.Fatalf("Error creating datatable buffer: %v", err)
	}
	if err := buf.AddTags(map[string]gologix.CIPType{
		"TestDint":   gologix.CIPTypeDINT,
		"TestReal":   gologix.CIPTypeREAL,
		"TestInt":    gologix.CIPTypeINT,
		"TestString": gologix.CIPTypeSTRING,
	}); err != nil {
		t.Fatalf("Error batch-adding tags: %v", err)
	}
	values, err = buf.ReadAll()
	if err != nil {
		t.Fatalf("Error reading all: %v", err)
	}
	if v, ok := values["TestString"]; !ok || v != "Something" {
		t.Errorf("TestString = %v, want Something", v)
	}
	if v, ok := values["TestDint"]; !ok || v != int32(36) {
		t.Errorf("TestDint = %v, want 36", v)
	}
	if v, ok := values["TestReal"]; !ok || v != float32(93.45) {
		t.Errorf("TestReal = %v, want 93.45", v)
	}
	if v, ok := values["TestInt"]; !ok || v != int16(999) {
		t.Errorf("TestInt = %v, want 999", v)
	}
	buf.Close()

	// 3. AddTagGroup and AddTagStruct
	buf, err = client.NewDataTableBuffer()
	if err != nil {
		t.Fatalf("Error creating datatable buffer: %v", err)
	}
	group := gologix.NewTagGroup(
		gologix.TagDef{Name: "TestDintArr[{0}]", Type: gologix.CIPTypeDINT, Elements: 3},
		gologix.TagDef{Name: "TestReal", Type: gologix.CIPTypeREAL},
	)
	if err := buf.AddTagGroup(group, 0); err != nil {
		t.Fatalf("Error adding tag group: %v", err)
	}

	type MyStruct struct {
		TestDint int32 `gologix:"program:gologix_tests.{0}.Field1"`
		TestInt  int16 `gologix:"TestInt"`
	}
	var s MyStruct
	if err := buf.AddTaggedStruct(&s, "readudt"); err != nil {
		t.Fatalf("Error adding tag struct: %v", err)
	}

	result, err := buf.ReadAllTyped()
	if err != nil {
		t.Fatalf("Error reading typed: %v", err)
	}
	if s.TestDint != 85456 {
		t.Errorf("struct value: TestDint = %v, want 85456", s.TestDint)
	}
	if s.TestInt != 999 {
		t.Errorf("struct value: TestInt = %v, want 999", s.TestInt)
	}
	dints, err := result.Int32Slice("TestDintArr[0]")
	if err != nil {
		t.Errorf("Error getting int32 slice: %v", err)
	} else {
		want := []int32{4351, 4352, 4353}
		for i, v := range dints {
			if v != want[i] {
				t.Errorf("TestDintArr[0][%d] = %v, want %v", i, v, want[i])
			}
		}
	}
	realVal, err := result.Float32("TestReal")
	if err != nil {
		t.Errorf("Error getting float32: %v", err)
	} else if realVal != float32(93.45) {
		t.Errorf("TestReal = %v, want 93.45", realVal)
	}
	buf.Close()
}
