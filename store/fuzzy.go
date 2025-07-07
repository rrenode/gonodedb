package store

import (
	"sort"
	"strings"

	"github.com/rrenode/gonodedb/model"

	"github.com/dgraph-io/badger/v4"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type MatchResult struct {
	Repo      model.Repo
	MatchText string
	Score     int
}

func FuzzySearchReposAcrossFields(db *badger.DB, fields, term string, minScore int) ([]MatchResult, error) {
	all, err := LoadAllRepos(db)
	if err != nil {
		return nil, err
	}

	var results []MatchResult
	fieldList := strings.Split(fields, "|")
	term = strings.ToLower(term)

	for _, repo := range all {
		var bestMatch MatchResult
		bestMatch.Score = -1

		for _, field := range fieldList {
			var value string
			switch field {
			case "name":
				value = repo.Name
			case "alias":
				value = repo.Alias
			case "path":
				value = repo.Path
			default:
				continue
			}

			valueLower := strings.ToLower(value)

			if strings.Contains(valueLower, term) {
				// Direct substring match — strong match
				bestMatch = MatchResult{
					Repo:      repo,
					MatchText: value,
					Score:     100,
				}
				break // skip other fields if one matches directly
			}

			score := fuzzy.LevenshteinDistance(term, valueLower)
			maxLen := max(len(term), len(valueLower))
			if maxLen == 0 {
				continue
			}
			similarity := (maxLen - score) * 100 / maxLen
			if similarity >= minScore && similarity > bestMatch.Score {
				bestMatch = MatchResult{
					Repo:      repo,
					MatchText: value,
					Score:     similarity,
				}
			}

		}

		if bestMatch.Score >= minScore {
			results = append(results, bestMatch)
		}
	}

	// Sort descending by Score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
