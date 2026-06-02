package jutil

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type KVPair struct {
	Key   string
	Value interface{}
}

type m map[string]interface{}

type OrderedMap2 struct {
	m
	l    *list.List
	keys map[string]*list.Element // the double linked list for delete and lookup to be O(1)
}

// Create a new OrderedMap
func NewOrderedMap() *OrderedMap2 {
	return &OrderedMap2{
		m:    make(map[string]interface{}),
		l:    list.New(),
		keys: make(map[string]*list.Element),
	}
}

// Create a new OrderedMap and populate from a list of key-value pairs
func NewOrderedMapFromKVPairs(pairs []*KVPair) *OrderedMap2 {
	om := NewOrderedMap()
	for _, pair := range pairs {
		om.Set(pair.Key, pair.Value)
	}
	return om
}

// return all keys
// func (om *OrderedMap) Keys() []string { return om.keys }

// set value for particular key, this will remember the order of keys inserted
// but if the key already exists, the order is not updated.
func (om *OrderedMap2) Set(key string, value interface{}) {
	if _, ok := om.m[key]; !ok {
		om.keys[key] = om.l.PushBack(key)
	}
	om.m[key] = value
}

// Check if value exists
func (om *OrderedMap2) Has(key string) bool {
	_, ok := om.m[key]
	return ok
}

// Get value for particular key, or nil if not exist; but don't rely on nil for non-exist; should check by Has or GetValue
func (om *OrderedMap2) Get(key string) interface{} {
	return om.m[key]
}

// Get value and exists together
func (om *OrderedMap2) GetValue(key string) (value interface{}, ok bool) {
	value, ok = om.m[key]
	return
}

// deletes the element with the specified key (m[key]) from the map. If there is no such element, this is a no-op.
func (om *OrderedMap2) Delete(key string) (value interface{}, ok bool) {
	value, ok = om.m[key]
	if ok {
		om.l.Remove(om.keys[key])
		delete(om.keys, key)
		delete(om.m, key)
	}
	return
}

// Iterate all key/value pairs in the same order of object constructed
func (om *OrderedMap2) EntriesIter() func() (*KVPair, bool) {
	e := om.l.Front()
	return func() (*KVPair, bool) {
		if e != nil {
			key := e.Value.(string)
			e = e.Next()
			return &KVPair{key, om.m[key]}, true
		}
		return nil, false
	}
}

// Iterate all key/value pairs in the reverse order of object constructed
func (om *OrderedMap2) EntriesReverseIter() func() (*KVPair, bool) {
	e := om.l.Back()
	return func() (*KVPair, bool) {
		if e != nil {
			key := e.Value.(string)
			e = e.Prev()
			return &KVPair{key, om.m[key]}, true
		}
		return nil, false
	}
}

// this implements type json.Marshaler interface, so can be called in json.Marshal(om)
// func (om *OrderedMap2) MarshalJSON() (res []byte, err error) {
func (om *OrderedMap2) MarshalJSON() (res []byte, err error) {

	res = append(res, '{')
	front, back := om.l.Front(), om.l.Back()

	for e := front; e != nil; e = e.Next() {

		k := e.Value.(string)
		res = append(res, fmt.Sprintf("%q:", k)...)
		var b []byte
		b, err = json.Marshal(om.m[k])
		// fmt.Println("1:", string(b))

		// 20230105, adly add this, prevent convert to unicode when string
		switch v := om.m[k].(type) {
		case string:

			str := fmt.Sprintf("%v", v)

			// b = []byte(strconv.QuoteToASCII(string(str)))
			b = []byte(strconv.QuoteToGraphic(string(str)))

			// fmt.Println("2:", string(b))
		case []interface{}:

			new := strings.Replace(string(b), "\\u0026", "&", -1)
			new = strings.Replace(new, "\\u003c", "<", -1)
			new = strings.Replace(new, "\\u003e", ">", -1)

			// fmt.Println("3:", new)

			b = []byte(new)

		default:
			// And here I'm feeling dumb. ;)
			// fmt.Printf("I don't know, ask stackoverflow.")
		}
		// end of 20230105, adly add this, prevent convert to unicode when string

		if err != nil {
			return
		}

		res = append(res, b...)

		if e != back {
			res = append(res, ',')
		}

	}
	res = append(res, '}')
	// fmt.Println(string(res))
	// bs = string(res)
	// fmt.Printf("marshalled: %v: %#v\n", res, res)
	return
}

// this implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, om)
func (om *OrderedMap2) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	err = om.parseobject(dec)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("expect end of JSON object but got more token: %T: %v or err: %v", t, t, err)
	}

	return nil
}

func (om *OrderedMap2) parseobject(dec *json.Decoder) (err error) {
	var t json.Token
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return err
		}

		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("expecting JSON key should be always a string: %T: %v", t, t)
		}

		t, err = dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var value interface{}
		value, err = handledelim(t, dec)
		if err != nil {
			return err
		}

		// om.keys = append(om.keys, key)
		om.keys[key] = om.l.PushBack(key)
		om.m[key] = value
	}

	t, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	return nil
}

func parsearray(dec *json.Decoder) (arr []interface{}, err error) {
	var t json.Token
	arr = make([]interface{}, 0)
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return
		}

		var value interface{}
		value, err = handledelim(t, dec)
		if err != nil {
			return
		}
		arr = append(arr, value)
	}
	t, err = dec.Token()
	if err != nil {
		return
	}
	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		err = fmt.Errorf("expect JSON array close with ']'")
		return
	}

	return
}

func handledelim(t json.Token, dec *json.Decoder) (res interface{}, err error) {
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			om2 := NewOrderedMap()
			err = om2.parseobject(dec)
			if err != nil {
				return
			}
			return om2, nil
		case '[':
			var value []interface{}
			value, err = parsearray(dec)
			if err != nil {
				return
			}
			return value, nil
		default:
			return nil, fmt.Errorf("unexpected delimiter: %q", delim)
		}
	}
	return t, nil
}
