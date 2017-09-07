# pfg
--
    import "github.com/victorhaggqvist/pfg"

Package pfg provides error helpers for detailed return errors to be used with
https://github.com/grpc-ecosystem/grpc-gateway

## Usage

#### func  ErrorDetail

```go
func ErrorDetail(code codes.Code, msg string, details ...proto.Message) error
```
ErrorDetail creates an error with bundled detail messages

#### func  ErrorHandler

```go
func ErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, rw http.ResponseWriter, req *http.Request, err error)
```
ErrorHandler for gRPC JSON Gateway runtime

Usage

    gw := runtime.NewServeMux(
    	runtime.WithProtoErrorHandler(pfg.ErrorHandler),
    	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONBuiltin{}),
    )
