package ui

import (
	"os"

	"github.com/rrenode/gonodedb/model"
	"github.com/rrenode/gonodedb/store"

	"github.com/dgraph-io/badger/v4"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func PrintRepoArray(repos []model.Repo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"Name", "Alias", "ID", "Path"})

	for _, r := range repos {
		t.AppendRow(table.Row{
			color.New(color.FgCyan).Sprint(r.Name),
			color.New(color.FgYellow).Sprint(r.Alias),
			color.New(color.FgGreen).Sprint(r.ID),
			color.New(color.FgMagenta).Sprint(r.Path),
		})
	}

	t.Render()
}

func PrintDB(db *badger.DB) error {
	repos, err := store.LoadAllRepos(db)
	if err != nil {
		return err
	}
	PrintRepoArray(repos)
	return nil
}
