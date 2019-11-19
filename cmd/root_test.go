package cmd

import (
	"fmt"
	"testing"
)

// func TestFlagShort(t *testing.T) {
// 	var cArgs []string
// 	c := &cobra.Command{
// 		Use: "hackernews",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			numberString, _ := cmd.Flags().GetString("posts")
// 			number, err := strconv.Atoi(numberString)
// 			ids, err := GetTopStories(number)
// 			if err != nil {
// 				fmt.Println(err)
// 				log.Fatal(err)
// 			}
// 			allStories, err := GetIndividualStory(ids)
// 			if err != nil {
// 				fmt.Println(err)
// 				log.Fatal(err)
// 			}
// 			fmt.Println(allStories)
// 		},
// 	}
// 	var intFlagValue int
// 	var stringFlagValue string
// 	c.Flags().IntVarP(&intFlagValue, "intf", "i", -1, "")
// 	c.Flags().StringVarP(&stringFlagValue, "sf", "s", "", "")
// 	output, err := executeCommand(c, "-i", "7", "-sabc", "one", "two")
// 	if output != "" {
// 		t.Errorf("Unexpected output: %v", err)
// 	}
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	if intFlagValue != 7 {
// 		t.Errorf("Expected flag value: %v, got %v", 7, intFlagValue)
// 	}
// 	if stringFlagValue != "abc" {
// 		t.Errorf("Expected stringFlagValue: %q, got %q", "abc", stringFlagValue)
// 	}
// 	got := strings.Join(cArgs, " ")
// 	expected := "one two"
// 	if got != expected {
// 		t.Errorf("Expected arguments: %q, got %q", expected, got)
// 	}
// }

func TestGetTopStories(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	ids := []int{21572552, 21572795}

	number, err := GetTopStories(2)
	if err != nil {
		g.Expect(err).To(gomega.HaveOccurred())
		fmt.Println(err)
		//g.Expect(err).To(gomega.Equal())
	} else {
		g.Expect(number).To(gomega.Equal(ids))
	}
}
