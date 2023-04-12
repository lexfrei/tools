package main

import (
	"log"
	"math/rand"

	"github.com/cockroachdb/errors"
)

var (
	errFromRuntimeOne   = errors.New("runtime error one")
	errFromRuntimeTwo   = errors.New("runtime error two")
	errFromRuntimeThree = errors.New("runtime error three")
)

func main() {
	err := dummyFunc(1, 2)
	if err != nil {
		switch {
		case errors.Is(err, errFromRuntimeOne):
			log.Println("error from runtime one")
			log.Fatalln(err)
		case errors.Is(err, errFromRuntimeTwo):
			log.Println("error from runtime two")
			log.Fatalln(err)
		case errors.Is(err, errFromRuntimeThree):
			log.Println("error from runtime three")
			log.Fatalln(err)
		default:
			log.Println("unexpected error")
			log.Fatalln(err)
		}
	}
}

// dummyFunc returns a random error from the list of runtime errors.
func dummyFunc(_, _ int) error {
	err := randomError()
	if err != nil {
		return errors.Wrap(err, "wow! expected error")
	}

	return nil
}

// randomError returns a random error from the list of runtime errors.
func randomError() error {
	// Get random 1..3
	//nolint:gomnd,gosec // not interesting
	switch rand.Intn(3) {
	case 1:
		return errFromRuntimeOne
	case 2:
		return errFromRuntimeTwo
	default:
		return errFromRuntimeThree
	}
}
