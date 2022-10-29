// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: app/warehouse_service/api/data.proto

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

// Validate checks the field values on CreateImportBillResponseData with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateImportBillResponseData) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ImportId

	return nil
}

// CreateImportBillResponseDataValidationError is the validation error returned
// by CreateImportBillResponseData.Validate if the designated constraints
// aren't met.
type CreateImportBillResponseDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateImportBillResponseDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateImportBillResponseDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateImportBillResponseDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateImportBillResponseDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateImportBillResponseDataValidationError) ErrorName() string {
	return "CreateImportBillResponseDataValidationError"
}

// Error satisfies the builtin error interface
func (e CreateImportBillResponseDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateImportBillResponseData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateImportBillResponseDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateImportBillResponseDataValidationError{}

// Validate checks the field values on GetImportBillResponseData with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GetImportBillResponseData) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetItem()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetImportBillResponseDataValidationError{
				field:  "Item",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// GetImportBillResponseDataValidationError is the validation error returned by
// GetImportBillResponseData.Validate if the designated constraints aren't met.
type GetImportBillResponseDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetImportBillResponseDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetImportBillResponseDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetImportBillResponseDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetImportBillResponseDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetImportBillResponseDataValidationError) ErrorName() string {
	return "GetImportBillResponseDataValidationError"
}

// Error satisfies the builtin error interface
func (e GetImportBillResponseDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetImportBillResponseData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetImportBillResponseDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetImportBillResponseDataValidationError{}

// Validate checks the field values on GetImportBillItem with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *GetImportBillItem) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for Code

	// no validation rules for LastActionById

	// no validation rules for LastActionByName

	// no validation rules for CreateById

	// no validation rules for CreateByName

	for idx, item := range m.GetItemDetails() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetImportBillItemValidationError{
					field:  fmt.Sprintf("ItemDetails[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// GetImportBillItemValidationError is the validation error returned by
// GetImportBillItem.Validate if the designated constraints aren't met.
type GetImportBillItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetImportBillItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetImportBillItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetImportBillItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetImportBillItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetImportBillItemValidationError) ErrorName() string {
	return "GetImportBillItemValidationError"
}

// Error satisfies the builtin error interface
func (e GetImportBillItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetImportBillItem.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetImportBillItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetImportBillItemValidationError{}

// Validate checks the field values on GetImportBillItemDetail with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GetImportBillItemDetail) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ProductId

	// no validation rules for ProductName

	// no validation rules for UomName

	// no validation rules for Quantity

	return nil
}

// GetImportBillItemDetailValidationError is the validation error returned by
// GetImportBillItemDetail.Validate if the designated constraints aren't met.
type GetImportBillItemDetailValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetImportBillItemDetailValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetImportBillItemDetailValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetImportBillItemDetailValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetImportBillItemDetailValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetImportBillItemDetailValidationError) ErrorName() string {
	return "GetImportBillItemDetailValidationError"
}

// Error satisfies the builtin error interface
func (e GetImportBillItemDetailValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetImportBillItemDetail.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetImportBillItemDetailValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetImportBillItemDetailValidationError{}
