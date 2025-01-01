module github.com/cheracc/fortress-grpc

go 1.23.4

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.35.1
)

require cloud.google.com/go/compute/metadata v0.5.2 // indirect

require (
	github.com/mattn/go-sqlite3 v1.14.24
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/oauth2 v0.24.0
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
)
