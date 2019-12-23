At locahost:8080
`go run main.go`

Query casefile looks like:

```
{
  casefile(id: 1000) {
    id
    number
    initiatedAt
    closedAt
    status
    flagCount
  }
}
```

Response:

```
{
  "data": {
    "casefile": {
      "id": "1000",
      "number": "1",
      "initiatedAt": "2019-01-02T15:04:05.000+0:00",
      "closedAt": "2019-02-02T15:04:05.000+0:00",
      "status": "Denied",
      "flagCount": "42"
    }
  }
}
```

possible IDs
1000, 2000, 3000, 4000

Query casefileEntry looks like:

```
{
  casefileEntry(id: 2000) {
    name
    id
    __typename
  }
}
```

Response:

```
{
  "data": {
    "casefileEntry": {
      "name": "Employment History",
      "id": "2000",
      "__typename": "SF86section"
    }
  }
}
```

possible IDs
1000, 1001, 1002, 1003, 1004 (financial entries)
2000, 2001 (SF86 sections)

Also, variants with and without fields. For example:

```
{
  casefileEntry(id: 1000) {
    __typename
  }
}
```

Response

```
{
  "data": {
    "casefileEntry": {
      "__typename": "FinancialEntry"
    }
  }
}
```
