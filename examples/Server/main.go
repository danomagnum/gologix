package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danomagnum/gologix"
)

func main() {

	// one memory based tag provider at slot 0 on the "backplane"
	p1 := gologix.MapTagProvider{}
	r := gologix.PathRouter{}
	path1, err := gologix.ParsePath("1,0")
	if err != nil {
		fmt.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.AddHandler(path1.Bytes(), &p1)

	// a different memory based tag provider at slot 1 on the "backplane"
	p2 := gologix.MapTagProvider{}
	path2, err := gologix.ParsePath("1,1")
	if err != nil {
		fmt.Printf("problem parsing path. %v", err)
		os.Exit(1)
	}
	r.AddHandler(path2.Bytes(), &p2)

	s := gologix.NewServer(&r)
	go s.Serve()

	t := time.NewTicker(time.Second * 5)
	for {
		<-t.C
		p1.Mutex.Lock()
		log.Printf("Data 1: %v", p1.Data)
		p1.Mutex.Unlock()

		p2.Mutex.Lock()
		log.Printf("Data 2: %v", p2.Data)
		p2.Mutex.Unlock()

	}

}
