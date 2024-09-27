package main

import (
	"errors"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/redcon"
)

func exists(db *badger.DB, _ redcon.Conn, cmd redcon.Command) (any, error) {
	if len(cmd.Args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'EXISTS' command")
	}
	log.Debug().Str("key", string(cmd.Args[1])).Msg("exists")

	exists := 0
	if err := db.View(func(txn *badger.Txn) error {
		item, gErr := txn.Get(cmd.Args[1])
		if gErr != nil {
			if errors.Is(gErr, badger.ErrKeyNotFound) {
				return nil
			}
			return gErr
		}

		if err := item.Value(func(val []byte) error {
			if len(val) != 0 {
				exists = 1
			}
			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return exists, nil
}
