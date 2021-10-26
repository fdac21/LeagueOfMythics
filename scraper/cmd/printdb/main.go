package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/dgraph-io/badger/v3"
)

func replaceAndPrint(s string) {
	parts := strings.SplitN(s, "___", 2)
	switch parts[0] {
	case "CHALLENGERI":
		fmt.Printf("a %v\n", parts[1])
		break
	case "GRANDMASTERI":
		fmt.Printf("b %v\n", parts[1])
		break
	case "MASTERI":
		fmt.Printf("c %v\n", parts[1])
		break
	case "DIAMONDI":
		fmt.Printf("d %v\n", parts[1])
		break
	case "DIAMONDII":
		fmt.Printf("e %v\n", parts[1])
		break
	case "DIAMONDIII":
		fmt.Printf("f %v\n", parts[1])
		break
	case "DIAMONDIV":
		fmt.Printf("g %v\n", parts[1])
		break
	case "PLATINUMI":
		fmt.Printf("h %v\n", parts[1])
		break
	case "PLATINUMII":
		fmt.Printf("i %v\n", parts[1])
		break
	case "PLATINUMIII":
		fmt.Printf("j %v\n", parts[1])
		break
	case "PLATINUMIV":
		fmt.Printf("k %v\n", parts[1])
		break
	case "GOLDI":
		fmt.Printf("l %v\n", parts[1])
		break
	case "GOLDII":
		fmt.Printf("m %v\n", parts[1])
		break
	case "GOLDIII":
		fmt.Printf("n %v\n", parts[1])
		break
	case "GOLDIV":
		fmt.Printf("o %v\n", parts[1])
		break
	case "SILVERI":
		fmt.Printf("p %v\n", parts[1])
		break
	case "SILVERII":
		fmt.Printf("q %v\n", parts[1])
		break
	case "SILVERIII":
		fmt.Printf("r %v\n", parts[1])
		break
	case "SILVERIV":
		fmt.Printf("s %v\n", parts[1])
		break
	}
}

func main() {

	db, err := badger.Open(badger.DefaultOptions("./databasedata"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		i := 0
		for it.Rewind(); it.Valid(); it.Next() {
			i++
			if (i % 20) != 0 {
				continue
			}
			replaceAndPrint(string(it.Item().Key()))

		}

		return nil
	})
}
