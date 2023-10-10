package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
)

type EnvConfig struct {
	env  string
	data *interface{}
}

var envConfig = &EnvConfig{}

// load configuration file
func LoadEnvConfig(file string, mode string) error {
	var _file string

	if mode != "" {
		_file = strings.Replace(file, "yaml", mode+".yaml", 1)
	}

	data, failed := loadFile(file)

	dataSrc, failedSrc := loadFile(_file)

	// 如果两个文件都加载成功，则进行merge操作，
	// 如果只有一个文件加载成功，则使用加载成功的文件
	if data != nil && dataSrc != nil && failed == nil && failedSrc == nil {
		envConfig.SetEnv(mode)
		envConfig.SetData(data)
		envConfig.MergeData(dataSrc)
	} else if data != nil && failed == nil {
		envConfig.SetEnv(mode)
		envConfig.SetData(data)
	} else if dataSrc != nil && failedSrc == nil {
		envConfig.SetEnv(mode)
		envConfig.SetData(dataSrc)
	} else {
		return fmt.Errorf("file [ %s ] load failed", file)
	}

	return nil
}

func loadFile(file string) (*interface{}, error) {
	if err := hasFile(file); err != nil {
		return nil, err
	}
	data, failed := ioutil.ReadFile(file)
	if failed != nil {
		return nil, failed
	}

	var _data = new(interface{})

	failed = yaml.Unmarshal(data, _data)
	if failed != nil {
		return nil, failed
	}

	return _data, nil
}

func (p *EnvConfig) SetEnv(mode string) {
	p.env = mode
}

func (p *EnvConfig) SetData(data any) {
	p.data = data.(*interface{})
}

func (p *EnvConfig) MergeData(data any) interface{} {
	srcData := *p.data
	distData := (data.(*interface{}))

	retData := mergeKeyForData(&srcData, distData)
	envConfig.SetData(&retData)

	return &retData
}

func (p *EnvConfig) GetRootValue() interface{} {
	return *p.data
}

// Get the value from root
func GetRootValue() interface{} {
	return envConfig.GetRootValue()
}

// Get the value from root for map
func GetRootMapValue() map[interface{}]interface{} {
	data := envConfig.GetRootValue()
	if data == nil {
		return nil
	}
	return data.(map[interface{}]interface{})
}

// Get the value by path like below
// file.yaml
// test: test
// test2:
//      test3: test333
/// GetValue("test2.test3") ===> test333
func GetValue(path string) (interface{}, error) {
	val, err := getValueByPath(path)
	return val, err
}

func getValueByPath(path string) (interface{}, error) {
	words := strings.Split(path, ".")
	data := GetRootMapValue()
	if data == nil {
		return nil, fmt.Errorf("No key [ %s ] found!", path)
	}

	len := len(words)
	var val interface{}
	var hasKey bool
	for idx, work := range words {
		val, hasKey = getValueByKey(work, data)

		if !hasKey {
			return nil, fmt.Errorf("No key [ %s ] found!", path)
		}

		if idx < len-1 {
			data = val.(map[interface{}]interface{})
		} else {
			return val, nil
		}
	}
	return nil, nil
}

func getValueByKey(key string, data map[interface{}]interface{}) (interface{}, bool) {
	keys := strings.FieldsFunc(key, func(r rune) bool {
		return string(r) == "["
	})
	var val interface{}
	var hasKey bool
	sliceLen := len(keys)

	if sliceLen > 1 {
		for i, k := range keys {
			if strings.HasSuffix(k, "]") {
				keys[i] = strings.Replace(k, "]", "", 1)
			}
		}

		val, hasKey = data[keys[0]]

		if !hasKey {
			return nil, false
		}

		_val, isOk := val.([]interface{})

		if !isOk {
			return nil, false
		}

		for i := 1; i < sliceLen; i++ {
			idx, err := strconv.Atoi(keys[i])
			if err != nil || len(_val) <= idx {
				return nil, false
			}

			newVal := _val[idx]

			if i < sliceLen-1 && isArray(newVal) {
				_val = newVal.([]interface{})
			} else if i < sliceLen-1 {
				return nil, false
			} else {
				return newVal, true
			}
		}
	}

	val, hasKey = data[keys[0]]

	return val, hasKey
}

func isArray(param any) bool {
	val := reflect.ValueOf(param)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		return true
	}
	return false
}

func kind(param any) reflect.Kind {
	kind := reflect.ValueOf(param)
	return kind.Kind()
}

func hasFile(file string) error {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func mergeKeyForData(a *interface{}, b *interface{}) interface{} {
	aa := (*a).(map[interface{}]interface{})
	bb := (*b).(map[interface{}]interface{})

	for key, val := range bb {
		newVal, hasKey := aa[key]
		// 如果不存在的key或者类型不同的key直接复制过去
		if !hasKey || kind(val) != kind(newVal) || kind(val) != reflect.Map {
			aa[key] = val
		} else {
			aa[key] = mergeKeyForData(&newVal, &val)
		}
	}
	return aa
}
