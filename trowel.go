package trowel

import (
	"regexp"
	"strconv"
	"strings"
)

type Trowel interface {
	// Get value from underlying array
	Index(idx int) Trowel
	// Get value from underlying dictionary
	Key(key string) Trowel
	// Get by a path string
	Path(path string) Trowel

	// Get underlying data
	Get() interface{}
	// Get underlying data as an array
	Array() ([]interface{}, error)
	// Get underlying data as a dictionary
	Map() (map[string]interface{}, error)

	// Get underlying errors
	HasErrors() bool
	Errors() []error
	Error() error
	addError(...error) Trowel
}

type trowelWrapper struct {
	data   interface{}
	errors []error
}

func (w *trowelWrapper) Get() interface{} {
	return w.data
}

func (w *trowelWrapper) Index(idx int) Trowel {
	if w.data == nil {
		w.addError(NewIndexError(w, "can't do index lookup (%d) on nil", idx))
		return w
	}
	var arr []interface{}
	var ok bool
	if arr, ok = w.data.([]interface{}); !ok {
		return NewTrowel(nil).addError(NewIndexError(w, "can't do index lookup (%d) on non-array %v", idx, w.data))
	}
	return &trowelWrapper{
		data:   arr[idx],
		errors: w.errors,
	}
}

func (w *trowelWrapper) Key(key string) Trowel {
	if w.data == nil {
		w.addError(NewKeyError(w, "can't do key lookup (%s) on nil", key))
		return w
	}
	var mp map[string]interface{}
	var ok bool
	if mp, ok = w.data.(map[string]interface{}); !ok {
		return NewTrowel(nil).addError(NewKeyError(w, "can't do key lookup (%s) on non-dictionary %v", key, w.data))
	}
	var value interface{}
	if value, ok = mp[key]; !ok {
		return NewTrowel(nil).addError(NewKeyError(w, "no key matching '%s'", key))
	}
	return &trowelWrapper{
		data:   value,
		errors: w.errors,
	}
}

func parsePath(pathStr string) ([]interface{}, error) {
	// Dot is optional in regex so we can match invalid keys
	const reRegularKey = `[.]?[\w-_]+`
	const reQuotedKey = `[.]?"[^"]+"`
	const reArrayIndex = `\[[^]]+\]`
	// Build full regex for splitting
	re := regexp.MustCompile(
		`(` +
			strings.Join([]string{
				reRegularKey,
				reQuotedKey,
				reArrayIndex,
			}, "|") +
			`)`)

	split := re.FindAllString(pathStr, -1)
	result := make([]interface{}, 0)
	for _, v := range split {
		switch v[0] {
		case '.':
			sanitised := strings.ReplaceAll(v[1:], "\"", "")
			result = append(result, sanitised)
			continue
		case '[':
			indexStr := v[1 : len(v)-1]
			idx, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, NewParseError("invalid index '%s' at '%s'", indexStr, v)
			}
			result = append(result, idx)
			continue
		default:
			return nil, NewParseError("invalid path element '%s'", v)
		}
	}
	return result, nil
}

func (w *trowelWrapper) Path(path string) Trowel {
	pathComponents, err := parsePath(path)
	if err != nil {
		return NewTrowel(nil).addError(err)
	}
	var child Trowel = w
	for _, component := range pathComponents {
		switch v := component.(type) {
		case int:
			child = child.Index(v)
		case string:
			child = child.Key(v)
		}
		if child.HasErrors() {
			return NewTrowel(nil).addError(child.Errors()...)
		}
	}
	return child
}

func (w *trowelWrapper) Array() ([]interface{}, error) {
	if w.data == nil {
		return nil, NewError(w, "underlying data is nil")
	}
	var arr []interface{}
	var ok bool
	if arr, ok = w.data.([]interface{}); !ok {
		return nil, NewError(w, "underlying data is not an array")
	}
	return arr, nil
}

func (w *trowelWrapper) Map() (map[string]interface{}, error) {
	if w.data == nil {
		return nil, NewError(w, "underlying data is nil")
	}
	var dict map[string]interface{}
	var ok bool
	if dict, ok = w.data.(map[string]interface{}); !ok {
		return nil, NewError(w, "underlying data is not a map")
	}
	return dict, nil
}

func (w *trowelWrapper) HasErrors() bool {
	return len(w.errors) > 0
}
func (w *trowelWrapper) Errors() []error {
	return w.errors
}
func (w *trowelWrapper) Error() error {
	if w.HasErrors() {
		return w.errors[len(w.errors)-1]
	}
	return nil
}
func (w *trowelWrapper) addError(err ...error) Trowel {
	if w.errors == nil {
		w.errors = []error{}
	}
	w.errors = append(w.errors, err...)

	return w
}

func NewTrowel(data interface{}) Trowel {
	return &trowelWrapper{
		data:   data,
		errors: make([]error, 0),
	}
}
