package errdetails

import "google.golang.org/grpc/codes"

func JSONParse(err error) error {
	return ErrInfo(codes.InvalidArgument, err.Error(), "JSON_PARSE_ERROR")
}
