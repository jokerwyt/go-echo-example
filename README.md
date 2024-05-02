# AppNet Interceptor Library

This library can be included by applications that will be used with AppNet gRPC elements. The library exposes an interceptor which performs the processing they specify via AppNet. Approximate LOC change required to applications ~= 1 line per gRPC connection. Allows rpc processing logic to be swapped at runtime via Go's plugin system.