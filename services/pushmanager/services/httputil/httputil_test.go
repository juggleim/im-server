package httputil_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	FormField         string            `json:"json_field" form:"form_field"`
	IgnoredField      string            `json:"-"`
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
		FormField:       "form value",
		IgnoredField:    "ignored value",
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
	values := httputil.StructToUrlValues(v)
	if got := values.Get("form_field"); got != "form value" {
		t.Fatalf("form tag was not respected: got %q", got)
	}
	if values.Has("json_field") {
		t.Fatal("json tag should not be used when a form tag is present")
	}
	if values.Has("IgnoredField") || values.Has("-") {
		t.Fatal("ignored field was encoded")
	}
}

func TestPostJSON(t *testing.T) {
	type payload struct {
		Message string `json:"message"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("unexpected content type: %q", got)
		}
		var request payload
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if request.Message != "hello" {
			t.Errorf("unexpected request message: %q", request.Message)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"message":"received"}`))
	}))
	defer server.Close()

	var response payload
	code, _, err := httputil.PostJSON(context.Background(), server.Client(), server.URL, &payload{Message: "hello"}, &response, nil)
	if err != nil {
		t.Fatalf("PostJSON failed: %v", err)
	}
	if code != http.StatusCreated || response.Message != "received" {
		t.Fatalf("unexpected response: code=%d message=%q", code, response.Message)
	}
}

func TestPostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Errorf("unexpected content type: %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		if got := r.Form.Get("message"); got != "hello" {
			t.Errorf("unexpected form value: %q", got)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	var response struct {
		OK bool `json:"ok"`
	}
	code, _, err := httputil.PostForm(context.Background(), server.Client(), server.URL, url.Values{"message": {"hello"}}, &response, nil)
	if err != nil {
		t.Fatalf("PostForm failed: %v", err)
	}
	if code != http.StatusOK || !response.OK {
		t.Fatalf("unexpected response: code=%d ok=%v", code, response.OK)
	}
}
