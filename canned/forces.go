package canned

import (
	"fmt"

	"github.com/danomagnum/gologix"
)

type ForceStatus uint16

func (f ForceStatus) Enabled() bool {
	return f&1 != 0
}
func (f ForceStatus) Exist() bool {
	return f&1 != 0
}

func GetForces(client *gologix.Client) (ForceStatus, error) {
	item, err := client.GetAttrList(gologix.CipObject_IOClass, 0, 9)
	if err != nil {
		return 0, err
	}

	if item == nil {
		return 0, fmt.Errorf("item is nil")
	}

	response := struct {
		ID     uint16
		Status uint16
		Forces ForceStatus
	}{}

	err = item.DeSerialize(&response)
	if err != nil {
		return 0, fmt.Errorf("could not deserialize response: %v", err)
	}

	return response.Forces, nil

}
