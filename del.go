package main

import (
	"bytes"
	"errors"
	"strings"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/redcon"
)

func del(db *badger.DB, conn redcon.Conn, cmd redcon.Command) (any, error) {
	if len(cmd.Args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'DEL' command")
	}

	log.Debug().Str("key", string(cmd.Args[1])).Msg("del")

	mu.Lock()
	defer mu.Unlock()

	amountDeleted := 0
	if err := db.Update(func(txn *badger.Txn) error {
		if !strings.Contains(string(cmd.Args[1]), "*") {
			i, err := txn.Get(cmd.Args[1])
			if err != nil {
				return err
			}

			if err := i.Value(func(val []byte) error {
				amountDeleted++
				return txn.Delete(cmd.Args[1])
			}); err != nil {
				return err
			}

			return nil
		}

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		if string(cmd.Args[1]) != "*" && strings.Contains(string(cmd.Args[1]), "*") {
			opts.Prefix = []byte(strings.ReplaceAll(string(cmd.Args[1]), "*", ""))
		}

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid() || it.ValidForPrefix(opts.Prefix); it.Next() {
			if bytes.Contains(it.Item().Key(), []byte("set::")) {
				continue
			}
			k := it.Item().KeyCopy(nil)
			if err := txn.Delete(k); err != nil {
				return err
			}

			amountDeleted++
		}

		return nil
	}); err != nil {
		return redcon.SimpleInt(0), nil
	}

	return amountDeleted, nil
}
