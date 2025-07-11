// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pb/sc.v1.proto

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

// Validate checks the field values on ListStorageClassRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListStorageClassRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListStorageClassRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListStorageClassRequestMultiError, or nil if none found.
func (m *ListStorageClassRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListStorageClassRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetClusterName()) < 1 {
		err := ListStorageClassRequestValidationError{
			field:  "ClusterName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Keyword

	// no validation rules for Page

	// no validation rules for PageSize

	// no validation rules for ListAll

	if len(errors) > 0 {
		return ListStorageClassRequestMultiError(errors)
	}

	return nil
}

// ListStorageClassRequestMultiError is an error wrapping multiple validation
// errors returned by ListStorageClassRequest.ValidateAll() if the designated
// constraints aren't met.
type ListStorageClassRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListStorageClassRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListStorageClassRequestMultiError) AllErrors() []error { return m }

// ListStorageClassRequestValidationError is the validation error returned by
// ListStorageClassRequest.Validate if the designated constraints aren't met.
type ListStorageClassRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListStorageClassRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListStorageClassRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListStorageClassRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListStorageClassRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListStorageClassRequestValidationError) ErrorName() string {
	return "ListStorageClassRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListStorageClassRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListStorageClassRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListStorageClassRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListStorageClassRequestValidationError{}

// Validate checks the field values on StorageClassItem with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *StorageClassItem) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StorageClassItem with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StorageClassItemMultiError, or nil if none found.
func (m *StorageClassItem) ValidateAll() error {
	return m.validate(true)
}

func (m *StorageClassItem) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := StorageClassItemValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for CreateTime

	// no validation rules for IsFollowed

	// no validation rules for Yaml

	// no validation rules for Provisioner

	// no validation rules for ReclaimPolicy

	// no validation rules for VolumeBindingMode

	// no validation rules for IsDefault

	if len(errors) > 0 {
		return StorageClassItemMultiError(errors)
	}

	return nil
}

// StorageClassItemMultiError is an error wrapping multiple validation errors
// returned by StorageClassItem.ValidateAll() if the designated constraints
// aren't met.
type StorageClassItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StorageClassItemMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StorageClassItemMultiError) AllErrors() []error { return m }

// StorageClassItemValidationError is the validation error returned by
// StorageClassItem.Validate if the designated constraints aren't met.
type StorageClassItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StorageClassItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StorageClassItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StorageClassItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StorageClassItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StorageClassItemValidationError) ErrorName() string { return "StorageClassItemValidationError" }

// Error satisfies the builtin error interface
func (e StorageClassItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStorageClassItem.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StorageClassItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StorageClassItemValidationError{}

// Validate checks the field values on ListStorageClassResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListStorageClassResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListStorageClassResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListStorageClassResponseMultiError, or nil if none found.
func (m *ListStorageClassResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListStorageClassResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Total

	for idx, item := range m.GetList() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListStorageClassResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListStorageClassResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListStorageClassResponseValidationError{
					field:  fmt.Sprintf("List[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListStorageClassResponseMultiError(errors)
	}

	return nil
}

// ListStorageClassResponseMultiError is an error wrapping multiple validation
// errors returned by ListStorageClassResponse.ValidateAll() if the designated
// constraints aren't met.
type ListStorageClassResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListStorageClassResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListStorageClassResponseMultiError) AllErrors() []error { return m }

// ListStorageClassResponseValidationError is the validation error returned by
// ListStorageClassResponse.Validate if the designated constraints aren't met.
type ListStorageClassResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListStorageClassResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListStorageClassResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListStorageClassResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListStorageClassResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListStorageClassResponseValidationError) ErrorName() string {
	return "ListStorageClassResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListStorageClassResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListStorageClassResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListStorageClassResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListStorageClassResponseValidationError{}
