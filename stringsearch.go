package main

import "strings"

func naiveSearch(searches *[]string, domain string) []string {
	var matches []string
	for _, s := range *searches {
		if strings.Contains(domain, s) {
			matches = append(matches, s)
		}
	}
	return stringSliceUniq(matches)
}

func trieSearch(runeTrie *RuneTrie, domain string) []string {
	var matches []string
	domainRunes := []rune(domain)

	for searchPos := 0; searchPos < len(domainRunes); searchPos++ {
		var searchTrie = runeTrie
		var matchedRunes []rune
		for endPos := searchPos; endPos < len(domainRunes); endPos++ {
			thisRune := domainRunes[endPos]
			var thisChild = searchTrie.children[thisRune]
			if thisChild == nil {
				break
			}
			matchedRunes = append(matchedRunes, thisRune)
			if thisChild.value == true {
				matches = append(matches, string(matchedRunes))
			}
			searchTrie = thisChild
		}
	}

	return matches
}


