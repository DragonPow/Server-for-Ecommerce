// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: app/warehouse_service/api/api.proto

package api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// Validate checks the field values on CreateImportRequestRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateImportRequestRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ProductId

	// no validation rules for Quantity

	// no validation rules for ActionById

	return nil
}

// CreateImportRequestRequestValidationError is the validation error returned
// by CreateImportRequestRequest.Validate if the designated constraints aren't met.
type CreateImportRequestRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateImportRequestRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateImportRequestRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateImportRequestRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateImportRequestRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateImportRequestRequestValidationError) ErrorName() string {
	return "CreateImportRequestRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateImportRequestRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateImportRequestRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateImportRequestRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateImportRequestRequestValidationError{}

// Validate checks the field values on CreateImportRequestResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateImportRequestResponse) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Code

	// no validation rules for Message

	if v, ok := interface{}(m.GetData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateImportRequestResponseValidationError{
				field:  "Data",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// CreateImportRequestResponseValidationError is the validation error returned
// by CreateImportRequestResponse.Validate if the designated constraints
// aren't met.
type CreateImportRequestResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateImportRequestResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateImportRequestResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateImportRequestResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateImportRequestResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateImportRequestResponseValidationError) ErrorName() string {
	return "CreateImportRequestResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateImportRequestResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateImportRequestResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateImportRequestResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateImportRequestResponseValidationError{}

// Validate checks the field values on CreateImportRequestResponse_Data with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *CreateImportRequestResponse_Data) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ImportId

	return nil
}

// CreateImportRequestResponse_DataValidationError is the validation error
// returned by CreateImportRequestResponse_Data.Validate if the designated
// constraints aren't met.
type CreateImportRequestResponse_DataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateImportRequestResponse_DataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateImportRequestResponse_DataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateImportRequestResponse_DataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateImportRequestResponse_DataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateImportRequestResponse_DataValidationError) ErrorName() string {
	return "CreateImportRequestResponse_DataValidationError"
}

// Error satisfies the builtin error interface
func (e CreateImportRequestResponse_DataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateImportRequestResponse_Data.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateImportRequestResponse_DataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateImportRequestResponse_DataValidationError{}
