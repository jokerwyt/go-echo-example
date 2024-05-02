package plugininterceptor

import (
	"fmt"
	"plugin"

	"google.golang.org/grpc"
)

type interceptInit interface {
	ClientInterceptor() grpc.UnaryClientInterceptor
	ServerInterceptor() grpc.UnaryServerInterceptor
}

func loadInterceptors(interceptorPluginPath string) interceptInit {
	// TODO: return err instead of panicking
	interceptorPlugin, err := plugin.Open(interceptorPluginPath)
	if err != nil {
		fmt.Printf("loading error: %v\n", err)
		panic("error loading interceptor plugin so")
	}

	symInterceptInit, err := interceptorPlugin.Lookup("InterceptInit")
	if err != nil {
		panic("error locating interceptor in plugin so")
	}

	interceptInit, ok := symInterceptInit.(interceptInit)
	if !ok {
		panic("error casting interceptInit")
	}

	return interceptInit
}
