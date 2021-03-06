package main

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestInstructionsTry(t *testing.T) {
	i := &instructions{
		Tries: 0,
	}

	if i.Try() != 1 {
		t.Logf("Try() must return 1 if a zero value is set")
	}

	i.Tries = 100
	if i.Try() != 100 {
		t.Logf("Try() must return the value specified if not set to zero")
	}
}

func TestInstructionsID(t *testing.T) {
	// I'm going to cheat a bit here. Instead of manually loading it up, I'm going to
	// load the entire dojo ... but in doing so, will just hardcode the id.
	// If I run into any more issues, I will simply add them to the instruction suiteA
	// test suite
	oogway := Dojo(
		&oogway{
			ChecksDir:       "t/suiteA",
			CheckInterval:   0 * time.Second,
			ChecksExtension: "com",
		})

	oogway.loadChecks()

	if oogway.checks["kcmerrill.com"].id() != "TXkgZGVzY3JpcHRpb24gd291bGQgZ28gaGVyZTFoMG0wczBzNHRvdWNoIC90bXAvc3VpdGVBLmNoZWNrLmtjbWVycmlsbC5jb20wMHRvdWNoIC90bXAvc3VpdGVBLndhcm5pbmcua2NtZXJyaWxsLmNvbTB0b3VjaCAvdG1wL3N1aXRlQS5jcml0aWNhbC5rY21lcnJpbGwuY29tMHRvdWNoIC90bXAvc3VpdGVBLmZpeC5rY21lcnJpbGwuY29tMnRvdWNoIC90bXAvc3VpdGVBLnJlY292ZXIua2NtZXJyaWxsLmNvbTA=" {
		fmt.Println(oogway.checks["kcmerrill.com"].id())
		log.Fatalf("Id() return has changed. It may need to be added to id()")
	}
}

func TestInstructionsEvery(t *testing.T) {
	i := &instructions{}

	if i.every() != 30*time.Second {
		log.Fatalf("Default value for every() should be 30s")
	}

	hour := 1 * time.Hour
	i.Every = hour
	if hour != i.every() {
		log.Fatalf("Every was set to 1h, expected 1h")
	}
}
