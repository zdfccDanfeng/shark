package reflect

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
	Age  int
	Id   int
}

func (u User) GetName(name string) {
	fmt.Println("我的名字叫: ", u.Name, " 传进来的是 ", name)
}

func TestReflectUser() {
	a := User{"xiaoming", 5, 6}
	a.GetName("测试1")
	Setname(&a)
	a.GetName("测试2")
	fmt.Println(a.Name)
	info(a)
	ra := reflect.ValueOf(a)
	rm := ra.MethodByName("GetName")
	args := []reflect.Value{reflect.ValueOf("测试三")} //通过反射调用该对象方法
	//这里再次回到切片与数组。到底什么是切片，什么是数组
	//再次总结一下，刚看了点资料,以下是链接
	//<https://www.zhihu.com/question/66673454>
	//以下是总结
	//1.数组要指明长度Array := [ArrayLength] ElementType，arr :=[4]int,数组拷贝为值拷贝
	//2.切片不需要指明长度SliceT := []ElementType,但是他有容量机制，拷贝为指针拷贝
	arr := [4]int{0, 1, 3, 4}
	arrs := arr[1:2]
	ccc := reflect.TypeOf(arrs).Kind() //通过反射判断是一个数组还是切片
	fmt.Println(ccc)
	rm.Call(args)
}
func info(o interface{}) { //获取这个实例化了的结构体的所有信息
	ot := reflect.TypeOf(o)  //通过反射获取类型
	ov := reflect.ValueOf(o) //获取值
	fmt.Println(ot, "\n")
	fmt.Println(ov, "\n")
	fmt.Println(ot.Name(), "\n")            //打印这个实现了空接口的结构名称
	for i := 0; i < ot.NumField(); i += 1 { //获取这个结构体内的值
		n := ot.Field(i)
		val := ov.Field(i)
		fmt.Println(n.Type, ":", n.Name, ",", val)
	}
	for i := 0; i < ot.NumMethod(); i += 1 { //打印这个结构体的每个方法
		f := ot.Method(i)
		fmt.Printf("%s:%v,%v\n", f.Name, f.Type, f.Index)
	}
}

func Setname(o interface{}) { //通过反射设置值
	v := reflect.ValueOf(o)
	fmt.Println(v.Elem())          //这里打印出来是这个结构体
	fmt.Println(v.Elem().CanSet()) //这里返回的是个布尔
	fmt.Println(v.Kind())          //这里打印出来的就是ptr
	fmt.Println(reflect.Ptr)       //也是ptr,引出个问题，为什么不用字符串"ptr"?以后再找答案,指针
	fmt.Println(v)                 //Elem返回v持有的接口保管的值的Value封装，
	//或者v持有的指针指向的值的Value封装。如果v的Kind不是Interface或Ptr会panic；如果v持有的值为nil，会返回Value零值。
	//Kind返回v持有的值的分类，如果v是Value零值，返回值为Invalid
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() { //如果是指针且不能被修改
		fmt.Println("xxx") //这个判断有问题,如果不是指针也会走else,或者可以修改但不是指针也会走else
		return
	} else {
		v = v.Elem()
	}
	n := v.FieldByName("Name")
	if n.IsValid() && n.Kind() == reflect.String { //如果取到了值并且类型是字符串则修改
		fmt.Println(n, "\n", n.Kind())
		n.SetString("lele") //修改名字//通过反射对值进行修改，肯定可以修改拉，传进来的是指针
	}
}
