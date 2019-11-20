package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/onsi/gomega"
)

//TestGetTopStories tests if the function returns the right top story IDs
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

//TestGetIndividualStory tests if the function returns the right top story details
func TestGetIndividualStory(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var responses []HackerNewsResponse
	ids := []int{21572552, 21574834, 21572622}
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
		responses = append(responses, result)

	}
	actual, err := GetIndividualStory(ids)
	if err != nil {
		g.Expect(err).To(gomega.HaveOccurred())
	} else {
		g.Expect(actual).To(gomega.Equal(responses))
	}
}

//TestResponseFormat tests if posts pass the requirements and convers it to the right format to return
func TestResponseFormat(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	tests := []struct {
		Message  string
		Response []HackerNewsResponse
		Error    string
	}{
		{Message: "should fail as author is emplty",
			Response: []HackerNewsResponse{
				{Author: "",
					Title: "TrueLayer",
				},
			},
			Error: "title and/or author are empty",
		},
		{Message: "should fail as title has more than 256 characters",
			Response: []HackerNewsResponse{
				{Author: "TrueLayer",
					Title: "incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
				},
			},
			Error: "title and/or author have more than 256 characters",
		},
		{Message: "should fail as URL is not valid",
			Response: []HackerNewsResponse{
				{Author: "TrueLayer",
					Title: "Data API",
					URI:   "jskd://invalid.url",
				},
			},
			Error: "invalid URI",
		},
		{Message: "should fail as points is negative",
			Response: []HackerNewsResponse{
				{Author: "TrueLayer",
					Points: -2,
					Title:  "Data API",
				},
			},
			Error: "points/comments/rank have negative number(s)",
		},
		{Message: "should not fail as all requirements are met",
			Response: []HackerNewsResponse{
				{Author: "TrueLayer",
					ID:       34234,
					Points:   12,
					Title:    "Data API",
					URI:      "https://truelayer.com/",
					Comments: []int{1, 2, 3, 4, 5},
				},
			},
			Error: "",
		},
	}

	expected := `[{"title":"Data API","uri":"https://truelayer.com/","author":"TrueLayer","points":12,"comments":5,"rank":1}]`

	for _, test := range tests {
		t.Run(test.Message, func(t *testing.T) {
			response, err := ResponseFormat(test.Response)
			if err != nil {
				g.Expect(err).To(gomega.HaveOccurred())
				g.Expect(err.Error()).To(gomega.Equal(test.Error))
			} else {
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(response).To(gomega.Equal(expected))
			}
		})
	}
}
