package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestGetAttrSingle(t *testing.T) {

	client := gologix.NewClient("192.168.2.241")
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
	if val != 0x97 {
		t.Errorf("Product Code should have been 0x97 (Compact Logix?) but was 0x%X", val)
	}

	//MajorRevision (USINT, attribute 4)
	//MinorRevision (USINT, attribute 4)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 4)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	major, err := i.Byte()
	if err != nil {
		t.Errorf("problem reading attr 3 major version: %v", err)
		return
	}
	minor, err := i.Byte()
	if err != nil {
		t.Errorf("problem reading attr 3 minor version: %v", err)
		return
	}
	log.Printf("Version:%d.%d", major, minor)
	if major != 33 || minor != 11 {
		t.Errorf("Version should have been 33.11 but was %d.%d", major, minor)
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
	if serial != 3223037449 {
		t.Errorf("Serial # should have been 3223037449 but was %d", serial)
	}

	//ProductName (STR_32, attribute 7)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 7)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	_, _ = i.Byte()
	name := string(i.Rest())
	log.Printf("ProductName: %s", name)
	wantName := "1769-L27ERM-QxC1B/A LOGIX5327ERM"
	if name != wantName {
		t.Errorf("Product Name should have been \n'%s' but was \n'%s'", wantName, name)
	}

}

func TestGetCtrlProps(t *testing.T) {
	client := gologix.NewClient("192.168.2.241")
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
}

func TestGetAttrList(t *testing.T) {

	client := gologix.NewClient("192.168.2.241")
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
	if results.ProdCode != 0x97 {
		t.Errorf("Product Code should have been 0x97 (Compact Logix?) but was 0x%X", results.ProdCode)
	}

	if results.MajorRevision != 33 || results.MinorRevision != 11 {
		t.Errorf("Version should have been 33.11 but was %d.%d", results.MajorRevision, results.MinorRevision)
	}

	if results.SerialNo != 3223037449 {
		t.Errorf("Serial # should have been 3223037449 but was %d", results.SerialNo)
	}

	name := string(i.Rest())
	log.Printf("ProductName: %s", name)
	wantName := "1769-L27ERM-QxC1B/A LOGIX5327ERM"
	if name != wantName {
		t.Errorf("Product Name should have been \n'%s' but was \n'%s'", wantName, name)
	}

}
