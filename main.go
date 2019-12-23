package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// dummy db data
var casefiles = []*casefile{
	{
		ID:          "1000",
		Number:      "1",
		InitiatedAt: "2019-01-02T15:04:05.000+0:00",
		ClosedAt:    "2019-02-02T15:04:05.000+0:00",
		Status:      "Denied",
		FlagCount:   "42",
	},
	{
		ID:          "2000",
		Number:      "2",
		InitiatedAt: "2019-01-02T15:04:05.000+0:00",
		Status:      "Approved",
		FlagCount:   "1",
	},
	{
		ID:          "3000",
		Number:      "3",
		InitiatedAt: "2019-01-02T15:04:05.000+0:00",
		ClosedAt:    "2019-03-02T15:04:05.000+0:00",
		Status:      "Denied",
		FlagCount:   "907",
	},
	{
		ID:          "4000",
		Number:      "4",
		InitiatedAt: "2019-01-02T15:04:05.000+0:00",
		FlagCount:   "0",
	},
}

type casefile struct {
	ID          graphql.ID // scalar by default
	Number      string
	InitiatedAt string
	ClosedAt    string
	Status      string
	FlagCount   string
}

// fake db map of dummy data
var casefileData = make(map[graphql.ID]*casefile)

type Resolver struct{}

// stores our casefile info query
type caseResolver struct {
	c *casefile
}

// Casefile method extracts info from the caseResolver struct to execute the query on a case by ID
func (r *Resolver) Casefile(args struct{ ID graphql.ID }) *caseResolver {
	if c := casefileData[args.ID]; c != nil {
		return &caseResolver{c}
	}
	return nil
}

// a reflection method has to resolve each field
func (r *caseResolver) ID() graphql.ID {
	return graphql.ID(r.c.ID) // coerce the string to the graphql.ID type
}

func (r *caseResolver) Number() string {
	return r.c.Number
}

func (r *caseResolver) InitiatedAt() string {
	return r.c.InitiatedAt
}

func (r *caseResolver) ClosedAt() *string {
	if r.c.Status == "" {
		return nil
	}
	k := string(r.c.ClosedAt)
	return &k
}

func (r *caseResolver) Status() *string {
	if r.c.Status == "" {
		return nil
	}
	k := string(r.c.Status)
	return &k
}

func (r *caseResolver) FlagCount() string {
	return r.c.FlagCount
}

var schema *graphql.Schema

func init() {
	// populating our dummy data casefile map (aka casefileData)
	for _, c := range casefiles {
		casefileData[c.ID] = c
	}

	s, err := getSchema("./schema.graphql")
	if err != nil {
		panic(err)
	}

	schema = graphql.MustParseSchema(s, &Resolver{})
}

func getSchema(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "oh no!", err
	}

	return string(b), nil
}

func main() {
	// handler that reads our query viewer page
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	// required relay handler for the `graph-gophers` server gets the schema passed in
	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// var page renders the Graphiql query viewer w/ the help of introspect json
var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
