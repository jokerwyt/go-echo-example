package plugininterceptor

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// TODO(nikolabo): synchronize access to these?
var currentClientChain grpc.UnaryClientInterceptor
var currentServerChain grpc.UnaryServerInterceptor
var highestFile string
var pluginPrefix string

func init() {
	go func() {
		for {
			if pluginPrefix != "" {
				updateChains(pluginPrefix)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func ClientInterceptor(pluginPrefixPath string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// if pluginPrefix == "" {
		// 	updateChains(pluginPrefixPath)
		// }
		pluginPrefix = pluginPrefixPath

		if currentClientChain == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return currentClientChain(ctx, method, req, reply, cc, invoker, opts...)
	}
}

func ServerInterceptor(pluginPrefixPath string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// if pluginPrefix == "" {
		// 	updateChains(pluginPrefixPath)
		// }
		pluginPrefix = pluginPrefixPath

		if currentServerChain == nil {
			return handler(ctx, req)
		}

		return currentServerChain(ctx, req, info, handler)
	}
}

func updateChains(prefix string) {
	var highestSeen string

	dir, prefix := filepath.Split(prefix)
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			if file.Name() > highestFile {
				highestSeen = file.Name()
			}
		}
	}

	if highestSeen != "" && highestSeen != highestFile {
		highestFile = highestSeen
		interceptInit := loadInterceptors(dir + highestFile)
		currentClientChain = interceptInit.ClientInterceptor()
		currentServerChain = interceptInit.ServerInterceptor()
	}
}
