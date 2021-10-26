package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/badger/v3"
	lolapi "github.com/fdac21/LeagueOfMythics/scraper/v2"
)

func main() {

	apiKey := os.Getenv("RGAPIKEY")
	if apiKey == "" {
		log.Fatal("No key found in $RGAPIKEY")
	}
	c := lolapi.NewApiClient(apiKey)
	db, err := badger.Open(badger.DefaultOptions("databasedata"))
	if err != nil {
		log.Fatalln(err)
	}

	type rankdiv struct {
		r lolapi.Rank
		d lolapi.Division
	}

	rds := []rankdiv{
		{r: lolapi.RankS, d: lolapi.Div1},
		{r: lolapi.RankS, d: lolapi.Div2},
		{r: lolapi.RankS, d: lolapi.Div3},
		{r: lolapi.RankS, d: lolapi.Div4},
	}
	for _, rd := range rds {
		log.Println("Currently on " + lolapi.Ranks[rd.r] + " " + lolapi.Divisions[rd.d])
		for i := 1; ; i++ {
			if i%10 == 0 {
				log.Printf("on page %v\n", i)
			}
			players, err := c.GetLeagueEntries(rd.r, rd.d, i)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintf(os.Stderr, "died on %v%v page %v\n", lolapi.Ranks[rd.r],
					lolapi.Divisions[rd.d], i)
				os.Exit(1)
			}
			if len(players) == 0 {
				break
			}
			// Insert into db
			for _, p := range players {
				db.Update(func(txn *badger.Txn) error {
					// Key on rank so its easy to know who is what rank
					k := []byte(lolapi.Ranks[rd.r] + lolapi.Divisions[rd.d] + "___" + p.SummonerId)
					v := []byte(p.SummonerName)
					_, err := txn.Get(k)
					if err == badger.ErrKeyNotFound {
						if err = txn.Set(k, v); err != nil {
							log.Fatal("insert error", err)
						}
					}
					return nil
				})
			}
		}
	}
	fmt.Println("done")
}
