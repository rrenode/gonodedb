// Example/Example.go
package Example

import (
	"fmt"
	"log"

	"github.com/rrenode/gonodedb/db"
	"github.com/rrenode/gonodedb/model"
	"github.com/rrenode/gonodedb/store"
	"github.com/rrenode/gonodedb/ui"
)

func RunExample() {
	badgerDB, err := db.Open("./badger.db", false)
	if err != nil {
		log.Fatal(err)
	}
	defer badgerDB.Close()
	defer db.RunGC(badgerDB, 0.5)

	repos := []model.Repo{
		{"43624344-5508-4065-9cbf-ba27d33c472d", "pph-python", "pph", "C:\\Projects\\pph"},
		{"11b42378-9ec5-4f9a-9e37-4df0ecb0af95", "alpha", "a", "/code/alpha"},
		{"b0b0defd-594c-4e5c-8545-a0890207f4e2", "beta", "b", "/code/beta"},
	}

	for _, r := range repos {
		if err := store.SaveRepoIfChanged(badgerDB, &r); err != nil {
			log.Println("Save failed:", r.Name)
		}
	}

	all, _ := store.LoadAllRepos(badgerDB)
	ui.PrintRepoArray(all)

	r, _ := store.LoadRepoByField(badgerDB, "alias", "a")
	fmt.Println("Queried:", r.Name)

	matches, err := store.FuzzySearchReposAcrossFields(badgerDB, "path", "od", 60)
	if err != nil {
		log.Fatal(err)
	}
	var matchedRepos []model.Repo
	for _, m := range matches {
		matchedRepos = append(matchedRepos, m.Repo)
	}

	fmt.Println("Fuzzy matches for 'o':")
	ui.PrintRepoArray(matchedRepos)
}
