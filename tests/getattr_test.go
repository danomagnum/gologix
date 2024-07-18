package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestGetAttrSingle(t *testing.T) {

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

			//VendorID (UINT, attribute 1)
			i, err := client.GetAttrSingle(gologix.CipObject_Identity, 1, 1)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
			val, err := i.Uint16()
			if err != nil {
				t.Errorf("problem reading attr 1: %v", err)
				return
			}
			log.Printf("vendor: 0x%X", val)

			if val != 0x01 {
				t.Errorf("vendor ID should have been 0x01 but was 0x%X", val)
			}

			//
			//DeviceType (UINT, attribute 2)
			i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 2)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
			val, err = i.Uint16()
			if err != nil {
				t.Errorf("problem reading attr 2: %v", err)
				return
			}
			log.Printf("device type: 0x%X", val)
			if val != 0x0E {
				t.Errorf("Device type should have been 0x0E (PLC) but was 0x%X", val)
			}

			//ProductCode (UINT, attribute 3)
			i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 3)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
			val, err = i.Uint16()
			if err != nil {
				t.Errorf("problem reading attr 3: %v", err)
				return
			}
			log.Printf("product code: 0x%X", val)
			if val != tc.ProductCode {
				t.Errorf("Product Code should have been %d (Compact Logix?) but was %d", tc.ProductCode, val)
			}

			//MajorRevision (USINT, attribute 4)
			//MinorRevision (USINT, attribute 4)
			i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 4)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
			major, err := i.Byte()
			if err != nil {
				t.Errorf("problem reading attr 4 major version: %v", err)
				return
			}
			minor, err := i.Byte()
			if err != nil {
				t.Errorf("problem reading attr 4 minor version: %v", err)
				return
			}
			log.Printf("Version:%d.%d", major, minor)
			if major != tc.SoftwareVersionMajor || minor != tc.SoftwareVersionMinor {
				t.Errorf("Version should have been %d.%d but was %d.%d", tc.SoftwareVersionMajor, tc.SoftwareVersionMinor, major, minor)
			}

			//Status (UINT, attribute 5)
			i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 5)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
			val, err = i.Uint16()
			if err != nil {
				t.Errorf("problem reading attr 5: %v", err)
				return
			}
			log.Printf("status: 0x%X", val)

			//SerialNumber (UDINT, attribute 6)
			i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 6)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
			serial, err := i.Uint32()
			if err != nil {
				t.Errorf("problem reading attr 6: %v", err)
				return
			}
			log.Printf("serial: %d", serial)
			if serial != tc.SerialNumber {
				t.Errorf("Serial # should have been %d but was %d", tc.SerialNumber, serial)
			}

			//ProductName (STR_32, attribute 7)
			i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 7)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
			_, _ = i.Byte()
			name := string(i.Rest())
			log.Printf("ProductName: %s", name)
			if name != tc.ProductName {
				t.Errorf("Product Name should have been \n'%s' but was \n'%s'", tc.ProductName, name)
			}

			// test multi-byte instance
			_, err = client.GetAttrSingle(gologix.CipObject_Assembly, 0x303, 0x03)
			if err != nil {
				t.Errorf("problem reading items: %v", err)
			}
		})
	}

}

func TestGetCtrlProps(t *testing.T) {
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

			props, err := client.GetControllerPropList()
			if err != nil {
				t.Errorf("problem getting controller prop list: %v", err)
			}
			log.Printf("Props: %+v", props)
		})
	}
}

func TestGetAttrList(t *testing.T) {

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

			//VendorID (UINT, attribute 1)
			i, err := client.GetAttrList(gologix.CipObject_Identity, 1,
				1, 2, 3, 4, 6, 7) // properties
			if err != nil {
				t.Errorf("problem reading items: %v", err)
				return
			}

			type CtrlAttrStruct struct {
				VendorID_AttrID uint16
				VendorID_Status uint16
				VendorID        uint16

				DeviceType_AttrID uint16
				DeviceType_Status uint16
				DeviceType        uint16

				ProdCode_AttrID uint16
				ProdCode_Status uint16
				ProdCode        uint16

				Revision_AttrID uint16
				Revision_Status uint16
				MajorRevision   byte
				MinorRevision   byte

				SerialNo_AttrID uint16
				SerialNo_Status uint16
				SerialNo        uint32

				ProductName_AttrID uint16
				ProductName_Status uint16
				Pad                byte
			}

			var results CtrlAttrStruct
			err = i.DeSerialize(&results)
			if err != nil {
				t.Errorf("problem deserializing: %v", err)
			}

			if results.VendorID != 0x01 {
				t.Errorf("vendor ID should have been 0x01 but was 0x%X", results.VendorID)
			}

			if results.DeviceType != 0x0E {
				t.Errorf("Device type should have been 0x0E (PLC) but was 0x%X", results.DeviceType)
			}
			if results.ProdCode != tc.ProductCode {
				t.Errorf("Product Code should have been %d (Compact Logix?) but was %d", tc.ProductCode, results.ProdCode)
			}

			if results.MajorRevision != tc.SoftwareVersionMajor || results.MinorRevision != tc.SoftwareVersionMinor {
				t.Errorf("Version should have been %d.%d but was %d.%d", tc.SoftwareVersionMajor, tc.SoftwareVersionMinor, results.MajorRevision, results.MinorRevision)
			}

			if results.SerialNo != tc.SerialNumber {
				t.Errorf("Serial # should have been %d but was %d", tc.SerialNumber, results.SerialNo)
			}

			name := string(i.Rest())
			log.Printf("ProductName: %s", name)
			if name != tc.ProductName {
				t.Errorf("Product Name should have been \n'%s' but was \n'%s'", tc.ProductName, name)
			}
		})
	}
}
