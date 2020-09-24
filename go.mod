module json2proto

go 1.12

require (
	comm v0.0.0
	github.com/droundy/goopt v0.0.0-20170604162106-0b8effe182da
	github.com/go-redis/redis/v7 v7.2.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/guregu/null v3.4.0+incompatible
	github.com/jinzhu/gorm v1.9.10
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/parnurzeal/gorequest v0.2.16
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
	moul.io/http2curl v1.0.0 // indirect
	proto v0.0.0
)

replace (
	comm => ../comm
	proto => ../proto
)
