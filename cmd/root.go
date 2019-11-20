package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

//HackerNewsResponse is the response of API endpoint for specific ID if successful
type HackerNewsResponse struct {
	Author   string `json:"by"`
	ID       int    `json:"id"`
	Points   int    `json:"score"`
	Title    string `json:"title"`
	URI      string `json:"url"`
	Comments []int  `json:"kids"`
}

//HackerNews is the response of the command
type HackerNews struct {
	Title    string `json:"title"`
	URI      string `json:"uri"`
	Author   string `json:"author"`
	Points   int    `json:"points"`
	Comments int    `json:"comments"`
	Rank     int    `json:"rank"`
}

// rootCmd represents the base command - hackernews
var rootCmd = &cobra.Command{
	Use:   "hackernews",
	Short: "posts from hackernews",
	Long:  `will return the top number of posts the user requests`,

	Run: func(cmd *cobra.Command, args []string) {
		//adds a flag called posts and gets the number entered after the flag
		number, _ := cmd.Flags().GetInt("posts")
		//calls API with the number received above and gets the number of top stories
		ids, err := GetTopStories(number)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		//calls API with the top story IDs to call another end point to receive full info
		responses, err := GetIndividualStory(ids)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		//validates and formats the response
		formatedResponses, err := ResponseFormat(responses)
		//prints on the command line
		fmt.Println(formatedResponses)
	},
}

// Execute the root command and sets the flag.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//defines flags
func init() {
	cobra.OnInitialize(initConfig)

	// Cobra supports persistent flags, which is global for application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hello-cobra.yaml)")

	// to run when this action is called directly
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().IntP("posts", "n", viper.GetInt("Posts"), "enter how many posts you would like to see")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hello-cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hello-cobra")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

//GetTopStories gets a numbner of posts from the command line and calls Hackernews /topstories endpoint to get IDs
func GetTopStories(i int) ([]int, error) {
	url := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var arr []int
	err = json.Unmarshal(body, &arr)
	if err != nil {
		return nil, err
	}

	var ids []int
	for index, id := range arr {
		if index < i {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

//GetIndividualStory gets top story IDs and calls item end poit to get detailed information
func GetIndividualStory(ids []int) ([]HackerNewsResponse, error) {
	var responses []HackerNewsResponse
	for _, id := range ids {
		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", id)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var result HackerNewsResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			return nil, err
		}
		responses = append(responses, result)
	}
	return responses, nil
}

//ResponseFormat gets posts details, validates and covert it to the response format for the command line.
func ResponseFormat(responses []HackerNewsResponse) (string, error) {
	rank := 1
	var allHackerNews []HackerNews
	for _, response := range responses {
		if response.Title == "" || response.Author == "" {
			return "", fmt.Errorf("title and/or author are empty")
		}

		if len(response.Title) > 256 || len(response.Author) > 256 {
			return "", fmt.Errorf("title and/or author have more than 256 characters")
		}

		if response.Points < 0 || len(response.Comments) < 0 || rank < 0 {
			return "", fmt.Errorf("points/comments/rank have negative number(s)")
		}

		u, err := url.Parse(response.URI)
		if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
			return "", fmt.Errorf("invalid URI")
		}

		hackerNews := HackerNews{
			Title:    response.Title,
			URI:      response.URI,
			Author:   response.Author,
			Points:   response.Points,
			Comments: len(response.Comments),
			Rank:     rank,
		}
		allHackerNews = append(allHackerNews, hackerNews)
		//adds the rank by the post order
		rank++
	}
	//formats it to JSON
	jsonHackerNews, err := json.Marshal(allHackerNews)
	if err != nil {
		return "", err
	}
	//converts it to string so that it is human readable in command line
	return string(jsonHackerNews), nil

}
