package plugininterceptor

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TODO(nikolabo): synchronize access to these?
var currentClientChain grpc.UnaryClientInterceptor
var currentServerChain grpc.UnaryServerInterceptor
var highestFile string
var pluginPrefix string
var pluginInterface interceptInit

type interceptInit interface {
	ClientInterceptor() grpc.UnaryClientInterceptor
	ServerInterceptor() grpc.UnaryServerInterceptor
	Kill() // call to disable weak synchronization goroutine in plugin
}

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
	pluginPrefix = pluginPrefixPath
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Add unique id to rpcs
		rpc_id := rand.Uint32()
		ctx = metadata.AppendToOutgoingContext(ctx, "appnet-rpc-id", strconv.FormatUint(uint64(rpc_id), 10))

		// Add config-version header
		configVersion := "1" // XZ: temp
		ctx = metadata.AppendToOutgoingContext(ctx, "appnet-config-version", configVersion)

		if currentClientChain == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return currentClientChain(ctx, method, req, reply, cc, invoker, opts...)
	}
}

func ServerInterceptor(pluginPrefixPath string) grpc.UnaryServerInterceptor {
	pluginPrefix = pluginPrefixPath
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
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
		intercept := loadInterceptors(dir + highestFile)
		if pluginInterface != nil {
			pluginInterface.Kill()
		}
		pluginInterface = intercept
		currentClientChain = intercept.ClientInterceptor()
		currentServerChain = intercept.ServerInterceptor()
	}
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
