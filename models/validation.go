package models

type ValidationErrors map[string][]string

func (errs ValidationErrors) Add(fieldName string, errorMsg string) {
	errs[fieldName] = append(errs[fieldName], errorMsg)
}