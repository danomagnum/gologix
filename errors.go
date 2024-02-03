package gologix

import (
	"errors"
	"fmt"
)

// Represents a CIP error.
type CIPError struct {
	Code     byte
	Extended uint16
}

func (err *CIPError) Error() string {
	ec := fmt.Sprintf("error %X%X: ", err.Code, err.Extended)
	switch err.Code {
	case 0x00:
		return ec + " no error?  This shouldn't happen :/"
	case 0x01:
		return ec + " connection failure"
	case 0x02:
		return ec + " resource unavailable"
	case 0x03:
		return ec + " bad parameter, size > 12 or size greater than size of element."
	case 0x04:
		return ec + " a syntax error was detected decoding the request path"
	case 0x05:
		return ec + " request path destination unknown: probably instance number is not present"
	case 0x06:
		return ec + " insufficient packet space: not enough room in the response buffer for all the data"
	case 0x07:
		return ec + " connection lost"
	case 0x08:
		return ec + " service is not supported for the object/instance"
	case 0x09:
		return ec + " could not write attribute data - possibly in valid or wrong type"
	case 0x0A:
		return ec + " attribute list error, generally attribute not supported. the status of the unsupported attribute is 0x14"
	case 0x10:
		switch err.Extended {
		case 0x2101:
			return ec + " device state conflict: keyswitch position: the requestor is changing force information in HARD RUN mode"
		case 0x2802:
			return ec + " device state conflict: safety status: the controller is in a state in which safety memory cannot be modified"
		}
	case 0x13:
		return ec + " insufficient Request Data: Data too short for expected param"
	case 0x16:
		return ec + " object does not exist"
	case 0x1A:
		return ec + " routing failure: request too large"
	case 0x1B:
		return ec + " routing failure: response too large"
	case 0x1C:
		return ec + " attribute list shortage: the list of attribute numbers was too few for the number of attributes parameter"
	case 0x26:
		return ec + " the request path size received was shorter or longer than expected."
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

func newMultiError(err error) multiError {
	if err != nil {
		return multiError{[]error{err}}
	}
	return multiError{[]error{}}
}

// combine multiple errors together.
type multiError struct {
	errs []error
}

func (e *multiError) Add(err error) {
	e.errs = append(e.errs, err)
}

func (e multiError) Error() string {
	err_str := ""
	for i := range e.errs {
		err_str = fmt.Sprintf("%s: %s", err_str, e.errs[i])
	}
	return err_str
}

func (e multiError) Unwrap() error {
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

func (e multiError) Is(target error) bool {
	return errors.Is(e.errs[0], target)
}
