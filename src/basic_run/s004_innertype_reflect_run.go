package basic_run

import (
	"fmt"
	"math/rand"
	"reflect"
)

// 验证slice 的append 和copy 内置函数相关
// 拷贝的长度=min(len(src), len(dst))
func AppendCopyRun() {
	fmt.Printf("Begin: =================ApppendCopyRun()=========================\n")
	srcSlice := []int{1, 2}
	for i := 0; i < 20; i++ {
		srcSlice = append(srcSlice, rand.Intn(20))
		fmt.Printf("Apppend Round %d len=%d cap=%d addr=%p\n", i, len(srcSlice), cap(srcSlice), &srcSlice)
	}
	dstSlice := make([]int, 10)
	copiedNum := copy(dstSlice, srcSlice)
	fmt.Printf("srcLen=%d, src:=%v\n", len(srcSlice), srcSlice)
	// 拷贝的长度=min(len(src), len(dst))
	fmt.Printf("dstLen=%d, dst:=%v\n", len(dstSlice), dstSlice)
	fmt.Printf("copiedNum=%d\n", copiedNum)
	fmt.Printf("End: =================ApppendCopyRun()=========================\n")
}

type Palne struct {
	Engine  string
	Company string
	Price   int64
}

func (p Palne) Fly(postion string) {
	fmt.Printf("Palne fly %+v @ %s\n", p, postion)
}

// 通过反射机制 来获取 interface{} 的具体值和type
func ReflectRun() {
	airbusA380 := Palne{"F5", "airbus", 144000002}
	airbusA380_ptr := &airbusA380

	fmt.Printf("airbusA380=%+v\n", *airbusA380_ptr)

	var iface interface{}
	iface = airbusA380
	// 此时i
	iface_val := reflect.ValueOf(iface) // 返回Value 类型 实际上是对应底层对象的指针 需要使用 Value.Interface 转换成interface{} 然后再强制转换成对应的对象
	iface_type := reflect.TypeOf(iface) // 返回Type 类型 可以调用 Type.Name() 方法转换成名称的字符串
	fmt.Printf("type of i_type=%s\n", iface_type)
	fmt.Printf("value of i_val=%+v\n", iface_val)

	//获取每个成员变量的值 并显示
	for i := 0; i < iface_type.NumField(); i++ {
		switch iface_type.Field(i).Name {
		case "int64":
			fmt.Printf("case int64: Field[%d] value=%+v\n", i, iface_val.Field(i).Interface().(int64))
		case "string":
			fmt.Printf("case string: Field[%d] value=%+v\n", i, iface_val.Field(i).Interface().(string))
		}
	}

	//获取对应类的方法 并调用
	//主要注意 *Plane 和 Plane 是两个不同的类型
	for i := 0; i < iface_type.NumMethod(); i++ {
		refMethod := iface_type.Method(i)
		fmt.Printf("get method %s %v\n", refMethod.Name, refMethod.Type)
		methodValue := iface_val.MethodByName(refMethod.Name)
		args := []reflect.Value{reflect.ValueOf("sky")}
		methodValue.Call(args)
	}
}
