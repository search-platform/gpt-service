package errdetails

import (
	"fmt"

	"github.com/search-platform/gpt-service/api/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
)

func NewError(code codes.Code, message string, details ...protoiface.MessageV1) error {
	st := status.New(codes.Code(code), message)
	if details == nil {
		return st.Err()
	}
	std, err := st.WithDetails(details...)
	if err != nil {
		return st.Err()
	}
	return std.Err()
}

func NewBadRequest(fieldViolations []*errdetails.BadRequest_FieldViolation) error {
	verr := errdetails.BadRequest{
		FieldViolations: fieldViolations,
	}
	st := status.New(codes.InvalidArgument, "bad request")
	std, err := st.WithDetails(&verr)
	if err != nil {
		return st.Err()
	}
	return std.Err()
}

func BadRequestFromError(err error) (*errdetails.BadRequest, bool) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}

	detail := st.Details()[0]
	badreq, ok := detail.(*errdetails.BadRequest)
	if !ok {
		return nil, false
	}

	return badreq, true
}

func NotFound(message string, args ...interface{}) error {
	return ErrInfo(codes.NotFound, fmt.Sprintf(message, args...), "NOT_FOUND")
}

func ErrInfo(code codes.Code, message, errType string, details ...protoiface.MessageV1) error {
	terr := errdetails.ErrorInfo{
		Reason:      errType,
		Description: message,
	}
	st := status.New(code, message)
	std, err := st.WithDetails(&terr)
	if err != nil {
		return st.Err()
	}
	return std.Err()
}
