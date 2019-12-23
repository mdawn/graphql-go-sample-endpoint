At locahost:8080

Query looks like:

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

possible IDs
1000, 2000, 3000, 4000

```
{
  casefileEntry (id: 2000){
    name
    id
    __typename
  }
}
```

possible IDs
1000, 1001, 1002, 1003, 1004
2000, 2001

Also, for example:

```
{
  casefileEntry(id: 1000) {
    __typename
  }
}
```

returns

```
{
  "data": {
    "casefileEntry": {
      "__typename": "FinancialEntry"
    }
  }
}
```
