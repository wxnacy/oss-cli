package goss

import (
    "reflect"
)

func InArray(val interface{}, array interface{}) (index int) {
    // 元素是否在数组中
    index = -1

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                index = i
                return
            }
        }
    }

    return
}
