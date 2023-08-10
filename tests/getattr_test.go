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
	log.Printf("items: %v", i)

	//
	//DeviceType (UINT, attribute 2)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 2)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	log.Printf("items: %v", i)
	//ProductCode (UINT, attribute 3)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 3)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	log.Printf("items: %v", i)
	//MajorRevision (USINT, attribute 4)
	//MinorRevision (USINT, attribute 4)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 4)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	log.Printf("items: %v", i)

	//Status (UINT, attribute 5)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 5)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	log.Printf("items: %v", i)

	//SerialNumber (UDINT, attribute 6)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 6)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	log.Printf("items: %v", i)

	//ProductName (STR_32, attribute 7)
	i, err = client.GetAttrSingle(gologix.CipObject_Identity, 1, 7)
	if err != nil {
		t.Errorf("problem reading items: %v", err)
	}
	log.Printf("items: %v", i)
	log.Printf("ProductName: %s", string(i[1].Data))

}
