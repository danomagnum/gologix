package main

import "log"

func main() {

	plc := &PLC{IPAddress: "192.168.2.241"}
	plc.Connect()
	//plc.read_single("program:Shed.Temp1", CIPTypeREAL, 1)
	ShedTemp, err := Read[float32](plc, "program:Shed.Temp1")
	if err != nil {
		log.Printf("Problem reading ShedTemp. %v", err)
	}
	log.Printf("Shed Temp: %v", ShedTemp)

}
