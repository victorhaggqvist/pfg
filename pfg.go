package pfg

import (
	"encoding/json"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type httpError struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Details []proto.Message `json:"details,omitempty"`
}

// ErrorDetail message
func ErrorDetail(code codes.Code, msg string, details ...proto.Message) error {
	anys := make([]*any.Any, 0, len(details))

	for _, det := range details {
		pany, err := ptypes.MarshalAny(det)
		if err != nil {
			return status.Error(codes.Internal, "Internal Server Error")
		}

		anys = append(anys, pany)
	}

	out := spb.Status{
		Code:    int32(code),
		Message: msg,
		Details: anys,
	}

	return status.ErrorProto(&out)
}

//var internalServerError = httpError{
//Code:    503,
//Message: http.StatusText(http.StatusInternalServerError),
//}

var internalServerErrorBuff = []byte(`{"code":503,"message":"Internal Server Error"}`)

// ErrorHandler for gRPC JSON Gateway runtime
func ErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, rw http.ResponseWriter, req *http.Request, err error) {

	errStatus, ok := status.FromError(err)
	if !ok {
		runtime.DefaultHTTPError(ctx, mux, marshaler, rw, req, err)
		return
	}

	returnError := httpError{
		Message: errStatus.Message(),
		Details: []proto.Message{},
	}

	httpCode := runtime.HTTPStatusFromCode(errStatus.Code())
	returnError.Code = httpCode

	rw.Header().Set("Content-Type", marshaler.ContentType())
	details := errStatus.Proto().GetDetails()
	for _, det := range details {
		var target ptypes.DynamicAny
		err := ptypes.UnmarshalAny(det, &target)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write(internalServerErrorBuff)
			return
		}

		returnError.Details = append(returnError.Details, target.Message)
	}

	outBuff, err := json.Marshal(returnError)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(internalServerErrorBuff)
		return
	}
	rw.WriteHeader(httpCode)
	rw.Write(outBuff)
}
