package main

import (
	"strings"
	"testing"
)

var domains = []string{
	"www-go0gledrive.co",
	"dzombak.com",
	"goodll.co",
	"www.www-micro-soft.com",
	"hjkhgrdhghghghmcimicrospfccm8978787097890yiuouihhjklfgjhklfgjhkfgkjhlfgdjhklfgdhjklfgdjm.xxx",
	"groveid.com",
	"groove.id",
	"cob.archive.nrtfa.fa.namdmz.dmzroot.net",
	"www-grooveid.online.tk",
	"testing.com",
	"example.com",
	"mising.info",
}

func buildSearches() []string {
	searchSources := []string{
		"google",
		"google.com",
		"microsoft",
		"microsoft.com",
		"grooveid",
		"grooveid.com",
		"groove.id",
		"missing.info",
	}
	var results = searchSources
	for _, s := range searchSources {
		results = append(results, additionAttack(s)...)
		results = append(results, vowelswapAttack(s)...)
		results = append(results, transpositionAttack(s)...)
		results = append(results, replacementAttack(s)...)
		results = append(results, repetitionAttack(s)...)
		results = append(results, omissionAttack(s)...)
		results = append(results, bitsquattingAttack(s)...)
		results = append(results, homographAttack(s)...)
		results = append(results, subdomainAttack(s)...)
		results = append(results, hyphenationAttack(s)...)
	}
	return results
}

func BenchmarkNaive(b *testing.B) {
	searches := buildSearches()
	println("begin naiveSearch w/", len(searches), "searches")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, d := range domains {
			result := naiveSearch(&searches, d)
			b.StopTimer()
			if n == 0 && len(result) > 0 {
				println(d, "->", strings.Join(result, ","))
			}
			b.StartTimer()
		}
	}
}

func BenchmarkTrie(b *testing.B) {
	searches := buildSearches()
	println("begin trieSearch w/", len(searches), "searches")

	searchTrie := NewRuneTrie()
	for _, s := range searches {
		searchTrie.Put(s, true)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, d := range domains {
			result := trieSearch(searchTrie, d)
			b.StopTimer()
			if n == 0 && len(result) > 0 {
				println(d, "->", strings.Join(result, ","))
			}
			b.StartTimer()
		}
	}
}

func TestUsingTrieSearchDoesNotAlterTrie(t *testing.T) {
	searches := buildSearches()
	searchTrie := NewRuneTrie()
	for _, s := range searches {
		searchTrie.Put(s, true)
	}
	searchTrie2 := NewRuneTrie()
	for _, s := range searches {
		searchTrie2.Put(s, true)
	}

	_ = trieSearch(searchTrie, "microsoft.com")

	if len(searchTrie.children) != len(searchTrie2.children) {
		t.Error("searching the trieSearch altered len(children)")
	}
	for k := range searchTrie.children {
		if searchTrie2.children[k] == nil {
			t.Errorf("key %v exists in trieSearch used for search, but not a newly created trieSearch", k)
		}
	}
}

func TestTrieSearch(t *testing.T) {
	searches := buildSearches()
	searchTrie := NewRuneTrie()
	for _, s := range searches {
		searchTrie.Put(s, true)
	}

	expected := map[string][]string {
		"www-go0gledrive.co": {"go0gle"},
		"www.www-micro-soft.com": {"micro-soft", "micro-soft.com"},
		"groveid.com": {"groveid", "groveid.com"},
		"groove.id": {"groove.id","roove.id","groove.i"},
		"www-grooveid.online.tk": {"grooveid","rooveid","groovei"},
		"dzombak.com": {},
		"mising.info": {"mising.info"},
	}

	for domain, expectedResult := range expected {
		result := trieSearch(searchTrie, domain)

		if len(expectedResult) != len(result) {
			t.Errorf("%s: expected %d results; got %d", domain, len(expectedResult), len(result))
		}
		for _, s := range expectedResult {
			if !stringSliceContains(result, s) {
				t.Errorf("%s: expected result %s is missing", domain, s)
			}
		}
	}
}

func TestTrieSearchMatchesNaiveSearch(t *testing.T) {
	searches := buildSearches()
	searchTrie := NewRuneTrie()
	for _, s := range searches {
		searchTrie.Put(s, true)
	}

	for _, d := range domains {
		trieResult := trieSearch(searchTrie, d)
		naiveResult := naiveSearch(&searches, d)

		if len(trieResult) != len(naiveResult) {
			t.Errorf("%s: trie got %d results; naive got got %d", d, len(trieResult), len(naiveResult))
		}
		for _, s := range naiveResult {
			if !stringSliceContains(trieResult, s) {
				t.Errorf("%s: naive result %s is missing from trie result", d, s)
			}
		}
	}
}
