package main

func main() {

	plc := PLC{IPAddress: "192.168.2.241"}
	plc.Connect()
	plc.read_single("program:Shed.Temp1", CIPTypeREAL, 1)

}
