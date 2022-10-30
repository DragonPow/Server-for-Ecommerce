// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: app/order_service/api/data.proto

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

// Validate checks the field values on GetOrderDetailResponseData with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GetOrderDetailResponseData) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// GetOrderDetailResponseDataValidationError is the validation error returned
// by GetOrderDetailResponseData.Validate if the designated constraints aren't met.
type GetOrderDetailResponseDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetOrderDetailResponseDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetOrderDetailResponseDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetOrderDetailResponseDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetOrderDetailResponseDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetOrderDetailResponseDataValidationError) ErrorName() string {
	return "GetOrderDetailResponseDataValidationError"
}

// Error satisfies the builtin error interface
func (e GetOrderDetailResponseDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetOrderDetailResponseData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetOrderDetailResponseDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetOrderDetailResponseDataValidationError{}
