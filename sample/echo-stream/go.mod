module github.com/appnet-org/golib/sample/echo-stream

go 1.22.1

replace github.com/appnet-org/golib/sample/echo-stream/pb => ./pb

require (
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
)
