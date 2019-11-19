package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestFlagShort(t *testing.T) {
	var cArgs []string
	c := &cobra.Command{
		Use: "hackernews",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	var intFlagValue int
	var stringFlagValue string
	c.Flags().IntP("posts", "n", viper.GetInt("Posts"), "Set your name")
	c.Flags().GetInt("posts")
	c.Flags().StringVarP(&stringFlagValue, "sf", "s", "", "")
	output, err := executeCommand(c, "-i", "7", "-sabc", "one", "two")
	if output != "" {
		t.Errorf("Unexpected output: %v", err)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if intFlagValue != 7 {
		t.Errorf("Expected flag value: %v, got %v", 7, intFlagValue)
	}
	if stringFlagValue != "abc" {
		t.Errorf("Expected stringFlagValue: %q, got %q", "abc", stringFlagValue)
	}
	got := strings.Join(cArgs, " ")
	expected := "one two"
	if got != expected {
		t.Errorf("Expected arguments: %q, got %q", expected, got)
	}
}

func TestGetTopStories(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	url := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var arr []int
	err = json.Unmarshal(body, &arr)

	var ids []int
	for index, id := range arr {
		if index < 2 {
			ids = append(ids, id)
		}
	}
	actual, err := GetTopStories(2)
	if err != nil {
		g.Expect(err).To(gomega.HaveOccurred())
	} else {
		g.Expect(ids).To(gomega.Equal(actual))
	}
}

func TestGetIndividualStory(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ids := []int{21572552, 21574834, 21572622}
	rank := 1
	var allHackerNews []HackerNews
	for _, id := range ids {
		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", id)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		var result HackerNewsResponse
		err = json.Unmarshal(body, &result)
		hackerNews := HackerNews{
			Title:    result.Title,
			URI:      result.URI,
			Author:   result.Author,
			Points:   result.Points,
			Comments: len(result.Comments),
			Rank:     rank,
		}
		allHackerNews = append(allHackerNews, hackerNews)
		rank++
	}
	jsonHackerNews, _ := json.Marshal(allHackerNews)
	actual, err := GetIndividualStory(ids)
	if err != nil {
		g.Expect(err).To(gomega.HaveOccurred())
	} else {
		g.Expect(actual).To(gomega.Equal(string(jsonHackerNews)))
	}
}
