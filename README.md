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
*****************************************************************************************
********************   GITHUB FOLLOWER/UNFOLLOWER FETCHER   *****************************
*****************************************************************************************
>>>> Please Enter your Github Username :    JohnDoe
```
Enter your github username and press Enter

Second you will be prompted to enter AuthToken

Token Authentication: The program requests a GitHub personal access token (with the user:read scope) to authenticate and access user data securely.

**HOWEVER IT DOES NOT STORE YOUR GITHUB AUTH TOKEN**

```console
>>>> Please Enter your Github Token :   Your_Github_Token 
"
```
Enter your github auth token and press Enter

## Data Retrieval – Using the GitHub API, it fetches:

A list of users who follow the authenticated user.

A list of users the authenticated user is following.

## Output

Three files (csv format) were created!

1. following.csv : containing the list of users **YOU FOLLOW** 

2. followers.csv : containing the list of users **WHO FOLLOW YOU**

3. WhoDoesntFollowBack.md : is a Markdown file containing the list of users **YOU FOLLOW BUT THEY DON'T FOLLOW YOU BACK**


