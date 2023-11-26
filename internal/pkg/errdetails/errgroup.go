package errdetails

import (
	"github.com/search-platform/gpt-service/api/errdetails"
	"google.golang.org/grpc/codes"
)

type BadRequestBuilder interface {
	Append(errtype errdetails.BadRequest_FieldViolationType, field, description string) BadRequestBuilder
	Required(field string, message string) BadRequestBuilder
	Invalid(field string, message string) BadRequestBuilder
	Unique(field string, message string) BadRequestBuilder
	NotEmpty() bool
	AsError() error
	WithContext(field string) BadRequestBuilder
}

type group struct {
	errdetails.BadRequest
}

func NewBadRequestBuilder() BadRequestBuilder {
	return &group{}
}

func (g *group) Append(errtype errdetails.BadRequest_FieldViolationType, field, description string) BadRequestBuilder {
	g.FieldViolations = append(g.FieldViolations, &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Type:        errtype,
		Description: description,
	})

	return g
}

func (g *group) Required(field string, description string) BadRequestBuilder {
	if description == "" {
		description = "required"
	}

	return g.Append(errdetails.BadRequest_REQUIRED, field, description)
}

func (g *group) Invalid(field string, description string) BadRequestBuilder {
	if description == "" {
		description = "not valid"
	}
	return g.Append(errdetails.BadRequest_NOT_VALID, field, description)
}

func (g *group) Unique(field string, description string) BadRequestBuilder {
	if description == "" {
		description = "not unique"
	}
	return g.Append(errdetails.BadRequest_UNIQUE, field, description)
}

func (g *group) NotEmpty() bool {
	return len(g.FieldViolations) > 0
}

func (g *group) AsError() error {
	if len(g.FieldViolations) == 0 {
		return nil
	}
	return NewError(codes.InvalidArgument, "bad request", &g.BadRequest)
}

func (g *group) WithContext(field string) BadRequestBuilder {
	return &groupContext{group: g, name: field}
}

type groupContext struct {
	*group

	name string
}

func (g *groupContext) field(field string) string {
	if field == "." {
		return g.name
	}
	if field[0] == '[' {
		return g.name + field
	}
	return g.name + "." + field
}

func (g *groupContext) Append(errtype errdetails.BadRequest_FieldViolationType, field, message string) BadRequestBuilder {
	g.group.Append(errtype, g.field(field), message)
	return g
}

func (g *groupContext) Required(field string, message string) BadRequestBuilder {
	g.group.Required(g.field(field), message)
	return g
}

func (g *groupContext) Invalid(field string, message string) BadRequestBuilder {
	g.group.Invalid(g.field(field), message)
	return g
}

func (g *groupContext) Unique(field string, description string) BadRequestBuilder {
	g.group.Unique(g.field(field), description)
	return g
}

func (g *groupContext) WithContext(field string) BadRequestBuilder {
	return &groupContext{group: g.group, name: g.field(field)}
}
