package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func ratelimit() *http.Client {
	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
	}
	return rateLimiter
}

// Fetch all the public organizations' membership of a user.
func fetchOrganizations(username string, client *github.Client) ([]*github.Organization, error) {
	//client := github.NewClient(ratelimit())
	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	return orgs, err
}

func fetchOrganizationAuditLogs(ctx context.Context, organization string, client *github.Client) []map[string]interface{} {
	//client := github.NewClient(nil)
	req, err := client.NewRequest(http.MethodGet, fmt.Sprintf("orgs/%s/audit-log?", organization), nil)
	if err != nil {
		log.Fatal(err)
	}

	var entries []map[string]interface{}

	res, err := client.Do(ctx, req, &entries)
	if err != nil {
		log.Fatal("get audit log: %w", err)
	}
	log.Println(res)
	return entries
}

func main() {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "......."},
	)
	tc := oauth2.NewClient(ctx, ts)
	// list all repositories for the authenticated user
	//repos, _, err := client.Repositories.List(ctx, "", nil)
	client := github.NewClient(tc)

	//fetchOrganizations("willnorris", client)

	r, resp, err := client.RateLimits(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fetchOrganizationAuditLogs(ctx, "BusinessName", client))

	log.Println(r.Core.Remaining)
	log.Println(resp.Rate.Remaining)

	log.Println(r.Search.Limit)
	log.Println(r.Search.Remaining)

}
