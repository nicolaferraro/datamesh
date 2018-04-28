package protobuf

import (
	"encoding/json"
	"reflect"
	"errors"
	"strconv"
)

type ObjectData struct {
	Path 	*Path
	Content interface{}
}

func (data *Data) Unmarshal() (*ObjectData, error) {
	if data == nil {
		return nil, nil
	} else {
		var jsonObject interface{}
		if err := json.Unmarshal(data.Content, &jsonObject); err != nil {
			return nil, err
		}

		return &ObjectData{
			Path: &(*data.Path),
			Content: jsonObject,
		}, nil
	}
}

func (data *ObjectData) Expand() ([]*ObjectData, error) {
	if data == nil {
		return make([]*ObjectData, 0), nil
	} else if reflect.ValueOf(data.Content).Kind() != reflect.Map && reflect.ValueOf(data.Content).Kind() != reflect.Slice {
		return []*ObjectData{data}, nil
	} else if slice, ok := data.Content.([]interface{}); ok {
		expanded := make([]*ObjectData, 0)
		for idx,v := range slice {
			subpath := *data.Path
			subpath.Location = subpath.Location + "[" + strconv.Itoa(idx) + "]"

			subdata := ObjectData{
				Path: &subpath,
				Content: v,
			}

			subexpanded, err := subdata.Expand()
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, subexpanded...)
		}

		return expanded, nil
	} else if m, ok := data.Content.(map[string]interface{}); ok {
		expanded := make([]*ObjectData, 0)
		for k,v := range m {
			subpath := *data.Path
			subpath.Location = subpath.Location + "." + k

			subdata := ObjectData{
				Path: &subpath,
				Content: v,
			}

			subexpanded, err := subdata.Expand()
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, subexpanded...)
		}

		return expanded, nil
	} else {
		return nil, errors.New("Unsupported map type. Keys must be strings")
	}
}