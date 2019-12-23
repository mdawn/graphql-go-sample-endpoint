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

type financialEntry struct {
	ID        graphql.ID
	Name      string
}


var financialEntries = []*financialEntry{
	{
		ID:        "1000",
		Name:      "Bankruptcy",
	},
	{
		ID:        "1001",
		Name:      "Outstanding Private Loan",
	},
	{
		ID:        "1002",
		Name:      "Paid Private Loan",
	},
	{
		ID:        "1003",
		Name:      "Student Loan",
	},
	{
		ID:        "1004",
		Name:      "Mortgage",
	},
}

var financialEntryData = make(map[graphql.ID]*financialEntry)

func init() {
	for _, h := range financialEntries {
		financialEntryData[h.ID] = h
	}
}

type sf86section struct {
	ID              graphql.ID
	Name            string
}

var sf86sections = []*sf86section{
	{
		ID:              "2000",
		Name:            "Employment History",
	},
	{
		ID:              "2001",
		Name:            "Credit History",
	},
}

var sf86sectionData = make(map[graphql.ID]*sf86section)

func init() {
	for _, d := range sf86sections {
		sf86sectionData[d.ID] = d
	}
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

func (r *Resolver) CasefileEntry(args struct{ ID graphql.ID }) *casefileEntryResolver {
	if h := financialEntryData[args.ID]; h != nil {
		return &casefileEntryResolver{&financialEntryResolver{h}}
	}
	if d := sf86sectionData[args.ID]; d != nil {
		return &casefileEntryResolver{&sf86sectionResolver{d}}
	}
	return nil
}

func (r *Resolver) FinancialEntry(args struct{ ID graphql.ID }) *financialEntryResolver {
	if h := financialEntryData[args.ID]; h != nil {
		return &financialEntryResolver{h}
	}
	return nil
}

func (r *Resolver) SF86section(args struct{ ID graphql.ID }) *sf86sectionResolver {
	if d := sf86sectionData[args.ID]; d != nil {
		return &sf86sectionResolver{d}
	}
	return nil
}

type casefileEntry interface {
	ID() graphql.ID
	Name() string
}

type casefileEntryResolver struct {
	casefileEntry
}

func (r *casefileEntryResolver) ToFinancialEntry() (*financialEntryResolver, bool) {
	c, ok := r.casefileEntry.(*financialEntryResolver)
	return c, ok
}

func (r *casefileEntryResolver) ToSF86section() (*sf86sectionResolver, bool) {
	c, ok := r.casefileEntry.(*sf86sectionResolver)
	return c, ok
}

type financialEntryResolver struct {
	h *financialEntry
}

func (r *financialEntryResolver) ID() graphql.ID {
	return r.h.ID
}

func (r *financialEntryResolver) Name() string {
	return r.h.Name
}

type sf86sectionResolver struct {
	d *sf86section
}

func (r *sf86sectionResolver) ID() graphql.ID {
	return r.d.ID
}

func (r *sf86sectionResolver) Name() string {
	return r.d.Name
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
