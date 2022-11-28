package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	plc := &PLC{IPAddress: "192.168.2.241"}
	plc.Connect()
	defer plc.Disconnect()
	//plc.ReadAll(1)
	//plc.read_single("program:Shed.Temp1", CIPTypeREAL, 1)
	//ReadAndPrint[float32](plc, "program:Shed.Temp1")
	//ReadAndPrint[int32](plc, "TestDint")
	//ReadAndPrint[int16](plc, "TestInt")
	//ReadAndPrint[bool](plc, "TestBool")
	//ReadAndPrint[float32](plc, "TestReal")
	tags := []string{"TestInt", "TestReal"}
	plc.read_multi(tags, CIPTypeDWORD, 1)

}

func ReadAndPrint[T GoLogixTypes](plc *PLC, path string) {
	value, err := Read[T](plc, path)
	if err != nil {
		log.Printf("Problem reading %s. %v", path, err)
		return
	}
	log.Printf("%s: %v", path, value)
}
