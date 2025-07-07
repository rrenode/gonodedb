package db

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

func RunGC(db *badger.DB, threshold float64) {
	for {
		err := db.RunValueLogGC(threshold)
		if err == badger.ErrNoRewrite {
			fmt.Println("No GC needed")
			break
		} else if err != nil {
			fmt.Println("GC error:", err)
			break
		}
		fmt.Println("GC ran")
	}
}
