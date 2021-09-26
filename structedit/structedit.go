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

func Render(name string, obj interface{}) {
	if imgui.TreeNode(name) {
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
					if f.IsNil() {
						imgui.Text(pad(fields[i], maxlen))
						imgui.SameLine()
						if imgui.Button("<nil>##" + fields[i]) {
							f.Set(reflect.New(t.Field(i).Type.Elem()))
						}
						if imgui.IsItemHovered() {
							imgui.SetTooltip("Instaniate?")
						}
					} else {
						Render(fields[i], f.Interface())
					}
					continue
				} else if k != reflect.Struct {
					imgui.Text(pad(fields[i], maxlen))
					imgui.SameLine()
				}

				a := f.Addr()
				inf := a.Interface()
				switch k {
				case reflect.String:
					imgui.InputText(id(fields[i]), inf.(*string))
				case reflect.Bool:
					imgui.Checkbox(id(fields[i]), inf.(*bool))
				case reflect.Float32:
					imgui.DragFloat(id(fields[i]), inf.(*float32))
				case reflect.Float64:
					t := float32(f.Float())
					imgui.DragFloat(id(fields[i]), &t)
					f.SetFloat(float64(t))
				case reflect.Int32:
					imgui.DragInt(id(fields[i]), inf.(*int32))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
					t := int32(f.Int())
					imgui.DragInt(id(fields[i]), &t)
					f.SetInt(int64(t))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					t := int32(f.Uint())
					imgui.DragInt(id(fields[i]), &t)
					f.SetUint(uint64(t))
				case reflect.Struct:
					Render(fields[i], inf)
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
