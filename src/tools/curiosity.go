// +build manually

package main

import (
	"encoding/json"
	"fmt"
	"hstools"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/boltdb/bolt"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) != 3 {
		log.Fatal("usage: curiosity keys.db {ip,key,colocated}")
	}

	keysDB, err := hstools.OpenKeysDb(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	keysDB.View(func(tx *bolt.Tx) error {
		switch os.Args[2] {
		case "ip":
			c := tx.Bucket([]byte("IPs")).Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var res hstools.IPMeta
				if err := json.Unmarshal(v, &res); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%d %s\n", len(res.Keys), k)
			}
		case "keys":
			c := tx.Bucket([]byte("Keys")).Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var res hstools.KeyMeta
				if err := json.Unmarshal(v, &res); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%d %s\n", len(res.IPs), hstools.ToHex(k))
			}
		case "colocated":
			c := tx.Bucket([]byte("Keys")).Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var res hstools.KeyMeta
				colocated := make(map[string]struct{})
				if err := json.Unmarshal(v, &res); err != nil {
					log.Fatal(err)
				}
				for _, ip := range res.IPs {
					ipMetaJSON := tx.Bucket([]byte("IPs")).Get([]byte(ip))
					var ipMeta hstools.IPMeta
					if err := json.Unmarshal(ipMetaJSON, &ipMeta); err != nil {
						log.Fatal(err)
					}
					for _, key := range ipMeta.Keys {
						colocated[hstools.ToHex(key)] = struct{}{}
					}
				}
				fmt.Printf("%d keys on %d IPs - %s %v\n",
					len(colocated), len(res.IPs), hstools.ToHex(k), res.IPs)
			}
		default:
			log.Fatal("usage: curiosity keys.db {ip,key,colocated}")
		}

		return nil
	})

	if err := keysDB.Close(); err != nil {
		log.Fatal(err)
	}
}
