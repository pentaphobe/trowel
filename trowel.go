package trowel

import "fmt"

type Trowel interface {
	// Get value from underlying array
	Index(idx int) (Trowel, error)
	// Get value from underlying dictionary
	Key(key string) (Trowel, error)

	// Get underlying data
	Get() interface{}
	// Get underlying data as an array
	Array() ([]interface{}, error)
	// Get underlying data as a dictionary
	Map() (map[string]interface{}, error)
}

type trowelWrapper struct {
	data interface{}
}

func (w *trowelWrapper) Get() interface{} {
	return w.data
}

func (w *trowelWrapper) Index(idx int) (Trowel, error) {
	if w.data == nil {
		return &trowelWrapper{}, fmt.Errorf("can't do index lookup (%d) on nil", idx)
	}
	var arr []interface{}
	var ok bool
	if arr, ok = w.data.([]interface{}); !ok {
		return &trowelWrapper{}, fmt.Errorf("can't do index lookup (%d) on non-array %v", idx, w.data)
	}
	return &trowelWrapper{
		data: arr[idx],
	}, nil
}

func (w *trowelWrapper) Key(key string) (Trowel, error) {
	if w.data == nil {
		return &trowelWrapper{}, fmt.Errorf("can't do key lookup (%s) on nil", key)
	}
	var mp map[string]interface{}
	var ok bool
	if mp, ok = w.data.(map[string]interface{}); !ok {
		return &trowelWrapper{}, fmt.Errorf("can't do key lookup (%s) on non-dictionary %v", key, w.data)
	}
	return &trowelWrapper{
		data: mp[key],
	}, nil
}

func (w *trowelWrapper) Array() ([]interface{}, error) {
	if w.data == nil {
		return nil, fmt.Errorf("underlying data is nil")
	}
	var arr []interface{}
	var ok bool
	if arr, ok = w.data.([]interface{}); !ok {
		return nil, fmt.Errorf("underlying data is not an array")
	}
	return arr, nil
}

func (w *trowelWrapper) Map() (map[string]interface{}, error) {
	if w.data == nil {
		return nil, fmt.Errorf("underlying data is nil")
	}
	var dict map[string]interface{}
	var ok bool
	if dict, ok = w.data.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("underlying data is not a map")
	}
	return dict, nil
}

func NewTrowel(data interface{}) Trowel {
	return &trowelWrapper{
		data: data,
	}
}
