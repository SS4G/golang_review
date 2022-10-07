package basic_run

import (
	"encoding/json"
	"fmt"
)

// 接口 动物
type Animal interface {
	Eat() string
	AverageLifeYear() int
	Serialize() string
}

type Hamster struct {
	Color    string `json:"color"`
	NickName string `json:"nick-name"`
	Age      int    `json:"age"`
}

func (p *Hamster) Eat() string {
	return fmt.Sprintf("hamster %s eat ant", p.NickName)
}

func (p *Hamster) AverageLifeYear() int {
	return 3
}

func (p *Hamster) Serialize() string {
	json_info, err := json.Marshal(p)
	if err != nil {
		return "{\"err\": true}"
	} else {
		return "{\"Hamster\":" + string(json_info) + "}"
	}
}

// 基类 People
type Human struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender bool   `json:"gender"`
}

func (p *Human) Eat() string {
	return fmt.Sprintf("human %s eat rice by mouth use knife", p.Name)
}

func (p *Human) AverageLifeYear() int {
	return 70
}

func (p *Human) Serialize() string {
	json_info, err := json.Marshal(p)
	if err != nil {
		return "{\"err\": true}"
	} else {
		return "{\"Human\":" + string(json_info) + "}"
	}
}

// People 派生类 Student
// 使用了匿名派生相当于继承 所以 BankClerk 也实现了 Animal 接口
// 如果不使用匿名字段就不是继承 也就不能继续使用 Animal 接口
type Student struct {
	Schoole string
	Human
}

// Student 定制了自己的序列化
func (p *Student) Serialize() string {
	var json_map map[string]interface{}
	ser_bytes, _ := json.Marshal(p.Human)
	//ser_bytes, _ := json.Marshal(p.hu)

	err := json.Unmarshal(ser_bytes, &json_map)
	if err != nil {
		fmt.Println(err)
		return "{\"err\": true}"
	}
	json_map["schoole"] = p.Schoole
	json_map["job"] = "student"
	json_info, err := json.Marshal(json_map)
	if err != nil {
		return "{\"err\": true}"
	} else {
		return "{\"Human\":" + string(json_info) + "}"
	}
}

// People 派生类 BankClerk
// 使用了匿名派生相当于继承 所以Student也实现了 Animal 接口
type BankClerk struct {
	Bank   string `json:"bank"`
	Salary int    `json:"salary"`
	Human         //`json:"human_core"` 如果注释掉这一句 Human 内部的 key 会被展平
}

// BankClerk 直接嵌套 序列化
func (p *BankClerk) Serialize() string {
	json_info, err := json.Marshal(p)
	if err != nil {
		return "{\"err\": true}"
	} else {
		return string(json_info)
	}
}

// 测试新建结体
func StructRun() {
	fmt.Printf("Begin: =================StructRun()=========================\n")
	student := &Student{"PKU", Human{Name: "HongkaiWang", Age: 30, Gender: true}}
	bankclerk := &BankClerk{"ICBC", 30000, Human{Name: "Xinyan", Age: 31, Gender: true}}
	hamster := &Hamster{NickName: "Laohuang", Color: "yellow", Age: 1}

	fmt.Println(student.Serialize())
	fmt.Println(bankclerk.Serialize())
	fmt.Println(hamster.Serialize())
	// 直接定义接口变量 实际上并没与实例化对象
	var animal1, animal2, animal3 Animal
	// 可以直接将实例化后的对象赋值给他们实现的接口
	animal1 = student
	animal2 = bankclerk
	animal3 = hamster
	fmt.Printf("Animal: %s\n", animal1.Serialize())
	fmt.Printf("Animal: %s\n", animal2.Serialize())
	fmt.Printf("Animal: %s\n", animal3.Serialize())

	// 直接通过 %+v 输出字段名以及
	fmt.Printf("Student: %+v\n", *student)
	bankclerk_bytes := []byte(animal2.Serialize())
	bank_clk_obj := BankClerk{}
	fmt.Printf("ser bytes=%v\n", bankclerk_bytes)
	err := json.Unmarshal(bankclerk_bytes, &bank_clk_obj)
	if err == nil {
		fmt.Printf("deserailized bank clerk:=%+v\n", bank_clk_obj)
	}
	fmt.Printf("End: =================StructRun()=========================\n")
}
