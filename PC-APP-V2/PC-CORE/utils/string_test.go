package utils

import (
    "testing"
    "reflect"
)

func TestRandomKeyGeneration(t *testing.T) {
    for i := 0; i < 100; i++ {
        var (
            str1 string = NewRandomString(32)
            str2 string = NewRandomString(32)
            empty string = NewRandomString(0)
        )
        if len(empty) != 0 {
            t.Errorf("Empty string should be empty")
        }
        t.Log("string 1 - NewRandomString : " + str1)
        t.Log("string 2 - NewRandomString : " + str2)
        if len(str1) == 0 || len(str2) == 0 {
            t.Error("Empty random string!")
        }
        if reflect.DeepEqual(str1, str2) {
            t.Error("Random strings are not different enough")
        }
    }
}

