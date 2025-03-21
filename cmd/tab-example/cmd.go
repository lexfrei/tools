package main

import "log"

func main() {
	for i := range 100 {
		log.Printf("%d\tabs\n", i)
	}
}
