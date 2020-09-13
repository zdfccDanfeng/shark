package bitmap

import "fmt"

//位图
type BitMap struct {
	bits []byte
	//这样写的目的是便于扩展 一些其它参数，如最大值、最小值、长度等
	//max int
	//min int
	//len int
	// ......
}

// 1亿里面找到重复的两个数
// 位图，1个字节存储（映射）8个字，1亿个字大概12M内存就够了。比较适合。
func find2sameBetween1yi() {
	//我们模拟产生1 - 1亿，来做实验
	maxnum := 100000000 //设置最大值
	//模拟产生1亿个数
	bif := make([]uint, maxnum)
	for i := 0; i < maxnum; i++ {
		bif[i] = uint(i)
	}
	//假设一个重复数
	bif[9999998] = 1007
	//创建最大值是1以的位图
	bit := NewBitMap(maxnum)
	for _, v := range bif {
		//判断是否存在
		if bit.IsExist(v) {
			fmt.Println("找到重复数：", v)
			return
		}
		//将不存在的数，添加到BitMap里
		bit.Add(v)
	}
}

//初始化一个BitMap
//一个byte有8位,代表8个数字,取余后加1便是存放最大数所需的容量
func NewBitMap(max int) *BitMap {
	bits := make([]byte, (max>>3)+1)
	return &BitMap{bits: bits}
}

//添加一个数字到BitMap
//计算添加数字在数组中的索引index,一个索引可以存放8个数字
//计算存放到索引下的第几个位置,一共0-7个位置
//原索引下的内容与1左移到指定位置后做或运算
func (b *BitMap) Add(num uint) {
	index := num >> 3
	pos := num & 0x07
	b.bits[index] |= 1 << pos
}

//判断一个数字是否存在
//找到数字所在的位置,然后做与运算
func (b *BitMap) IsExist(num uint) bool {
	index := num >> 3
	pos := num & 0x07
	return b.bits[index]&(1<<pos) != 0
}
