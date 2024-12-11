package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestWrite(t *testing.T) {
	tcs := getTestConfig()
	for _, tc := range tcs.PlcList {
		t.Run(tc.PlcAddress, func(t *testing.T) {
			client := gologix.NewClient(tc.PlcAddress)
			err := client.Connect()
			if err != nil {
				t.Error(err)
				return
			}
			defer func() {
				err := client.Disconnect()
				if err != nil {
					t.Errorf("problem disconnecting. %v", err)
				}
			}()

			write_and_check[int16](t, client, "program:gologix_tests.writeint", 0, 10, 20)
			write_and_check[int32](t, client, "program:gologix_tests.writedint", 0, 10, 20)
			//write_and_check[byte](t, client, "writesint", 0, 10, 20)

			write_and_check[bool](t, client, "program:gologix_tests.writebool", false, true, false, true, false)

			write_and_check[float32](t, client, "program:gologix_tests.writereal", 0, 12.4, 5353.0281, 4)

			write_and_check[int32](t, client, "program:gologix_tests.writeudt.field1", 0, 5, 281, 46)
			write_and_check[float32](t, client, "program:gologix_tests.writeudt.field2", 0, 12.4, 5353.0281, 4)

			write_and_check[string](t, client, "program:gologix_tests.MultiWriteString", "a", "b", "c")
		})

	}
}

func write_and_check[T gologix.GoLogixTypes](t *testing.T, client *gologix.Client, tag string, values ...T) {
	var err error
	var have T
	t.Run(tag, func(t *testing.T) {
		for _, want := range values {
			err = client.Write(tag, want)
			if err != nil {
				t.Errorf("problem writing. %v", err)
			}
			err = client.Read(tag, &have)
			if err != nil {
				t.Errorf("problem reading back. %v", err)
			}
			if have != want {
				t.Errorf("want %v. have %v", want, have)
			}

		}
	})

}

func TestMultiWrite(t *testing.T) {
	tcs := getTestConfig()
	for _, tc := range tcs.PlcList {
		t.Run(tc.PlcAddress, func(t *testing.T) {
			client := gologix.NewClient(tc.PlcAddress)
			err := client.Connect()
			if err != nil {
				t.Error(err)
				return
			}
			defer func() {
				err := client.Disconnect()
				if err != nil {
					t.Errorf("problem disconnecting. %v", err)
				}
			}()

			write_map := make(map[string]any)
			write_map["program:gologix_tests.MultiWriteInt"] = int16(123)
			write_map["program:gologix_tests.MultiWriteReal"] = float32(456.7)
			write_map["program:gologix_tests.MultiWriteDint"] = int32(891011)
			write_map["program:gologix_tests.MultiWriteString"] = "Test String!"
			write_map["program:gologix_tests.MultiWriteBool"] = true

			err = client.WriteMap(write_map)
			if err != nil {
				log.Printf("error writing to multiple tags at once: %v", err)
			}

			vals, err := client.ReadList([]string{"program:gologix_tests.MultiWriteInt",
				"program:gologix_tests.MultiWriteReal",
				"program:gologix_tests.MultiWriteDint",
				"program:gologix_tests.MultiWriteString",
				"program:gologix_tests.MultiWriteBool"},
				[]gologix.CIPType{gologix.CIPTypeINT, gologix.CIPTypeREAL, gologix.CIPTypeDINT, gologix.CIPTypeSTRING, gologix.CIPTypeBOOL},
				[]int{1, 1, 1, 1, 1})

			if err != nil {
				t.Errorf("problem reading tags back: %v", err)
			}

			// verify the values read back correctly.

			i16 := vals[0].(int16)
			if i16 != 123 {
				t.Errorf("Int read incorrectly. wanted %d got %d", 123, i16)
			}

			f32 := vals[1].(float32)
			if f32 != 456.7 {
				t.Errorf("Real read incorrectly. wanted %f got %f", 456.7, f32)
			}

			i32 := vals[2].(int32)
			if i32 != 891011 {
				t.Errorf("DINT read incorrectly. wanted %d got %d", 891011, i32)
			}

			s := vals[3].(string)
			if s != "Test String!" {
				t.Errorf("String read incorrectly. wanted %s got %s", "Test String!", s)
			}

			b := vals[4].(bool)
			if b != true {
				t.Errorf("BOOL read incorrectly. wanted %v got %v", true, b)
			}

			write_map["program:gologix_tests.MultiWriteInt"] = int16(321)
			write_map["program:gologix_tests.MultiWriteReal"] = float32(7.654)
			write_map["program:gologix_tests.MultiWriteDint"] = int32(110198)
			write_map["program:gologix_tests.MultiWriteString"] = "String Test!"
			write_map["program:gologix_tests.MultiWriteBool"] = false

			err = client.WriteMap(write_map)
			if err != nil {
				log.Printf("error writing to multiple tags at once: %v", err)
			}

			vals, err = client.ReadList([]string{"program:gologix_tests.MultiWriteInt",
				"program:gologix_tests.MultiWriteReal",
				"program:gologix_tests.MultiWriteDint",
				"program:gologix_tests.MultiWriteString",
				"program:gologix_tests.MultiWriteBool"},
				[]gologix.CIPType{gologix.CIPTypeINT, gologix.CIPTypeREAL, gologix.CIPTypeDINT, gologix.CIPTypeSTRING, gologix.CIPTypeBOOL}, []int{1, 1, 1, 1, 1})

			if err != nil {
				t.Errorf("problem reading tags back: %v", err)
			}

			// verify the values read back correctly.

			i16 = vals[0].(int16)
			if i16 != 321 {
				t.Errorf("Int read incorrectly. wanted %d got %d", 321, i16)
			}

			f32 = vals[1].(float32)
			if f32 != 7.654 {
				t.Errorf("Real read incorrectly. wanted %f got %f", 7.654, f32)
			}

			i32 = vals[2].(int32)
			if i32 != 110198 {
				t.Errorf("DINT read incorrectly. wanted %d got %d", 110198, i32)
			}

			s = vals[3].(string)
			if s != "String Test!" {
				t.Errorf("String read incorrectly. wanted %s got %s", "String Test!", s)
			}

			b = vals[4].(bool)
			if b != false {
				t.Errorf("BOOL read incorrectly. wanted %v got %v", false, b)
			}
		})
	}
}

func TestWriteUdt1(t *testing.T) {

	type TestUDT struct {
		Field1 int32
		Field2 float32
	}

	tcs := getTestConfig()
	for _, tc := range tcs.PlcList {
		t.Run(tc.PlcAddress, func(t *testing.T) {

			client := gologix.NewClient(tc.PlcAddress)
			err := client.Connect()
			if err != nil {
				t.Error(err)
				return
			}
			defer func() {
				err := client.Disconnect()
				if err != nil {
					t.Errorf("problem disconnecting. %v", err)
				}
			}()

		})
	}
}

func TestWriteUdt(t *testing.T) {

	type TestUDT struct {
		Field1 int32
		Field2 float32
	}

	type TestUDT2 struct {
		Field1 int32
		Flag1  bool
		Flag2  bool
		Field2 int32
	}

	type NestedUDT struct {
		Field1  int32
		SubStr1 TestUDT
		SubStr2 TestUDT2
		Field2  int32
	}

	tcs := getTestConfig()
	for _, tc := range tcs.PlcList {
		t.Run(tc.PlcAddress, func(t *testing.T) {

			client := gologix.NewClient(tc.PlcAddress)
			err := client.Connect()
			if err != nil {
				t.Error(err)
				return
			}
			defer func() {
				err := client.Disconnect()
				if err != nil {
					t.Errorf("problem disconnecting. %v", err)
				}
			}()

			var t1 TestUDT
			var t1_readback TestUDT

			t1.Field1 = 123
			t1.Field2 = 456.7

			err = client.Write("program:gologix_tests.writeudt", t1)
			if err != nil {
				t.Errorf("problem writing first time: %v", err)
			}

			err = client.Read("program:gologix_tests.writeudt", &t1_readback)
			if err != nil {
				t.Errorf("problem reading first time: %v", err)
			}
			if t1_readback != t1 {
				t.Errorf("want %v. have %v", t1, t1_readback)
			}

			t1.Field1 = 321
			t1.Field2 = 7.654

			err = client.Write("program:gologix_tests.writeudt", t1)
			if err != nil {
				t.Errorf("problem writing second time: %v", err)
			}

			err = client.Read("program:gologix_tests.writeudt", &t1_readback)
			if err != nil {
				t.Errorf("problem reading second time: %v", err)
			}
			if t1_readback != t1 {
				t.Errorf("want %v. have %v", t1, t1_readback)
			}

			var t2 TestUDT2
			var t2_readback TestUDT2

			t2.Field1 = 123
			t2.Flag1 = true
			t2.Flag2 = false
			t2.Field2 = 456

			err = client.Write("program:gologix_tests.writeudt2", t2)
			if err != nil {
				t.Errorf("problem writing first time: %v", err)
			}

			err = client.Read("program:gologix_tests.writeudt2", &t2_readback)
			if err != nil {
				t.Errorf("problem reading first time: %v", err)
			}
			if t2_readback != t2 {
				t.Errorf("want %v. have %v", t2, t2_readback)
			}

			t2.Field1 = 321
			t2.Flag1 = false
			t2.Flag2 = true
			t2.Field2 = 765

			err = client.Write("program:gologix_tests.writeudt2", t2)
			if err != nil {
				t.Errorf("problem writing second time: %v", err)
			}

			err = client.Read("program:gologix_tests.writeudt2", &t2_readback)
			if err != nil {
				t.Errorf("problem reading second time: %v", err)
			}
			if t2_readback != t2 {
				t.Errorf("want %v. have %v", t2, t2_readback)
			}

			var t3 NestedUDT
			var t3_readback NestedUDT

			t3.Field1 = 123
			t3.SubStr1.Field1 = 456
			t3.SubStr1.Field2 = 1.2
			t3.SubStr2.Field1 = 5
			t3.SubStr2.Flag1 = true
			t3.SubStr2.Flag2 = false
			t3.SubStr2.Field2 = 789
			t3.Field2 = 8

			err = client.Write("program:gologix_tests.writeudt3", t3)
			if err != nil {
				t.Errorf("problem writing first time: %v", err)
			}

			err = client.Read("program:gologix_tests.writeudt3", &t3_readback)
			if err != nil {
				t.Errorf("problem reading first time: %v", err)
			}
			if t3_readback != t3 {
				t.Errorf("want %v. have %v", t3, t3_readback)
			}

			t3.Field1 = 321
			t3.SubStr1.Field1 = 54
			t3.SubStr1.Field2 = 2.1
			t3.SubStr2.Field1 = 22
			t3.SubStr2.Flag1 = false
			t3.SubStr2.Flag2 = true
			t3.SubStr2.Field2 = 987
			t3.Field2 = 65

			err = client.Write("program:gologix_tests.writeudt3", t3)
			if err != nil {
				t.Errorf("problem writing second time: %v", err)
			}

			err = client.Read("program:gologix_tests.writeudt3", &t3_readback)
			if err != nil {
				t.Errorf("problem reading second time: %v", err)
			}
			if t3_readback != t3 {
				t.Errorf("want %v. have %v", t3, t3_readback)
			}

		})
	}
}
