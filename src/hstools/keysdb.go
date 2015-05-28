package hstools

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

type KeyMeta struct {
	FirstSeen Hour
	LastSeen  Hour
	IPs       []string
}

type IPMeta struct {
	Keys [][]byte
}

type KeysDB struct {
	db *bolt.DB
}

func OpenKeysDb(filename string) (*KeysDB, error) {
	d := &KeysDB{}
	var err error

	if d.db, err = bolt.Open(filename, 0664, nil); err != nil {
		return nil, err
	}

	d.db.MaxBatchDelay = 5 * time.Second

	if err = d.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("Keys")); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte("IPs")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *KeysDB) Seen(keys []Hash, ips []string, h Hour, wg *sync.WaitGroup) {
	fn := func(tx *bolt.Tx) error {
		for i, k := range keys {
			ip := ips[i]
			b := tx.Bucket([]byte("Keys"))
			oldJSON := b.Get(k[:])
			meta := KeyMeta{
				FirstSeen: h,
				LastSeen:  h,
				IPs:       []string{ip},
			}
			if oldJSON != nil {
				var oldMeta KeyMeta
				if err := json.Unmarshal(oldJSON, &oldMeta); err != nil {
					return err
				}
				if oldMeta.FirstSeen < meta.FirstSeen {
					meta.FirstSeen = oldMeta.FirstSeen
				}
				if oldMeta.LastSeen > meta.LastSeen {
					meta.LastSeen = oldMeta.LastSeen
				}
				for _, oldIP := range oldMeta.IPs {
					if oldIP != ip {
						meta.IPs = append(meta.IPs, oldIP)
					}
				}
			}
			encoded, err := json.Marshal(meta)
			if err != nil {
				return err
			}
			if err := b.Put(k[:], encoded); err != nil {
				return err
			}

			b = tx.Bucket([]byte("IPs"))
			oldJSON = b.Get([]byte(ip))
			ipMeta := IPMeta{
				Keys: [][]byte{k[:]},
			}
			if oldJSON != nil {
				var oldMeta IPMeta
				if err := json.Unmarshal(oldJSON, &oldMeta); err != nil {
					return err
				}
				for _, oldKey := range oldMeta.Keys {
					if !bytes.Equal(oldKey, k[:]) {
						ipMeta.Keys = append(ipMeta.Keys, oldKey)
					}
				}
			}
			encoded, err = json.Marshal(ipMeta)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(ip), encoded); err != nil {
				return err
			}
		}
		return nil
	}
	go func() {
		wg.Add(1)
		if err := d.db.Batch(fn); err != nil {
			log.Fatal(err)
		} else {
			log.Println("recorded", HourToTime(h))
			wg.Done()
		}
	}()
}

func (d *KeysDB) Lookup(key Hash) (res KeyMeta, err error) {
	err = d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keys"))
		v := b.Get([]byte(key[:]))
		if err := json.Unmarshal(v, &res); err != nil {
			return err
		}
		return nil
	})
	return
}

func (d *KeysDB) View(fn func(tx *bolt.Tx) error) error {
	return d.db.View(fn)
}

func (d *KeysDB) Close() error {
	return d.db.Close()
}
