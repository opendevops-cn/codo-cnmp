// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pb/user.v1.proto

package pb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
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
	_ = sort.Sort
)

// Validate checks the field values on FollowItem with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *FollowItem) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FollowItem with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in FollowItemMultiError, or
// nil if none found.
func (m *FollowItem) ValidateAll() error {
	return m.validate(true)
}

func (m *FollowItem) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for FollowType

	// no validation rules for FollowValue

	// no validation rules for Id

	// no validation rules for CreateTime

	// no validation rules for ClusterName

	if len(errors) > 0 {
		return FollowItemMultiError(errors)
	}

	return nil
}

// FollowItemMultiError is an error wrapping multiple validation errors
// returned by FollowItem.ValidateAll() if the designated constraints aren't met.
type FollowItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FollowItemMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FollowItemMultiError) AllErrors() []error { return m }

// FollowItemValidationError is the validation error returned by
// FollowItem.Validate if the designated constraints aren't met.
type FollowItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FollowItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FollowItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FollowItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FollowItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FollowItemValidationError) ErrorName() string { return "FollowItemValidationError" }

// Error satisfies the builtin error interface
func (e FollowItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFollowItem.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FollowItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FollowItemValidationError{}

// Validate checks the field values on ListUserFollowRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListUserFollowRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListUserFollowRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListUserFollowRequestMultiError, or nil if none found.
func (m *ListUserFollowRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListUserFollowRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for FollowType

	// no validation rules for Page

	// no validation rules for PageSize

	// no validation rules for ListAll

	// no validation rules for Keyword

	// no validation rules for FollowValue

	if len(errors) > 0 {
		return ListUserFollowRequestMultiError(errors)
	}

	return nil
}

// ListUserFollowRequestMultiError is an error wrapping multiple validation
// errors returned by ListUserFollowRequest.ValidateAll() if the designated
// constraints aren't met.
type ListUserFollowRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListUserFollowRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListUserFollowRequestMultiError) AllErrors() []error { return m }

// ListUserFollowRequestValidationError is the validation error returned by
// ListUserFollowRequest.Validate if the designated constraints aren't met.
type ListUserFollowRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListUserFollowRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListUserFollowRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListUserFollowRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListUserFollowRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListUserFollowRequestValidationError) ErrorName() string {
	return "ListUserFollowRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListUserFollowRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListUserFollowRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListUserFollowRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListUserFollowRequestValidationError{}

// Validate checks the field values on ListUserFollowResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListUserFollowResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListUserFollowResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListUserFollowResponseMultiError, or nil if none found.
func (m *ListUserFollowResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListUserFollowResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetList() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListUserFollowResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListUserFollowResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListUserFollowResponseValidationError{
					field:  fmt.Sprintf("List[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Total

	if len(errors) > 0 {
		return ListUserFollowResponseMultiError(errors)
	}

	return nil
}

// ListUserFollowResponseMultiError is an error wrapping multiple validation
// errors returned by ListUserFollowResponse.ValidateAll() if the designated
// constraints aren't met.
type ListUserFollowResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListUserFollowResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListUserFollowResponseMultiError) AllErrors() []error { return m }

// ListUserFollowResponseValidationError is the validation error returned by
// ListUserFollowResponse.Validate if the designated constraints aren't met.
type ListUserFollowResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListUserFollowResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListUserFollowResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListUserFollowResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListUserFollowResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListUserFollowResponseValidationError) ErrorName() string {
	return "ListUserFollowResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListUserFollowResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListUserFollowResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListUserFollowResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListUserFollowResponseValidationError{}

// Validate checks the field values on UserFollowRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserFollowRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserFollowRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserFollowRequestMultiError, or nil if none found.
func (m *UserFollowRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UserFollowRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for FollowType

	// no validation rules for FollowValue

	// no validation rules for ClusterName

	if len(errors) > 0 {
		return UserFollowRequestMultiError(errors)
	}

	return nil
}

// UserFollowRequestMultiError is an error wrapping multiple validation errors
// returned by UserFollowRequest.ValidateAll() if the designated constraints
// aren't met.
type UserFollowRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserFollowRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserFollowRequestMultiError) AllErrors() []error { return m }

// UserFollowRequestValidationError is the validation error returned by
// UserFollowRequest.Validate if the designated constraints aren't met.
type UserFollowRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserFollowRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserFollowRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserFollowRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserFollowRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserFollowRequestValidationError) ErrorName() string {
	return "UserFollowRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UserFollowRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserFollowRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserFollowRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserFollowRequestValidationError{}

// Validate checks the field values on UserFollowResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UserFollowResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserFollowResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserFollowResponseMultiError, or nil if none found.
func (m *UserFollowResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UserFollowResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return UserFollowResponseMultiError(errors)
	}

	return nil
}

// UserFollowResponseMultiError is an error wrapping multiple validation errors
// returned by UserFollowResponse.ValidateAll() if the designated constraints
// aren't met.
type UserFollowResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserFollowResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserFollowResponseMultiError) AllErrors() []error { return m }

// UserFollowResponseValidationError is the validation error returned by
// UserFollowResponse.Validate if the designated constraints aren't met.
type UserFollowResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserFollowResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserFollowResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserFollowResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserFollowResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserFollowResponseValidationError) ErrorName() string {
	return "UserFollowResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UserFollowResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserFollowResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserFollowResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserFollowResponseValidationError{}

// Validate checks the field values on DeleteUserFollowRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteUserFollowRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteUserFollowRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteUserFollowRequestMultiError, or nil if none found.
func (m *DeleteUserFollowRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteUserFollowRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for FollowType

	// no validation rules for FollowValue

	// no validation rules for ClusterName

	if len(errors) > 0 {
		return DeleteUserFollowRequestMultiError(errors)
	}

	return nil
}

// DeleteUserFollowRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteUserFollowRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteUserFollowRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteUserFollowRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteUserFollowRequestMultiError) AllErrors() []error { return m }

// DeleteUserFollowRequestValidationError is the validation error returned by
// DeleteUserFollowRequest.Validate if the designated constraints aren't met.
type DeleteUserFollowRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteUserFollowRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteUserFollowRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteUserFollowRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteUserFollowRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteUserFollowRequestValidationError) ErrorName() string {
	return "DeleteUserFollowRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteUserFollowRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteUserFollowRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteUserFollowRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteUserFollowRequestValidationError{}

// Validate checks the field values on DeleteUserFollowResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteUserFollowResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteUserFollowResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteUserFollowResponseMultiError, or nil if none found.
func (m *DeleteUserFollowResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteUserFollowResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DeleteUserFollowResponseMultiError(errors)
	}

	return nil
}

// DeleteUserFollowResponseMultiError is an error wrapping multiple validation
// errors returned by DeleteUserFollowResponse.ValidateAll() if the designated
// constraints aren't met.
type DeleteUserFollowResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteUserFollowResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteUserFollowResponseMultiError) AllErrors() []error { return m }

// DeleteUserFollowResponseValidationError is the validation error returned by
// DeleteUserFollowResponse.Validate if the designated constraints aren't met.
type DeleteUserFollowResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteUserFollowResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteUserFollowResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteUserFollowResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteUserFollowResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteUserFollowResponseValidationError) ErrorName() string {
	return "DeleteUserFollowResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteUserFollowResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteUserFollowResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteUserFollowResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteUserFollowResponseValidationError{}
