package common

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type variable struct {
	// Output dump string
	Out string

	// Indent counter
	indent int64
}

// refer: https://github.com/liudng/godump/blob/master/dump.go
func (v *variable) dump(val reflect.Value, name string) {
	v.indent++

	if val.IsValid() && val.CanInterface() {
		typ := val.Type()

		switch typ.Kind() {
		case reflect.Array, reflect.Slice:
			v.printType(name, val.Interface())
			l := val.Len()
			for i := 0; i < l; i++ {
				v.dump(val.Index(i), strconv.Itoa(i))
			}
		case reflect.Map:
			v.printType(name, val.Interface())
			//l := val.Len()
			keys := val.MapKeys()
			for _, k := range keys {
				v.dump(val.MapIndex(k), k.Interface().(string))
			}
		case reflect.Ptr:
			v.printType(name, val.Interface())
			v.dump(val.Elem(), name)
		case reflect.Struct:
			v.printType(name, val.Interface())
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				v.dump(val.FieldByIndex([]int{i}), field.Name)
			}
		default:
			v.printValue(name, val.Interface())
		}
	} else {
		v.printValue(name, "")
	}

	v.indent--
}

func (v *variable) printType(name string, vv interface{}) {
	v.printIndent()
	v.Out = fmt.Sprintf("%s%s(%T)\n", v.Out, name, vv)
}

func (v *variable) printValue(name string, vv interface{}) {
	v.printIndent()
	v.Out = fmt.Sprintf("%s%s(%T) %#v\n", v.Out, name, vv, vv)
}

func (v *variable) printIndent() {
	var i int64
	for i = 0; i < v.indent; i++ {
		v.Out = fmt.Sprintf("%s    ", v.Out)
	}
}

// Print to standard out the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Dump(v interface{}) {
	val := reflect.ValueOf(v)
	dump := &variable{indent: -1}
	dump.dump(val, "")
	fmt.Printf("%s", dump.Out)
}

// Return the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Sdump(v interface{}) string {
	val := reflect.ValueOf(v)
	dump := &variable{indent: -1}
	dump.dump(val, "")
	return dump.Out
}

//Wrap os.Getwd()
func Getwd() string {
	w, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//replace backslash(\) to slash(/) on windows platform
	if runtime.GOOS == "windows" {
		w = strings.Replace(w, "\\", "/", -1)
	}

	return w
}

func Getenv(key string) string {
	value := os.Getenv(key)
	if runtime.GOOS == "windows" {
		value = strings.Replace(value, "\\", "/", -1)
	}

	return value
}

func Setenv(key, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	return nil
}

func Chdir(path string) error {
	err := os.Chdir(path)
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	return nil
}

func ParseText(txt string) []string {
	//support space in path
	args := make([]string, 0)
	node := ""
	colon := false
	for _, c := range txt {
		t := string(c)
		//log.Printf("%#v\n", t)
		if t == "'" {
			if colon == true {
				colon = false
			} else {
				colon = true
			}

			continue
		}

		if (t == " " && colon == true) || t != " " {
			node += t
		} else {
			args = append(args, node)
			node = ""
		}
	}

	if node != "" {
		args = append(args, node)
	}

	//args := strings.Split(txt, " ")

	return args
}

// Console parameters
func Arguments(app string) (string, string) {
	var c, h, p string
	flag.StringVar(&c, "c", "./example.json", "Usage: mplus -c=/path/to/example.json")
	flag.StringVar(&h, "h", "nil", "Usage: example -h")
	flag.StringVar(&p, "p", "", "Usage: example -p=Passport/User/Login&id=1")
	flag.Parse()

	if h != "nil" {
		fmt.Println(
			fmt.Sprintf("Usage: %s [OPTION]...", app),
			"{example} is the name of the application, you can change in a real environment.",
			"",
			"  -c  The path of the configuration file.",
			"  -h  Display this help and exit.",
			"  -p  Console application action path. Separated by a slash.")
		os.Exit(0)
	}

	return c, p
}
