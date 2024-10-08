package main

import (
	"errors"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/redcon"
)

func get(db *badger.DB, _ redcon.Conn, cmd redcon.Command) (any, error) {
	if len(cmd.Args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'GET' command")
	}

	log.Debug().Str("key", string(cmd.Args[1])).Msg("get")

	var bdgrVal []byte
	if err := db.View(func(txn *badger.Txn) error {
		item, gErr := txn.Get(cmd.Args[1])
		if gErr != nil {
			return gErr
		}

		if err := item.Value(func(val []byte) error {
			bdgrVal = append([]byte{}, val...)
			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, nil
	}

	return bdgrVal, nil
}
