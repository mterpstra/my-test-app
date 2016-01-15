echo $GOPATH
ls -R $GOPATH
go get github.com/go-sql-driver/mysql
go install github.com/go-sql-driver/mysql
go get github.com/tools/godep
go install github.com/tools/godep
go build -o bin/application
