package gologix

import (
	"errors"
	"fmt"
)

type CIPError struct {
	Code     byte
	Extended uint16
}

func (err *CIPError) Error() string {
	ec := fmt.Sprintf("error %X%X: ", err.Code, err.Extended)
	switch err.Code {
	case 0x03:
		return ec + "Bad parameter, size > 12 or size greater than size of element."
	case 0x04:
		return ec + "A syntax error was detected decoding the Request Path."
	case 0x05:
		return ec + "Request Path destination unknown: probably instance number is not present."
	case 0x06:
		return ec + "Insufficient Packet Space: Not enough room in the response buffer for all the data."
	case 0x0A:
		return ec + "Attribute list error, generally attribute not supported. The status of the unsupported attribute is 0x14."
	case 0x10:
		switch err.Extended {
		case 0x2101:
			return ec + "Device state conflict: keyswitch position: The requestor is changing force information in HARD RUN mode."
		case 0x2802:
			return ec + "Device state conflict: Safety Status: The controller is in a state in which Safety Memory cannot be modified."
		}
	case 0x13:
		return ec + "Insufficient Request Data: Data too short for expected param."
	case 0x1C:
		return ec + "Attribute List Shortage: The list of attribute numbers was too few for the number of attributes parameter"
	case 0x26:
		return ec + "The Request Path Size received was shorter or longer than expected."
	case 0xFF:
		switch err.Extended {
		case 0x2104:
			return ec + "General Error: Offset is beyond end of the requested tag."
		case 0x2105:
			return ec + "General Error: Number of Elements or Byte Offset is beyond the end of the requested tag."
		case 0x2107:
			return ec + "General Error: Tag type used n request does not match the target tag's data type."
		}
	}

	return ec + "Unknown Error"
}

func NewMultiError(err error) MultiError {
	if err != nil {
		return MultiError{[]error{err}}
	}
	return MultiError{[]error{}}
}

// combine multiple errors together.
type MultiError struct {
	errs []error
}

func (e *MultiError) Add(err error) {
	e.errs = append(e.errs, err)
}

func (e MultiError) Error() string {
	err_str := ""
	for err := range e.errs {
		err_str = fmt.Sprintf("%s: %s", err_str, err)
	}
	return err_str
}

func (e MultiError) Unwrap() error {
	if len(e.errs) > 2 {
		e.errs = e.errs[1:]
		return e
	}
	if len(e.errs) == 2 {
		return e.errs[1]
	}
	if len(e.errs) == 1 {
		return e.errs[0]
	}
	return nil
}

func (e MultiError) Is(target error) bool {
	return errors.Is(e.errs[0], target)
}
