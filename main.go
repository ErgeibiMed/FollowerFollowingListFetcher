package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	githubCredentials := cli()
	const followersEndpoint = "followers"
	var followers Result
	getDatafromGithub(githubCredentials, followersEndpoint, &followers)
	filefollowers := "followers.txt"
	if followers.err == nil {
		fmt.Printf("writing to file %s ...\n", filefollowers)
		file, err := os.Create(filefollowers)
		if err != nil {
			fmt.Printf("failed to create file %s because err: %s\n", filefollowers, err.Error())
		}
		defer file.Close()
		for _, v := range followers.data {
			str := fmt.Sprintf("userlogin = %s  |  id = %8d  url = %s \n", v.Login, v.ID, v.HTMLURL)
			file.WriteString(str)
			delimeter := []byte("------------------------------------------------------------------------\n")
			file.Write(delimeter)
		}
	}
	const followingEndpoint = "following"
	var following Result
	getDatafromGithub(githubCredentials, followingEndpoint, &following)
	filefollowing := "following.txt"
	if following.err == nil {
		fmt.Printf("writing to file %s ...\n", filefollowing)
		file, err := os.Create(filefollowing)
		if err != nil {
			fmt.Printf("failed to create file %s because err: %s\n", filefollowing, err.Error())
		}
		defer file.Close()
		for _, v := range following.data {
			str := fmt.Sprintf("userlogin = %s  |  id = %8d  url = %s \n", v.Login, v.ID, v.HTMLURL)
			file.WriteString(str)
			delimeter := []byte("------------------------------------------------------------------------\n")
			file.Write(delimeter)
		}
	}
	fmt.Println("===========================================================================")
	fmt.Printf("Two files were created : %s and %s \n", filefollowers, filefollowing)
	fmt.Println("===========================================================================")
}

/// // ////////followButTheydontfollowBck
// TODO: you follow But They dont follow back in one file or output to terminal
/// func followButTheydontfollowBck(input1 Result, input2 Result) [][]string {
/// 	len1 := len(input1.data)
/// 	len2 := len(input2.data)
/// 	uniqueUSername := make(map[string][]string)
/// 	for i := range len1 {
/// 		v1 := input1.data[i].Login
/// 		htmlURL1 := input1.data[i].HtmlURL
/// 		id1 := input1.data[i].ID
/// 		id1AsStr := strconv.Itoa(int(id1))
/// 		uniqueUSername[v1] = append(uniqueUSername[v1], v1)
/// 		uniqueUSername[v1] = append(uniqueUSername[v1], id1AsStr)
/// 		uniqueUSername[v1] = append(uniqueUSername[v1], htmlURL1)
/// 	}
/// 	for i := range len2 {
/// 		v2 := input2.data[i].Login
/// 		htmlURL2 := input2.data[i].HtmlURL
/// 		id2 := input2.data[i].ID
/// 		id2AsStr := strconv.Itoa(int(id2))
/// 		uniqueUSername[v2] = append(uniqueUSername[v2], v2)
/// 		uniqueUSername[v2] = append(uniqueUSername[v2], id2AsStr)
/// 		uniqueUSername[v2] = append(uniqueUSername[v2], htmlURL2)
/// 	}
/// 	result := make([][]string, 30)
/// 	for key := range uniqueUSername {
/// 		value, ok := uniqueUSername[key]
/// 		if ok {
/// 			result = append(result, value)
/// 		}
/// 	}
/// 	return result
/// }

func cli() [2]string {
	fmt.Println("-----------------------Github follower/Unfollower tool--------------------------------")
	var inputs [2]string
	for {
		fmt.Printf("Please Enter your Github Username : ")
		s := ""
		n, err1 := fmt.Scanln(&s)
		inputs[0] = strings.TrimSpace(s)

		s = ""
		fmt.Printf("Please Enter your Github Token : ")
		m, err2 := fmt.Scanln(&s)
		inputs[1] = strings.TrimSpace(s)
		if err1 != nil || err2 != nil {
			log.Fatal(err1, err2)
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

func getDatafromGithub(githubCredential [2]string, endpoint string, result *Result) {
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
	fmt.Println("---------------------------")
	fmt.Println(" resp.Status=", resp.Status)
	fmt.Println("---------------------------")
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error occured while reading the request.Body into the buffer err : ", err)
		result.data = nil
		result.err = err
		return
	}
	//	fmt.Println("-----respBody---------")
	//	fmt.Println(string(respBody))
	//	fmt.Println("-----respBody---------")
	var response []Response // the array of struct that we decode to
	error := json.Unmarshal(respBody, &response)
	if error != nil {
		fmt.Println("Error occured while decoding json response err : ", err)
		result.data = nil
		result.err = err
		return
	}
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
	ID    uint64 `json:"id"`
	//////
	//iam ignoring these fields
	/////
	//nodeID            string
	//avatarURL         string
	//gravataID         string
	//URL string `json:"url"`
	HTMLURL string `json:"html_url"`
	// followersURL      string
	// followingURL      string
	// gistsURL          string
	// starredURL        string
	// subscriptionsURL  string
	// organizationsURL  string
	// reposURL          string
	// eventsURL         string
	// receivedEventsURL string
	// typeUser          string
	// siteAdmin         bool
}
