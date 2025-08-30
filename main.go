package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	githubCredentials := cli()
	//////////////////////////////////////:
	var wg sync.WaitGroup
	wg.Add(2)
	const followersEndpoint = "followers"
	const followingEndpoint = "following"
	var followers, following Result
	FileNameOfFollowers := "followers.csv"
	FileNameOfFollowing := "following.csv"
	FileNameOfWhoDoesntFollowBack := "WhoDoesntFollowBack.md"
	go getDatafromGithub(githubCredentials, followersEndpoint, &followers, &wg)
	go getDatafromGithub(githubCredentials, followingEndpoint, &following, &wg)
	wg.Wait()
	if following.err == nil && followers.err == nil {
		wg.Add(2)
		go createAndWriteToFile(FileNameOfFollowers, &followers, &wg)
		go createAndWriteToFile(FileNameOfFollowing, &following, &wg)
		wg.Wait()
		getWhoDosntFollowBack(FileNameOfWhoDoesntFollowBack, FileNameOfFollowers, FileNameOfFollowing)
	}
	/////////////////////////////////////////////////////////////////////
	fmt.Println("                                                                           ")
	fmt.Printf("Three files (csv format) were created!\n")
	fmt.Printf("%s : containing the list of users <YOU FOLLOW>  \n", FileNameOfFollowing)
	fmt.Printf("%s : containing the list of users <WHO FOLLOW YOU>  \n", FileNameOfFollowers)
	fmt.Printf("%s : containing the list of users <YOU FOLLOW BUT THEY DON'T FOLLOW YOU BACK>  \n", FileNameOfWhoDoesntFollowBack)
	fmt.Println("                                                                           ")
}

func createAndWriteToFile(nameOfFile string, dataFromGithub *Result, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf(">>>> begin writing to file %s ...\n", nameOfFile)
	file, erro := os.Create(nameOfFile)
	if erro != nil {
		log.Fatalf("ERROR: Could not create file %s because of error = %s\n", nameOfFile, erro.Error())
	}
	defer file.Close()
	str := "login,id,html_url\n"
	file.WriteString(str)
	for _, v := range dataFromGithub.data {
		str := fmt.Sprintf("%s,%s\n", v.Login, v.HtmlURL)
		file.WriteString(str)
	}
	fmt.Printf(">>>> finished writing to file %s ...\n", nameOfFile)
}

// IT GENERAT MARKDOWN FILE OF USERS WHO DONT FOLLOW BACK
func getWhoDosntFollowBack(nameOfFile string, FileNameOfFollowers string, FileNameOfFollowing string) {
	file, erro := os.Create(nameOfFile)
	if erro != nil {
		log.Fatalf("ERROR: Could not create file %s because of error = %s\n", nameOfFile, erro.Error())
	}
	defer file.Close()
	followersBytes, err1 := os.ReadFile(FileNameOfFollowers)
	if err1 != nil {
		log.Fatalf("ERROR: Could not create file %s because of error = %s\n", nameOfFile, err1.Error())
	}
	/////////////////////////////////:
	followingBytes, err2 := os.ReadFile(FileNameOfFollowing)
	if err2 != nil {
		log.Fatalf(
			"ERROR: Could not create file %s because of error = %s\n",
			nameOfFile,
			err2.Error(),
		)
	}
	newline := []byte("\n")
	linesFollowers := bytes.Split(followersBytes[1:len(followersBytes)-1], newline)
	linesFollowing := bytes.Split(followingBytes[1:len(followingBytes)-1], newline)
	comma := []byte(",")
	numberofCommas := bytes.Count(linesFollowers[0], comma)
	result := make(map[string]string)
	for _, line := range linesFollowers {
		vSplitInTwo := bytes.SplitN(line, comma, numberofCommas)
		key := strings.TrimSpace(string(vSplitInTwo[0]))
		value := strings.TrimSpace(string(vSplitInTwo[1]))
		_, ok := result[key]
		if !ok {
			result[key] = value
		}

	}
	for _, line := range linesFollowing {
		vSplitInTwo := bytes.SplitN(line, comma, numberofCommas)
		key := strings.TrimSpace(string(vSplitInTwo[0]))
		value := strings.TrimSpace(string(vSplitInTwo[1]))
		_, ok := result[key]
		if ok {
			delete(result, key)
		} else {
			result[key] = value
		}

	}
	fmt.Printf(">>>> begin writing to file %s ...\n", nameOfFile)

	Title := "# List of users <YOU FOLLOW BUT THEY DON'T FOLLOW YOU BACK>\n"
	file.WriteString(Title)
	for k, v := range result {
		tempStr := fmt.Sprintf("[%s](%s)\n", k, v)
		file.WriteString(tempStr)
	}
	fmt.Printf(">>>> finished writing to file %s \n", nameOfFile)
}

func cli() [2]string {
	fmt.Println("*****************************************************************************************")
	fmt.Println("********************   GITHUB FOLLOWER/UNFOLLOWER FETCHER   *****************************")
	fmt.Println("*****************************************************************************************")
	var inputs [2]string
	for {
		fmt.Printf(">>>> Please Enter your Github Username :    ")
		s := ""
		n, err1 := fmt.Scanln(&s)
		inputs[0] = strings.TrimSpace(s)

		s = ""
		fmt.Printf(">>>> Please Enter your Github Token :    ")
		m, err2 := fmt.Scanln(&s)
		inputs[1] = strings.TrimSpace(s)
		if err1 != nil || err2 != nil {
			if err1 != nil {
				log.Fatalf("ERROR : An error occured while processing the inputs err= %s\n", err1)
			}
			if err2 != nil {
				log.Fatalf("ERROR : An error occured while processing the inputs err= %s\n", err2)
			}
		}
		if n == 1 && m == 1 {
			break
		}
	}
	return inputs
}

type Result struct {
	data []Response
	err  error
}

//	curl -L \
//	 -H "Accept: application/vnd.github+json" \
//	 -H "Authorization: Bearer <YOUR-TOKEN>" \
//	 -H "X-GitHub-Api-Version: 2022-11-28" \
//	 https://api.github.com/repos/OWNER/REPO/issues
//
// /Assume theres only one page

func getDatafromGithub(githubCredential [2]string, endpoint string, result *Result, wg *sync.WaitGroup) {
	defer wg.Done()
	const GITHUBAPIURL = "https://api.github.com"
	// url := fmt.Sprint(GithubAPIURL, "/users/", githubCredential[0], "/", endpoint, "?per_page=100")
	url := fmt.Sprintf("%s/users/%s/%s", GITHUBAPIURL, githubCredential[0], endpoint)
	reqwst, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error occured while creating newRequest err : ", err)
		result.data = nil
		result.err = err
		return
	}
	reqwst.Header.Set("Accept", "application/vnd.github+json")

	authHeader := fmt.Sprintf("%s <%s>", githubCredential[0], githubCredential[1])
	reqwst.Header.Set("Authorization", authHeader)

	reqwst.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// use an http.Client to send the request:
	client := &http.Client{}
	resp, err := client.Do(reqwst)
	if err != nil {
		fmt.Println("Error occured while sending the request err : ", err)
		result.data = nil
		result.err = err
		return
	}
	defer resp.Body.Close()
	fmt.Println("                                                      ")
	fmt.Println(">>>>Fetching data from  ", resp.Request.URL)
	fmt.Println(">>>>resp.Status   ", resp.Status)
	fmt.Println("                                                      ")
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request.Body into the buffer err : ", err)
		result.data = nil
		result.err = err
		return
	}
	var response []Response // the array of struct that we decode to
	fmt.Println("                                                      ")
	fmt.Println("Begin Deserializing data ...")

	error := json.Unmarshal(respBody, &response)
	if error != nil {
		fmt.Println("Error occured while decoding json response err : ", err)
		result.data = nil
		result.err = err
		return
	}
	fmt.Println("Finished Deserializing data with no errors ...")
	fmt.Println("                                                      ")
	/////
	result.data = response
	result.err = nil
}

// "Respons ...
//
//	{
//	   "login": "octocat",
//	   "id": 1,
//	   "node_id": "MDQ6VXNlcjE=",
//	   "avatar_url": "https://github.com/images/error/octocat_happy.gif",
//	   "gravatar_id": "",
//	   "url": "https://HOSTNAME/users/octocat",
//	   "html_url": "https://github.com/octocat",
//	   "followers_url": "https://HOSTNAME/users/octocat/followers",
//	   "following_url": "https://HOSTNAME/users/octocat/following{/other_user}",
//	   "gists_url": "https://HOSTNAME/users/octocat/gists{/gist_id}",
//	   "starred_url": "https://HOSTNAME/users/octocat/starred{/owner}{/repo}",
//	   "subscriptions_url": "https://HOSTNAME/users/octocat/subscriptions",
//	   "organizations_url": "https://HOSTNAME/users/octocat/orgs",
//	   "repos_url": "https://HOSTNAME/users/octocat/repos",
//	   "events_url": "https://HOSTNAME/users/octocat/events{/privacy}",
//	   "received_events_url": "https://HOSTNAME/users/octocat/received_events",
//	   "type": "User",
//	   "site_admin": false
//	}""

type Response struct {
	Login string `json:"login"`
	// i am ignoring these fields
	// ID                uint64 `json:"id"`
	// NodeID            string `json:"node_id"`
	// AvatarURL         string `json:"avatar_url"`
	// GravatarID        string `json:"gravatar_url"`
	// URL               string `json:"url"`
	HtmlURL string `json:"html_url"`
	// FollowersURL      string `json:"followers_url"`
	// FollowingURL      string `json:"following_url"`
	// GistsURL          string `json:"gists_url"`
	// StarredURL        string `json:"starred_url"`
	// SubscriptionsURL  string `json:"subscriptions_url"`
	// OrganizationsURL  string `json:"organizations_url"`
	// ReposURL          string `json:"repos_url"`
	// EventsURL         string `json:"events_url"`
	// ReceivedEventsURL string `json:"received_events_url"`
	// Type              string `json:"type"`
	// SiteAdmin         bool   `json:"site_admin"`
}
