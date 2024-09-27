package main

import (
	"bytes"
	"errors"
	"slices"
	"strings"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/redcon"
)

func keys(db *badger.DB, conn redcon.Conn, cmd redcon.Command) (any, error) {
	if len(cmd.Args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'KEYS' command")
	}

	log.Debug().Str("pattern", string(cmd.Args[1])).Msg("keys")

	keys := []string{}
	if err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		if string(cmd.Args[1]) != "*" {
			opts.Prefix = []byte(cmd.Args[1])
		}

		if strings.Contains(string(cmd.Args[1]), "*") {
			opts.Prefix = []byte(strings.ReplaceAll(string(cmd.Args[1]), "*", ""))
		}

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid() || it.ValidForPrefix(opts.Prefix); it.Next() {
			if bytes.Contains(it.Item().Key(), []byte("set::")) {
				continue
			}
			k := string(it.Item().Key())
			if slices.Contains(keys, k) {
				continue
			}
			keys = append(keys, k)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}

	return keys, nil
}
