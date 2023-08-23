package gologix_tests

import (
	"log"
	"testing"

	"github.com/danomagnum/gologix"
)

func TestGetAttrSingle(t *testing.T) {

	client := gologix.NewClient("localhost")
	//client := gologix.NewClient("192.168.2.241")
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

	//ProductName (STR_32, attribute 7)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 7)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	log.Printf("ProductName: %s", string(i.Rest()))

}
