package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func PostJSON(ctx context.Context, client *http.Client, url string, req, res interface{}, headers map[string]string) (code int, body string, err error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, "", err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return 0, "", err
	}
	request.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	if client == nil {
		client = &http.Client{}
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}

	if err = json.Unmarshal(resBody, res); err != nil {
		return response.StatusCode, string(resBody), err
	}
	return response.StatusCode, string(resBody), nil
}

func PostForm(ctx context.Context, client *http.Client, url string, req url.Values, res interface{}, headers map[string]string) (code int, body string, err error) {
	reqBody := req.Encode()
	request, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(reqBody))
	if err != nil {
		return 0, "", err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	if client == nil {
		client = &http.Client{}
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}

	if err = json.Unmarshal(resBody, res); err != nil {
		return response.StatusCode, string(resBody), err
	}
	return response.StatusCode, string(resBody), nil
}

func StructToUrlValues(intf interface{}) (values url.Values) {
	values = url.Values{}
	if intf == nil || reflect.ValueOf(intf).Kind() != reflect.Ptr {
		return values
	}
	if reflect.ValueOf(intf).IsNil() {
		return values
	}

	iValue := reflect.ValueOf(intf).Elem()
	iType := iValue.Type()

	for i := 0; i < iValue.NumField(); i++ {
		fieldValue := iValue.Field(i)
		structField := iType.Field(i)

		tag := structField.Tag.Get("json")
		name, opts := parseTag(tag)
		if !isValidTag(name) {
			name = ""
		}
		if name == "" {
			name = structField.Name
		}

		var v string
		switch field := fieldValue.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(fieldValue.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(fieldValue.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(fieldValue.Float(), 'f', 4, 64)
		case []byte:
			v = string(fieldValue.Bytes())
		case string:
			v = fieldValue.String()
		case map[string]string:
			for k, v := range field {
				values.Set(fmt.Sprintf("%s.%s", name, k), v)
			}
		default:
			subValues := StructToUrlValues(field)
			for subKey, subValue := range subValues {
				for _, v := range subValue {
					values.Add(fmt.Sprintf("%s.%s", name, subKey), v)
				}
			}
		}
		if (v == "" || v == "0") && opts.Contains("omitempty") {
			continue
		}
		values.Set(name, v)
	}
	return values
}
