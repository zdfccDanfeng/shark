package exterSort

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var startTime time.Time

func Init() {
	startTime = time.Now()
}

// 从数组里面读取数据作为数据源
func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out) // 表示数据送完
	}()
	return out
}

// 只进不出，返回只出不进的chan
func InMemorySort(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// read into memory
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read in memory done !. ", time.Now().Sub(startTime))
		sort.Ints(a)
		fmt.Println(" memory sort done !. ", time.Now().Sub(startTime))

		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// 两个有序channel上进行归并排序操作
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// 	等待两个上游数据排序完成
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		fmt.Println(" memory merge done !. ", time.Now().Sub(startTime))

		close(out)
	}()
	return out
}

// 从文件读取数据， 最多读取chunkSize字节大小个
func ReadSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024) // chan 增加	buffer
	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			// 最多读取chunkSize个
			if err != nil || (chunkSize != -1 && bytesRead > chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

// 数据输出
func Sink(writer io.Writer, out <-chan int) {
	for v := range out {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		_, _ = writer.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}
