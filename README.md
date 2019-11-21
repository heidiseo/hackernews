# TrueLayer Tech Challenge

This Command Line Interface that has a command `hackernews` and a flag `--posts` that consumes the number of posts user would like to see and returns that number of top posts.

## Built With:
* `go` version go1.13


## Libraries Used
* `github.com/spf13/cobra` for Command Line Interface. Easy to add Flags and modify the interface.
* `github.com/spf13/viper` for Command Line Interface Flags. Works well with the Cobra library above to modify Flag.
* `github.com/mitchellh/go-homedir` for Config, Env variables and finding home directory from the Command Line.

## Installing and Running

* Download Golang `https://golang.org/dl/`
* Setup Go Workspace `export PATH=$PATH:$GOPATH/bin` Add the bin directory from your workspace to your PATH(in your .profile file)
* To get libraries `go get {GitHub library path}`
* To create and add to go.mod file `go mod {GitHub library path}`
* To build the app `go build`
* To build and install the app to the bin directory `go install` (to run `hackernews` in command line)


## Functions(root.go)
    `rootCmd` variable
        * Root command
        * `posts` flag that takes a number of posts
        * checks if the number is higher than 100
        * gets information from other functions and returns in JSON format(string when printing to command line to be human readable)
    
    `GetTopStories` function
        * sends a GET request to Hackernews top story endpoint with the number received from the Root command
        * returns the IDs of the top stories

    `getScoredCards` function
        * sends a GET request to Hackernews item endpoint with the IDs received from the `GetTopStories`
        * returns the details of the posts
    
    `getScoredCards` function
        * validates the result from `getScoredCards`
        * converts the result to the right format


## Tests(root_test.go)

To run the tests do `go test`.