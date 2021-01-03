## migrate

This is a tool for easier small-scale data migration from PostgreSQL to Couchbase 7.0 with scopes and collections.   
It uses a JSON-based approach where the tables of a particular schema are organised into scopes and collections in a JSON file. Each collection is mapped to a particular table and each scope to a schema, based on the mapping in [this article on the blog](https://blog.couchbase.com/scopes-and-collections-for-modern-multi-tenant-applications-couchbase-7-0/).  
Once generated, this can be modified by the user, with modifications ranging from changing scope/collection names, adding multiple tables to a collection, re-organizing scopes and collections, deleting scopes/collections, etc.    
The final JSON then forms the basis of the organisation of the specific Couchbase bucket.  

### Steps to run:  
1. Modify `config.sh` variables according to requirement and have PostgreSQL and a docker container running Couchbase 7.0-beta running.     
2. Run `go run db.go`.  
3. Open `public.json` and modify the organisation according to requirements.   
4. Run `go run migrate.go`.  
5. View your bucket - the data should be imported!  
