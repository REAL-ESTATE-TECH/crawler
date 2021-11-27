package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const Depth = 4

func TestCrawl(t *testing.T) {
	t.Run("Fetch all annual reports on https://www.brfvärnhem1.se", func(t *testing.T) {
		want := []string{
			"https://drive.google.com/file/d/1esZlXo18xauddHwBPYQEqGLx9o9jkM38/view?usp=sharing",
			"https://drive.google.com/file/d/1JSOVunUt02CLMFu9Y4qWs4aN-PBu_j3Z/view?usp=sharing",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2018.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2017.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2016.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2015.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2014.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2013.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2012.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2011.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2008.pdf",
			"https://xn--brfvrnhem1-t5a.se/themes/BrfVarnhem/assets/files/Årsredovisning_BRF_Värnhem1_2007.pdf",
		}

		crawler := New()
		res := crawler.Crawl("https://xn--brfvrnhem1-t5a.se/", 2, 8)
		t.Log(res)
		assert.True(t, subset(want, res))
	})
}

// helper

// subset returns true if the first array is completely
// contained in the second array.
func subset(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value] += 1
	}
	for _, value := range first {
		if _, ok := set[value]; !ok {
			return false
		}
	}
	return true
}
