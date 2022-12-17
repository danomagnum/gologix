package gologix

type CIPAddress byte

func (p CIPAddress) Bytes() []byte {
	return []byte{byte(p)}
}

func (p CIPAddress) Len() int {
	return 1
}
