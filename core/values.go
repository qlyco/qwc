package core

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/qlyco/qwc/internal/res"
)

type templateData map[string]any

func (values templateData) text(key string) (string, bool) {
	text, isString := values[key].(string)

	return text, isString
}

func (values templateData) date() time.Time {
	date, _ := values[res.DateKey].(time.Time)

	return date
}

func (values templateData) isDisabled(key string) bool {
	text, _ := values.text(key)

	return text == res.DisabledValue
}

func (values templateData) visibility() string {
	text, isString := values.text(res.VisibilityKey)

	if !isString {
		return res.VisibleValue
	}

	if !slices.Contains(res.VisibilityValues, text) {
		return res.NoneValue
	}

	return text
}

func (values templateData) isText() bool {
	text, _ := values.text(res.TypeKey)

	return text == res.TextValue
}

func (values templateData) problems() []string {
	var problems []string

	date, isString := values.text(res.DateKey)

	if isString {
		problems = append(problems, fmt.Sprintf(res.InvalidDateMessage, date))
	}

	pageType, isString := values.text(res.TypeKey)

	if isString && !slices.Contains(res.TypeValues, pageType) {
		problems = append(problems, fmt.Sprintf(res.InvalidTypeMessage, pageType))
	}

	pageSize, isString := values.text(res.PageSizeKey)
	count, err := strconv.Atoi(pageSize)

	if isString && (err != nil || count < 1) {
		problems = append(problems, fmt.Sprintf(res.InvalidPageSizeMessage, pageSize, res.PageSizeDefault))
	}

	visibility, isString := values.text(res.VisibilityKey)

	if isString && !slices.Contains(res.VisibilityValues, visibility) {
		problems = append(problems, fmt.Sprintf(res.InvalidVisibilityMessage, visibility))
	}

	return problems
}
