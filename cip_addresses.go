package gologix

type CIPAddress byte

func (p CIPAddress) Bytes() []byte {
	return []byte{byte(p)}
}
