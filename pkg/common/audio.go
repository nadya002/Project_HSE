package main

import (
	"log"
	"os"
)

func main() {
	rf, wf, err := os.Pipe()
	if err != nil {
		log.Println(err)
		return
	}
	rf.Fd()
	wf.Fd()

}
