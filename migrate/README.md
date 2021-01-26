## migrate

These are tools for easier small-scale data migration to scopes and collections in Couchbase 7.0, from PostgreSQL and earlier versions of Couchbase. They use a JSON-based approach for representing table and bucket organization respectively.  

### Migrating from PostgreSQL: 
The tables of a particular schema are organised into scopes and collections in a JSON file. Each collection is mapped to a particular table and each scope to a schema, based on the mapping in [this article on the blog](https://blog.couchbase.com/scopes-and-collections-for-modern-multi-tenant-applications-couchbase-7-0/).  
Once generated, this can be modified by the user, with modifications ranging from changing scope/collection names, adding multiple tables to a collection, re-organizing scopes and collections, deleting scopes/collections, etc.    
The final JSON then forms the basis of the organisation of the specific Couchbase bucket.  
The data is then imported to a CSV file as an intermediate step and finally, imported to the specified collection based on the JSON file generated. 

### Migrating from earlier versions of Couchbase:  
Upgrading to Couchbase 7.0 will move all data to the `_default` collection. This is used for conveniently separating data in the `_default` scope to separate collections. 
The user specifies the name of the scopes and collections to be created, along with a key and a value of the key for each collection.Indices are created based using this key as the field. Documents are added to a JSON file before upload.`cbimport` is used here, over an `INSERT-SELECT`, to upload documents as JSON objects due to better performance when the number of documents is large.

### Steps to run: 
#### Migration from PostgreSQL to Couchbase:  
1. `.couchgres` is the config file for postgreSQL credentials and `config.sh` is for the Couchbase container. Modify `.couchgres` and `config.sh` variables according to requirement and have PostgreSQL and a docker container running Couchbase 7.0-beta.     
2. Run `go run db.go`.  
3. Open `public.json` and modify the organisation according to requirements.   
4. Run `go run migrate.go`.  
5. View your bucket - the data should be imported!   

#### Migration to multiple scopes and collections:  
1. Modify the variables in `config.sh` based on the Couchbase container and have a docker container running Couchbase 7.0-beta. 
2. Create `split.json` using a similar format: 
```
[
 {
  "Scope": "scope_name",
  "Key": "column_name",
  "Collections": [
   {
    "coll1": "value_for_coll1",
    "coll2": "value_for_coll2"
   }
  ]
 }
]
```   
3. Run `go run default.go`. This can take a few minutes to run with buckets with more than 5000 documents.  
