package forms

// The errors type holds the validation error messages for forms.
// The name of the form field will be used as the key of the map.
type errors map[string][]string

// Add replaces or add error messages for a given field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves error messages for a given field.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
