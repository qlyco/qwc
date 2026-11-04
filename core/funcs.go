package core

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/qlyco/qwc/internal/res"
)

func dict(pairs ...any) (templateData, error) {
	if len(pairs)%res.DictPairSize != 0 {
		return nil, errors.New(res.DictPairsMessage)
	}

	result := templateData{}

	for i := 0; i < len(pairs); i += res.DictPairSize {
		key, isString := pairs[i].(string)

		if !isString {
			return nil, fmt.Errorf(res.DictKeyMessage, pairs[i])
		}

		result[key] = pairs[i+1]
	}

	return result, nil
}

func slice(items any, indexes ...int) (any, error) {
	value := reflect.ValueOf(items)

	if value.Kind() != reflect.Slice && value.Kind() != reflect.String {
		return nil, fmt.Errorf(res.SliceTypeMessage, items)
	}

	if len(indexes) > res.SliceIndexLimit {
		return nil, errors.New(res.SliceIndexesMessage)
	}

	bounds := []int{0, value.Len()}

	for i, index := range indexes {
		bounds[i] = min(max(index, 0), value.Len())
	}

	return value.Slice(min(bounds[0], bounds[1]), bounds[1]).Interface(), nil
}
