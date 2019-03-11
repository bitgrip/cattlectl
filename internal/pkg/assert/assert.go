// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package assert

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	yaml "gopkg.in/yaml.v2"
)

var Update = flag.Bool("update", false, "update .golden files")
var typeOfBytes = reflect.TypeOf([]byte(nil))

func AssertGoldenFile(tb testing.TB, testName string, actual []byte) {
	golden := fmt.Sprintf("testdata/golden-files/%s.golden", testName)
	if *Update {
		ioutil.WriteFile(golden, actual, 0644)
	}
	expected, err := ioutil.ReadFile(golden)
	Ok(tb, err)
	Equals(tb, string(expected), string(actual))
}

// assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// not ok fails the test if an err is nil.
func NotOk(tb testing.TB, err error, expecteError string) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: missing error\033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
	}
	if "" != expecteError {
		Equals(tb, expecteError, err.Error())
	}
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)

		fmt.Printf(
			"\033[31m%s:%d:\ndiff:\033[39m\n\n\t%s\n\n",
			filepath.Base(file),
			line,
			strings.Replace(buildDiff(exp, act), "\n", "\n\t", -1),
		)
		tb.FailNow()
	}
}

func buildDiff(exp, act interface{}) string {
	fmt.Println(reflect.TypeOf(exp))
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(buildCompareString(act), buildCompareString(exp), true)
	return dmp.DiffPrettyText(diffs)
}
func buildCompareString(input interface{}) string {
	result := buildCompareStringForType(input)
	if result == "" {
		return fmt.Sprint(input)
	}
	return result
}

func buildCompareStringForType(input interface{}) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Struct:
		fmt.Println(input)
		result, err := yaml.Marshal(input)
		if err != nil {
			return fmt.Sprint(input)
		}
		return string(result)
	case reflect.Ptr:
		if v.IsNil() {
			return "nil ptr"
		}
		return buildCompareString(v.Elem())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Slice:
		if v.Type() == typeOfBytes {
			return buildCompareString(string(input.([]byte)))
		}
		result, err := yaml.Marshal(input)
		if err != nil {
			return fmt.Sprint(input)
		}
		return string(result)
	case reflect.String:
		return input.(string)
	}
	return fmt.Sprint(input)
}

func FailInStub(tb testing.TB, stubStacks int, msg string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(stubStacks + 1)
	fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	tb.FailNow()
}

func Fail(tb testing.TB, msg string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	tb.FailNow()
}
