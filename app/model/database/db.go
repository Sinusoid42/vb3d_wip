package database

import(
	"github.com/leesper/couchdb-golang"
	
)


var Database *couchdb.Database

func init() {
	var err error
	Database, err = couchdb.NewDatabase("http://localhost:5984/vb3d")

	if err != nil {
		print(err)
	}
}


