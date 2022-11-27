package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
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
	plc.ReadAll(1)
	log.Fatal("Done")
	//plc.read_single("program:Shed.Temp1", CIPTypeREAL, 1)
	i := int16(0)
	t0 := time.Now()
	count := int16(1000)
	for {
		//ReadAndPrint[int32](plc, "TestDint")
		//ReadAndPrint[int16](plc, "TestInt")
		//ReadAndPrint[bool](plc, "TestBool")
		//ReadAndPrint[float32](plc, "TestReal")

		err := plc.Write_single("TestInt", i)
		if err != nil {
			log.Printf("error writing. %v", err)
		}
		Read[int16](plc, "TestInt")
		/*
			compare, err := Read[int16](plc, "TestInt")
			if err != nil {
				fmt.Printf("Error with read of TestInt. %v", err)
				panic("Dag, yo.")
			}
			if compare != i {
				fmt.Printf("Read didn't match write!. Expected %v. Got %v", i, compare)
				panic("woops!")
			}
			//fmt.Printf("TestInt: %v", compare)
		*/
		i += 1
		if i == count {

			fmt.Printf("Done in %v. \n", time.Since(t0))
			break
		}
	}

}

func ReadAndPrint[T GoLogixTypes](plc *PLC, path string) {
	value, err := Read[T](plc, path)
	if err != nil {
		log.Printf("Problem reading %s. %v", path, err)
		return
	}
	log.Printf("%s: %v", path, value)
}
