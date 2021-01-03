# Scope Vis

Used for showing the collections and scopes in a particular bucket, with percentages based on number of records. Used with Couchbase 7.0 and designed for its new feature, scopes and collections.

### Prerequisites:  
1. docker - running CB:enterprise-7.0.0 container.  
2. `travel-sample` or any other bucket with scopes and collections loaded and `config.sh` details set.  
3. go 1.13+   
4. `jq` - to format JSON output

### Steps to Run:
1. Run: `go run processing.go`
2. Open `coll.html`


