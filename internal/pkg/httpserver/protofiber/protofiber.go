package protofiber

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Protofiber struct {
	mOpts *protojson.MarshalOptions
	uOpts *protojson.UnmarshalOptions
}

// Opts is an options struct for marshaller and unmarshaller
type Opts struct {
	// UseProtoNames uses proto field name instead of lowerCamelCase name in JSON
	// field names.
	UseProtoNames bool

	// UseEnumNumbers emits enum values as numbers.
	UseEnumNumbers bool
}

func NewProtofiber(opts ...Opts) *Protofiber {
	p := &Protofiber{
		mOpts: &protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
		uOpts: &protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	if len(opts) == 0 {
		return p
	}
	opt := opts[0]
	if opt.UseProtoNames {
		p.mOpts.UseProtoNames = true
	}
	if opt.UseEnumNumbers {
		p.mOpts.UseEnumNumbers = true
	}
	return p
}

func (p *Protofiber) Marhshal(value interface{}) ([]byte, error) {
	v, ok := value.(proto.Message)
	if !ok {
		return json.Marshal(value)
	}
	return p.mOpts.Marshal(v)
}

func (p *Protofiber) Unmarshal(data []byte, out interface{}) error {
	o := out.(proto.Message)
	return p.uOpts.Unmarshal(data, o)
}

func (p *Protofiber) CtxMarshal(ctx *fiber.Ctx, m proto.Message) error {
	b, err := p.mOpts.Marshal(m)
	if err != nil {
		return err
	}
	ctx.Response().SetBodyRaw(b)
	ctx.Response().Header.SetContentType(fiber.MIMEApplicationJSON)
	return nil
}

func (p *Protofiber) CtxUnmarshal(ctx *fiber.Ctx, m proto.Message) error {
	ctype := utils.ToLower(utils.UnsafeString(ctx.Request().Header.ContentType()))

	ctype = utils.ParseVendorSpecificContentType(ctype)

	if !strings.HasPrefix(ctype, fiber.MIMEApplicationJSON) {
		return fiber.ErrUnprocessableEntity
	}

	return p.uOpts.Unmarshal(ctx.Body(), m)
}
