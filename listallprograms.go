package gologix

import (
	"strings"
)

func (client *Client) ListAllPrograms() error {
	if verbose {
		client.Logger.Printf("listing all programs")
	}

	// for generic messages we need to create the cip path ourselves.  The serialize function can be used to do this.
	path, err := Serialize(CipObject_Programs, CIPInstance(0))
	if err != nil {
		client.Logger.Printf("could not serialize path: %v", err)
		return err
	}

	number_of_attr_to_receive := 1
	attr_28_program_name := 28
	msg_data, err := Serialize(uint16(number_of_attr_to_receive), uint16(attr_28_program_name))
	if err != nil {
		client.Logger.Printf("could not serialize message data: %v", err)
		return err
	}

	resp, err := client.GenericCIPMessage(CIPService_GetInstanceAttributeList, path.Bytes(), msg_data.Bytes())
	if err != nil {
		client.Logger.Printf("problem reading programs: %v", err)
		return err
	}

	results := make(map[string]*KnownProgram)

	for {
		var hdr listprograms_resp_header
		err = resp.DeSerialize(&hdr)
		if err != nil {
			client.Logger.Printf("got last item")
			break
		}

		// read the name
		name := make([]byte, hdr.NameLen)
		err = resp.DeSerialize(&name)
		if err != nil {
			client.Logger.Printf("could not read name: %v", err)
			return err
		}

		// convert the name to a string
		namestr := strings.TrimSpace(string(name))
		results[namestr] = &KnownProgram{ID: CIPInstance(hdr.InstanceID), Name: namestr}
	}

	client.KnownPrograms = results

	return nil
}

type listprograms_resp_header struct {
	InstanceID uint16
	Padding    uint16
	NameLen    uint32
}
