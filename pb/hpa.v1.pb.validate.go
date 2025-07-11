// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pb/hpa.v1.proto

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

// Validate checks the field values on ListHpaRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ListHpaRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListHpaRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ListHpaRequestMultiError,
// or nil if none found.
func (m *ListHpaRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListHpaRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetClusterName()) < 1 {
		err := ListHpaRequestValidationError{
			field:  "ClusterName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Namespace

	// no validation rules for Keyword

	// no validation rules for Page

	// no validation rules for PageSize

	// no validation rules for ListAll

	if len(errors) > 0 {
		return ListHpaRequestMultiError(errors)
	}

	return nil
}

// ListHpaRequestMultiError is an error wrapping multiple validation errors
// returned by ListHpaRequest.ValidateAll() if the designated constraints
// aren't met.
type ListHpaRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListHpaRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListHpaRequestMultiError) AllErrors() []error { return m }

// ListHpaRequestValidationError is the validation error returned by
// ListHpaRequest.Validate if the designated constraints aren't met.
type ListHpaRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListHpaRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListHpaRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListHpaRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListHpaRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListHpaRequestValidationError) ErrorName() string { return "ListHpaRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListHpaRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListHpaRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListHpaRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListHpaRequestValidationError{}

// Validate checks the field values on HpaItem with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *HpaItem) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HpaItem with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in HpaItemMultiError, or nil if none found.
func (m *HpaItem) ValidateAll() error {
	return m.validate(true)
}

func (m *HpaItem) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := HpaItemValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetNamespace()) < 1 {
		err := HpaItemValidationError{
			field:  "Namespace",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for WorkloadType

	// no validation rules for Workload

	// no validation rules for MinReplicas

	// no validation rules for MaxReplicas

	// no validation rules for TargetCpuUtilization

	// no validation rules for CurrentCpuUtilization

	// no validation rules for TargetMemoryUtilization

	// no validation rules for CurrentMemoryUtilization

	// no validation rules for Labels

	// no validation rules for CreateTime

	// no validation rules for Yaml

	// no validation rules for Annotations

	// no validation rules for UpdateTime

	// no validation rules for CurrentReplicas

	// no validation rules for IsFollowed

	if len(errors) > 0 {
		return HpaItemMultiError(errors)
	}

	return nil
}

// HpaItemMultiError is an error wrapping multiple validation errors returned
// by HpaItem.ValidateAll() if the designated constraints aren't met.
type HpaItemMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HpaItemMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HpaItemMultiError) AllErrors() []error { return m }

// HpaItemValidationError is the validation error returned by HpaItem.Validate
// if the designated constraints aren't met.
type HpaItemValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HpaItemValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HpaItemValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HpaItemValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HpaItemValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HpaItemValidationError) ErrorName() string { return "HpaItemValidationError" }

// Error satisfies the builtin error interface
func (e HpaItemValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHpaItem.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HpaItemValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HpaItemValidationError{}

// Validate checks the field values on ListHpaResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListHpaResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListHpaResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListHpaResponseMultiError, or nil if none found.
func (m *ListHpaResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListHpaResponse) validate(all bool) error {
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
					errors = append(errors, ListHpaResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListHpaResponseValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListHpaResponseValidationError{
					field:  fmt.Sprintf("List[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListHpaResponseMultiError(errors)
	}

	return nil
}

// ListHpaResponseMultiError is an error wrapping multiple validation errors
// returned by ListHpaResponse.ValidateAll() if the designated constraints
// aren't met.
type ListHpaResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListHpaResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListHpaResponseMultiError) AllErrors() []error { return m }

// ListHpaResponseValidationError is the validation error returned by
// ListHpaResponse.Validate if the designated constraints aren't met.
type ListHpaResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListHpaResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListHpaResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListHpaResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListHpaResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListHpaResponseValidationError) ErrorName() string { return "ListHpaResponseValidationError" }

// Error satisfies the builtin error interface
func (e ListHpaResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListHpaResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListHpaResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListHpaResponseValidationError{}

// Validate checks the field values on CreateOrUpdateHpaByYamlRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateOrUpdateHpaByYamlRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateOrUpdateHpaByYamlRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// CreateOrUpdateHpaByYamlRequestMultiError, or nil if none found.
func (m *CreateOrUpdateHpaByYamlRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateOrUpdateHpaByYamlRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetClusterName()) < 1 {
		err := CreateOrUpdateHpaByYamlRequestValidationError{
			field:  "ClusterName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetYaml()) < 1 {
		err := CreateOrUpdateHpaByYamlRequestValidationError{
			field:  "Yaml",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateOrUpdateHpaByYamlRequestMultiError(errors)
	}

	return nil
}

// CreateOrUpdateHpaByYamlRequestMultiError is an error wrapping multiple
// validation errors returned by CreateOrUpdateHpaByYamlRequest.ValidateAll()
// if the designated constraints aren't met.
type CreateOrUpdateHpaByYamlRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateOrUpdateHpaByYamlRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateOrUpdateHpaByYamlRequestMultiError) AllErrors() []error { return m }

// CreateOrUpdateHpaByYamlRequestValidationError is the validation error
// returned by CreateOrUpdateHpaByYamlRequest.Validate if the designated
// constraints aren't met.
type CreateOrUpdateHpaByYamlRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateOrUpdateHpaByYamlRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateOrUpdateHpaByYamlRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateOrUpdateHpaByYamlRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateOrUpdateHpaByYamlRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateOrUpdateHpaByYamlRequestValidationError) ErrorName() string {
	return "CreateOrUpdateHpaByYamlRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateOrUpdateHpaByYamlRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateOrUpdateHpaByYamlRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateOrUpdateHpaByYamlRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateOrUpdateHpaByYamlRequestValidationError{}

// Validate checks the field values on CreateOrUpdateHpaByYamlResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateOrUpdateHpaByYamlResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateOrUpdateHpaByYamlResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// CreateOrUpdateHpaByYamlResponseMultiError, or nil if none found.
func (m *CreateOrUpdateHpaByYamlResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateOrUpdateHpaByYamlResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return CreateOrUpdateHpaByYamlResponseMultiError(errors)
	}

	return nil
}

// CreateOrUpdateHpaByYamlResponseMultiError is an error wrapping multiple
// validation errors returned by CreateOrUpdateHpaByYamlResponse.ValidateAll()
// if the designated constraints aren't met.
type CreateOrUpdateHpaByYamlResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateOrUpdateHpaByYamlResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateOrUpdateHpaByYamlResponseMultiError) AllErrors() []error { return m }

// CreateOrUpdateHpaByYamlResponseValidationError is the validation error
// returned by CreateOrUpdateHpaByYamlResponse.Validate if the designated
// constraints aren't met.
type CreateOrUpdateHpaByYamlResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateOrUpdateHpaByYamlResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateOrUpdateHpaByYamlResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateOrUpdateHpaByYamlResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateOrUpdateHpaByYamlResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateOrUpdateHpaByYamlResponseValidationError) ErrorName() string {
	return "CreateOrUpdateHpaByYamlResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateOrUpdateHpaByYamlResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateOrUpdateHpaByYamlResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateOrUpdateHpaByYamlResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateOrUpdateHpaByYamlResponseValidationError{}

// Validate checks the field values on DeleteHpaRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DeleteHpaRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteHpaRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteHpaRequestMultiError, or nil if none found.
func (m *DeleteHpaRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteHpaRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetClusterName()) < 1 {
		err := DeleteHpaRequestValidationError{
			field:  "ClusterName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := DeleteHpaRequestValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetNamespace()) < 1 {
		err := DeleteHpaRequestValidationError{
			field:  "Namespace",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteHpaRequestMultiError(errors)
	}

	return nil
}

// DeleteHpaRequestMultiError is an error wrapping multiple validation errors
// returned by DeleteHpaRequest.ValidateAll() if the designated constraints
// aren't met.
type DeleteHpaRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteHpaRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteHpaRequestMultiError) AllErrors() []error { return m }

// DeleteHpaRequestValidationError is the validation error returned by
// DeleteHpaRequest.Validate if the designated constraints aren't met.
type DeleteHpaRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteHpaRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteHpaRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteHpaRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteHpaRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteHpaRequestValidationError) ErrorName() string { return "DeleteHpaRequestValidationError" }

// Error satisfies the builtin error interface
func (e DeleteHpaRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteHpaRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteHpaRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteHpaRequestValidationError{}

// Validate checks the field values on DeleteHpaResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DeleteHpaResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteHpaResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteHpaResponseMultiError, or nil if none found.
func (m *DeleteHpaResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteHpaResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return DeleteHpaResponseMultiError(errors)
	}

	return nil
}

// DeleteHpaResponseMultiError is an error wrapping multiple validation errors
// returned by DeleteHpaResponse.ValidateAll() if the designated constraints
// aren't met.
type DeleteHpaResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteHpaResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteHpaResponseMultiError) AllErrors() []error { return m }

// DeleteHpaResponseValidationError is the validation error returned by
// DeleteHpaResponse.Validate if the designated constraints aren't met.
type DeleteHpaResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteHpaResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteHpaResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteHpaResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteHpaResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteHpaResponseValidationError) ErrorName() string {
	return "DeleteHpaResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteHpaResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteHpaResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteHpaResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteHpaResponseValidationError{}

// Validate checks the field values on HpaDetailRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *HpaDetailRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HpaDetailRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// HpaDetailRequestMultiError, or nil if none found.
func (m *HpaDetailRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *HpaDetailRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetClusterName()) < 1 {
		err := HpaDetailRequestValidationError{
			field:  "ClusterName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := HpaDetailRequestValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetNamespace()) < 1 {
		err := HpaDetailRequestValidationError{
			field:  "Namespace",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return HpaDetailRequestMultiError(errors)
	}

	return nil
}

// HpaDetailRequestMultiError is an error wrapping multiple validation errors
// returned by HpaDetailRequest.ValidateAll() if the designated constraints
// aren't met.
type HpaDetailRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HpaDetailRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HpaDetailRequestMultiError) AllErrors() []error { return m }

// HpaDetailRequestValidationError is the validation error returned by
// HpaDetailRequest.Validate if the designated constraints aren't met.
type HpaDetailRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HpaDetailRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HpaDetailRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HpaDetailRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HpaDetailRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HpaDetailRequestValidationError) ErrorName() string { return "HpaDetailRequestValidationError" }

// Error satisfies the builtin error interface
func (e HpaDetailRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHpaDetailRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HpaDetailRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HpaDetailRequestValidationError{}

// Validate checks the field values on HpaDetailResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *HpaDetailResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HpaDetailResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// HpaDetailResponseMultiError, or nil if none found.
func (m *HpaDetailResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *HpaDetailResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetDetail()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, HpaDetailResponseValidationError{
					field:  "Detail",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, HpaDetailResponseValidationError{
					field:  "Detail",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDetail()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HpaDetailResponseValidationError{
				field:  "Detail",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return HpaDetailResponseMultiError(errors)
	}

	return nil
}

// HpaDetailResponseMultiError is an error wrapping multiple validation errors
// returned by HpaDetailResponse.ValidateAll() if the designated constraints
// aren't met.
type HpaDetailResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HpaDetailResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HpaDetailResponseMultiError) AllErrors() []error { return m }

// HpaDetailResponseValidationError is the validation error returned by
// HpaDetailResponse.Validate if the designated constraints aren't met.
type HpaDetailResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HpaDetailResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HpaDetailResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HpaDetailResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HpaDetailResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HpaDetailResponseValidationError) ErrorName() string {
	return "HpaDetailResponseValidationError"
}

// Error satisfies the builtin error interface
func (e HpaDetailResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHpaDetailResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HpaDetailResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HpaDetailResponseValidationError{}
