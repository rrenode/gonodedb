package store

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/rrenode/gonodedb/model"

	"github.com/dgraph-io/badger/v4"
)

func autoIndexRepo(txn *badger.Txn, repo *model.Repo) error {
	indexes := map[string]string{
		"name":  repo.Name,
		"alias": repo.Alias,
		"path":  repo.Path,
	}
	for field, value := range indexes {
		key := fmt.Sprintf("%s:%s", field, value)
		if err := txn.Set([]byte(key), []byte(repo.ID)); err != nil {
			return err
		}
	}
	return nil
}

func SaveRepo(db *badger.DB, repo *model.Repo) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(repo)
		if err != nil {
			return err
		}
		if err := txn.Set([]byte("repo:"+repo.ID), data); err != nil {
			return err
		}
		return autoIndexRepo(txn, repo)
	})
}

func SaveRepoIfChanged(db *badger.DB, repo *model.Repo) error {
	existing, err := LoadRepoByID(db, repo.ID)
	if err == nil && reflect.DeepEqual(existing, repo) {
		return nil
	}
	return SaveRepo(db, repo)
}

func LoadRepoByID(db *badger.DB, id string) (*model.Repo, error) {
	var repo model.Repo
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("repo:" + id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &repo)
		})
	})
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func LoadRepoByField(db *badger.DB, field, value string) (*model.Repo, error) {
	var repoID string
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(field + ":" + value))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			repoID = string(val)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return LoadRepoByID(db, repoID)
}

func LoadAllRepos(db *badger.DB) ([]model.Repo, error) {
	var repos []model.Repo
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("repo:")

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var r model.Repo
				if err := json.Unmarshal(v, &r); err != nil {
					return err
				}
				repos = append(repos, r)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return repos, err
}
