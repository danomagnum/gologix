package main

import "log"

func main() {

	plc := &PLC{IPAddress: "192.168.2.241"}
	plc.Connect()
	//plc.read_single("program:Shed.Temp1", CIPTypeREAL, 1)
	//ReadAndPrint[int32](plc, "TestDint")
	//ReadAndPrint[int16](plc, "TestInt")
	//ReadAndPrint[bool](plc, "TestBool")
	//ReadAndPrint[float32](plc, "TestReal")
	plc.conn.Disconnect()
	/*
		TestDint, err := Read[float32](plc, "TestDint")
		if err != nil {
			log.Printf("Problem reading ShedTemp. %v", err)
		}
		log.Printf("Dint: %v", TestDint)
	*/

}

func ReadAndPrint[T GoLogixTypes](plc *PLC, path string) {
	value, err := Read[T](plc, path)
	if err != nil {
		log.Printf("Problem reading %s. %v", path, err)
		return
	}
	log.Printf("%s: %v", path, value)
}
