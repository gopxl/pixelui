package structedit

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/inkyblackness/imgui-go"
)

func pad(s string, l int) string {
	if len(s) == l {
		return s
	}
	return s + strings.Repeat(" ", l-len(s))
}

func id(s string) string {
	return "##" + s
}

func Inspect(name string, obj interface{}) {
	if imgui.Begin("Inspect: " + name) {
		Render("Struct", obj)
	}
	imgui.End()
}

func Render(name string, obj interface{}) {
	render(name, obj)
}

func editString(field string, inf interface{}) {
	imgui.InputText(id(field), inf.(*string))
}

func editBool(field string, inf interface{}) {
	imgui.Checkbox(id(field), inf.(*bool))
}

func editFloat(field string, f reflect.Value) {
	switch f.Kind() {
	case reflect.Float32:
		imgui.DragFloat(id(field), f.Addr().Interface().(*float32))
	case reflect.Float64:
		t := float32(f.Float())
		if imgui.DragFloat(id(field), &t) {
			f.SetFloat(float64(t))
		}
	}
}

func editInt(field string, f reflect.Value) {
	switch f.Kind() {
	case reflect.Int32:
		imgui.DragInt(id(field), f.Addr().Interface().(*int32))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		t := int32(f.Uint())
		imgui.DragInt(id(field), &t)
		f.SetUint(uint64(t))
	default:
		t := int32(f.Int())
		imgui.DragInt(id(field), &t)
		f.SetInt(int64(t))
	}
}

func editPtr(field string, t reflect.Type, i, padlen int, f reflect.Value) {
	if f.IsNil() {
		imgui.Text(pad(field, padlen))
		imgui.SameLine()
		if imgui.Button("<nil>##" + field) {
			f.Set(reflect.New(t.Field(i).Type.Elem()))
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Instaniate?")
		}
	} else {
		Render(field, f.Interface())
	}
}

func editArray(field string, t reflect.Type, i int, f reflect.Value) {
	if imgui.TreeNodeV(field, imgui.TreeNodeFlagsDefaultOpen) {

		if imgui.Button("+##" + field) {
			f.Set(reflect.Append(f, reflect.New(t.Field(i).Type.Elem()).Elem()))
		}
		imgui.TreePop()
	}
}

func render(name string, obj interface{}) {
	var tree = false
	if name != "" {
		tree = imgui.TreeNodeV(name, imgui.TreeNodeFlagsDefaultOpen)
	}
	if tree {
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Ptr {
			v := reflect.ValueOf(obj).Elem()
			t = v.Type()
			maxlen := 0
			num := v.NumField()
			fields := make([]string, num)

			for i := 0; i < num; i++ {
				f := t.Field(i).Name
				maxlen = int(math.Max(float64(maxlen), float64(len(f))))
				fields[i] = f
			}

			for i := 0; i < num; i++ {
				f := v.Field(i)
				k := f.Kind()

				if k == reflect.Ptr {
					editPtr(fields[i], t, i, maxlen, f)
					continue
				} else if !(k == reflect.Struct || k == reflect.Array || k == reflect.Slice) {
					imgui.Text(pad(fields[i], maxlen))
					imgui.SameLine()
				}

				a := f.Addr()
				inf := a.Interface()
				switch k {
				case reflect.String:
					editString(fields[i], inf)
				case reflect.Bool:
					editBool(fields[i], inf)
				case reflect.Float32, reflect.Float64:
					editFloat(fields[i], f)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					editInt(fields[i], f)
				case reflect.Struct:
					Render(fields[i], inf)
				case reflect.Array, reflect.Slice:
					editArray(fields[i], t, i, f)
				default:
					imgui.Text(fmt.Sprintf("<unrenderable type: %s>", fields[i]))
				}
			}
		} else {
			v := reflect.ValueOf(obj)
			maxlen := 0
			num := v.NumField()
			fields := make([]string, num)

			for i := 0; i < num; i++ {
				f := t.Field(i).Name
				maxlen = int(math.Max(float64(maxlen), float64(len(f))))
				fields[i] = f
			}

			for i := 0; i < num; i++ {
				f := v.Field(i)
				k := f.Kind()
				if k == reflect.Ptr && !f.IsNil() {
					Render(fields[i], f.Elem().Interface())
				} else if k == reflect.Struct {
					Render(fields[i], f.Interface())
				} else {
					imgui.Text(pad(fields[i], maxlen))
					imgui.SameLine()
					imgui.Text(fmt.Sprint(f))
				}
			}
		}
		imgui.TreePop()
	}
}
