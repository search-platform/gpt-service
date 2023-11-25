package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/search-platform/gpt-service/internal/pkg/httpserver/protofiber"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var pf = protofiber.NewProtofiber()

func convertCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return fiber.StatusBadRequest
	// case codes.Unauthenticated: // пока не актуально
	// 	return fiber.StatusUnauthorized
	case codes.DeadlineExceeded:
		return fiber.StatusGone
	case codes.Internal:
		return fiber.StatusInternalServerError
	case codes.NotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusBadRequest
	}
}

// ErrorHandler функция обработчик ошибок
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	log.Error().Stack().Err(err).Msg("request error")

	if st, ok := status.FromError(err); ok {
		if len(st.Details()) == 0 {
			code := fiber.StatusInternalServerError
			switch st.Code() {
			case codes.InvalidArgument:
				code = fiber.StatusBadRequest
			case codes.Unauthenticated:
				code = fiber.StatusUnauthorized
			}

			return ctx.Status(code).JSON(fiber.Map{
				"message": st.Message(),
				"code":    code,
			})
		}
	}
	if strings.HasPrefix(err.Error(), "proto:") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"code":    fiber.StatusBadRequest,
		})
	}
	if err, ok := err.(*fiber.Error); ok {
		return ctx.Status(err.Code).JSON(fiber.Map{
			"message": err.Message,
			"code":    err.Code,
		})
	}

	return ctx.Status(500).JSON(fiber.Map{
		"message": "internal server error",
		"code":    500,
	})
}
