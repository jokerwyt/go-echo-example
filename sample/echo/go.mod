module github.com/appnet-org/golib/sample/echo

go 1.22.1

// toolchain go1.22.2

require (
	github.com/appnet-org/golib/interceptor v0.0.0-00010101000000-000000000000
	github.com/appnet-org/golib/sample/echo-pb v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.24.0
	google.golang.org/grpc v1.63.2
)

replace github.com/appnet-org/golib/sample/echo-pb => ../echo-pb

replace github.com/appnet-org/golib/interceptor => ../../interceptor
