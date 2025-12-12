module 78-pflops/services/http_gateway

go 1.25.1

require (
	78-pflops/services/ad_service v0.0.0
	78-pflops/services/user_service v0.0.0
	github.com/golang/protobuf v1.5.4
	google.golang.org/grpc v1.76.0
)

require (
	golang.org/x/net v0.45.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace 78-pflops/services/user_service => ../user_service

replace 78-pflops/services/ad_service => ../ad_service
