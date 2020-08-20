package hasura

import (
	"os"

	"github.com/shurcooL/graphql"
)

// Client returns a GraphQL client for Hasura.
func Client() *graphql.Client {
	graphqlURL := "http://localhost:8080/hasura/v1/graphql"
	if val, ok := os.LookupEnv("graphql_url"); ok {
		graphqlURL = val
	}
	client := graphql.NewClient(graphqlURL, nil)
	return client
}
