package handling

import (
	"encoding/json"
	"reflect"
)

type BaseResponse[T any] struct {
	Message string   `json:"message"`
	Result  bool     `json:"result"`
	Failed  []string `json:"failed,omitempty"`
	Token   string   `json:"token,omitempty"`
	Data    T        `json:"data,omitempty"`
	Code    int      `json:"-"`
}

func (r BaseResponse[T]) MarshalJSON() ([]byte, error) {
	type Alias BaseResponse[T]
	raw := struct {
		Alias
		Data any `json:"data,omitempty"`
	}{
		Alias: Alias(r),
	}

	// reflection check
	rv := reflect.ValueOf(r.Data)
	kind := rv.Kind()

	switch kind {
	case reflect.Slice, reflect.Array:
		// SELALU kirim slice/array (meskipun kosong)
		raw.Data = r.Data
	default:
		// If data is nil pointer or nil interface → skip
		if !rv.IsValid() {
			break
		}

		if (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface) && rv.IsNil() {
			break
		}

		zero := reflect.Zero(rv.Type()).Interface()

		if !reflect.DeepEqual(r.Data, zero) {
			raw.Data = r.Data
		}
	}

	return json.Marshal(raw)
}
