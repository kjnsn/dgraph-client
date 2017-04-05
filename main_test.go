package dgraph

import (
	"flag"
	"fmt"
	"log"
	"testing"
)

var address = flag.String("address", "127.0.0.1:8080", "Dgraph server address")

func TestXXX(t *testing.T) {
	conn, err := Connect(*address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	xd, err := conn.Query(`
        mutation {
            set {
                <piano> <name> "Piano" .
                <guitar> <name> "Guitar" .
                <bob> <name> "Bob" .
                <bob> <age> "27"^^<xs:int> .
                <bob> <sign> "Gemini" .
                <bob> <plays> <piano> .
                <mike> <name> "Mike" .
                <mike> <age> "42"^^<xs:int> .
                <mike> <sign> "Scorpio" .
                <mike> <plays> <guitar> .
                <clara> <name> "Clara" .
                <clara> <age> "32"^^<xs:int> .
                <clara> <sign> "Taurus" .
                <clara> <follows> <bob> .
                <clara> <follows> <mike> .
                <bob> <follow> <clara> .
            }
        }
        query {
            me(id: clara) {
                _uid_
                name
                age
                sign
                follows {
                    _uid_
                    name
                    age
                    sign
                    plays {
                        name
                    }
                }
            }
        }
    `)

	if err != nil {
		log.Fatal(err)
	}

	// xd := xdgraph.ReadResponse(resp)

	fmt.Printf("Clara's UID: %d\n", xd.Attribute("me").UID())
	fmt.Printf("Clara's name: %s\n", xd.Attribute("me").Property("name").ToString())
	fmt.Printf("Clara's sign (using First()): %s\n", xd.First().Property("sign").ToString())

	// Using the model above, we cannot use types until #594 is fixed and released
	// See https://github.com/dgraph-io/dgraph/issues/594
	fmt.Printf("Clara's age: %d // If 0 see https://github.com/dgraph-io/dgraph/issues/594\n", xd.Attribute("me").Property("age").ToInt())

	fmt.Printf("Clara follows:\n")
	xd.First().Attribute("follows").Each(func(r Response) {
		fmt.Printf(" - %s (uid %d)\n", r.Property("name").ToString(), r.UID())
	})

	fmt.Printf("Signs of people Clara follows:\n")
	for _, v := range xd.First().Attribute("follows").Properties("sign") {
		fmt.Printf(" - %s\n", v.ToString())
	}

	fmt.Printf("Instruments played by people Clara follows:\n")
	xd.First().Attribute("follows").Attribute("plays").Each(func(r Response) {
		fmt.Printf(" - %s\n", r.Property("name").ToString())
	})

	fmt.Printf("Which instrument does each person Clara follows plays:\n")
	xd.First().Attribute("follows").Each(func(r Response) {
		if r.Attribute("plays").IsNil() {
			// This person does not play any instrument
			return
		}
		fmt.Printf(" - %s plays %s\n", r.Property("name").ToString(), r.Attribute("plays").Property("name").ToString())
	})

	fmt.Printf("Json (full graph):\n%s\n", xd.JSON())
	fmt.Printf("Json (sub graph):\n%s\n", xd.Attribute("me").Attribute("follows").JSON())
}
