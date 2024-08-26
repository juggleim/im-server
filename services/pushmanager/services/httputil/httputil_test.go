package httputil_test

import (
	"context"
	"net/http"
	"testing"

	"im-server/services/pushmanager/services/httputil"
)

type TestType1 struct {
	FieldInt          int               `json:"field_int,omitempty"`
	FieldUint16       uint16            `json:"field_uint_16"`
	FieldUint32       uint32            `json:"field_uint_32"`
	FieldFloat32      float32           `json:"field_float_32"`
	FieldFloat64      float64           `json:"field_float_64"`
	FieldBytes        []byte            `json:"field_bytes"`
	FieldString       string            `json:"field_string"`
	FieldWithoutTag   string            ``
	FieldStringMap    map[string]string `json:"field_string_map"`
	FieldStringMapNil map[string]string `json:"field_string_map_nil"`
	TestType2         *TestType2        `json:"test_type_2"`
	TestType2Nil      *TestType2        `json:"test_type_2_nil"`
}
type TestType2 struct {
	FieldInt       int               `json:"field_int"`
	FieldStringMap map[string]string `json:"field_string_map"`
	TestType3      *TestType3        `json:"test_type_3"`
	TestType3Nil   *TestType2        `json:"test_type_3_nil"`
}

type TestType3 struct {
	FieldString    string            `json:"field_string"`
	FieldStringMap map[string]string `json:"field_string_map"`
	UnhandleType   chan struct{}     `json:"unhandle_type"`
}

func TestStructToUrlValues(t *testing.T) {
	v := &TestType1{
		FieldInt:        0,
		FieldUint16:     16,
		FieldUint32:     32,
		FieldFloat32:    32,
		FieldFloat64:    64,
		FieldBytes:      []byte("test bytes"),
		FieldString:     "test string",
		FieldWithoutTag: "field without tag",
		FieldStringMap: map[string]string{
			"key1": "value1",
		},
		TestType2: &TestType2{
			FieldInt: 1,
			FieldStringMap: map[string]string{
				"key2": "value2",
			},
			TestType3: &TestType3{
				FieldString: "test3 string",
				FieldStringMap: map[string]string{
					"key3": "value3",
				},
			},
		},
	}
	_ = httputil.StructToUrlValues(v)
}

func TestPostJSON(t *testing.T) {
	_, _, _ = httputil.PostJSON(context.Background(), &http.Client{}, "https://ipinfo.io/", &struct{}{}, &struct{}{}, nil)
}

func TestPostForm(t *testing.T) {
	_, _, _ = httputil.PostForm(context.Background(), &http.Client{}, "https://ipinfo.io/", nil, &struct{}{}, nil)
}
