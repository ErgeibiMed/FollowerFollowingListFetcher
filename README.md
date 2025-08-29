# GitHub Follower/Following Fetcher 

## Description

This program allows users to fetch the list of account that follow them
and the list of their followers.

First clone this repo

```console
  git clone https://github.com/ErgeibiMed/FollowerFollowingListFetcher.git
  cd FollowerFollowingListFetcher 
  go run main.go
```
Upon execution 

First you will be prompted to enter your github USername

```console
-----------------------Github follower/Unfollower tool--------------------------------
Please Enter your Github Username : JohnDoe
```
Enter your github username and press Enter

Second you will be prompted to enter AuthToken

Token Authentication: The program requests a GitHub personal access token (with the user:read scope) to authenticate and access user data securely.

HOWEVER IT DOES NOT STORE YOUR GITHUB AUTH TOKEN

```console
Please Enter your Github Token : Yourgithub_token
```
Enter your github auth token and press Enter

## Data Retrieval – Using the GitHub API, it fetches:

A list of users who follow the authenticated user.

A list of users the authenticated user is following.

## Output Generation 

Two text files are generated in the working directory:

followers.txt –> Contains the usernames of the followers.

following.txt –> Contains the usernames of the accounts the user follows.

## Disclaimer
This tool is useful to get rid of those annoying account looking for (literally) clout thinking this is instagram.

