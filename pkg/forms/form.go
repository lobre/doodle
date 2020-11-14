package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX is a regular expression for sanity checking the format
// of an email address. This pattern is the one currently recommended
// by the W3C and Web Hypertext Application Technology Working Group.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form will validate form data against a particular set of rules.
// If an error occurs, it will store an error message associated with
// the field. Only a few validation rules are implemented so far.
// To implement more rules, have a look at the following blog post:
// https://www.alexedwards.net/blog/validation-snippets-for-go
type Form struct {
	url.Values
	Errors errors
}

// New creates a new Form taking data as entry.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks that specific fields in the form
// data are present and not blank. If any fields fail this check,
// add the appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength checks that a specific field in the form contains
// a minimum number of characters. If the check fails, then add
// the appropriate message to the form errors.
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// MaxLength checks that a specific field in the form contains
// a maximum number of characters. If the check fails, then add
// the appropriate message to the form errors.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	// check proper characters instead of bytes
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// PermittedValues checks that a specific field in the form matches
// one of a set of specific permitted values. If the check fails,
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// MatchesPattern checks that a specific field in the form matches
// a regular expression. If the check fails, then add the appropriate
// message to the form errors.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// Valid returns true if there are no errors in the form.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
