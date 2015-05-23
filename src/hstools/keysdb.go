package hstools

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type KeyMeta struct {
	FirstSeen Hour
	LastSeen  Hour
}

type KeysDB struct {
	db *bolt.DB
}

func OpenKeysDb(filename string) (*KeysDB, error) {
	// os.Remove(filename)

	d := &KeysDB{}
	var err error

	if d.db, err = bolt.Open(filename, 0664, nil); err != nil {
		return nil, err
	}

	if err = d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Keys"))
		return err
	}); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *KeysDB) Seen(keys []Hash, h Hour) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		for _, k := range keys {
			b := tx.Bucket([]byte("Keys"))
			oldJSON := b.Get(k[:])
			meta := KeyMeta{
				FirstSeen: h,
				LastSeen:  h,
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
			}
			encoded, err := json.Marshal(meta)
			if err != nil {
				return err
			}
			if err := b.Put(k[:], encoded); err != nil {
				return err
			}
		}
		return nil
	})

}

func (d *KeysDB) Close() error {
	return d.db.Close()
}
