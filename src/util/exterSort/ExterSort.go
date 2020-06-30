package exterSort

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func TestPipLineLocal() {
	file, err := os.Create("small.in")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	source := RandomSource(6000)
	writer := bufio.NewWriter(file)
	Sink(writer, source)
	writer.Flush()
	// /Users/zhangdanfeng/GoPro/Web/small.in
	wd, _ := os.Getwd()
	fmt.Println(wd)
	file, err = os.Open("small.in")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := ReadSource(bufio.NewReader(file), -1)
	// 全部读取
	count := 0
	for y := range p {
		if count >= 100 {
			break
		}
		fmt.Println(y)
		count++
	}
	fmt.Println("================")
	ExterSort()
}

func MergeDemo() {
	souce := Merge(InMemorySort(ArraySource(3, 2, 6, 4)), InMemorySort(ArraySource(9, 10, 2, 8)))
	for v := range souce {
		fmt.Println(v)
	}
}
func ExterSort() {
	p := createPipline("small.in", 100000, 4)
	writeToFile(p, "small.out")
	printFile("small.out")
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := ReadSource(file, -1)
	for v := range p {
		fmt.Println(v)
	}
}

func writeToFile(ints <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	Sink(writer, ints)
}

func createPipline(filename string, fileSize, chunkCount int) <-chan int {

	chunkSize := fileSize / chunkCount
	Init()
	sortedResults := []<-chan int{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		_, _ = file.Seek(int64(i*chunkSize), 0)
		source := ReadSource(bufio.NewReader(file), chunkSize)
		sortedResults = append(sortedResults, InMemorySort(source))
	}
	return MergeN(sortedResults...)
}

// 分布式并行排序-- 跨主机网络通信
func createNetPipline(filename string, fileSize, chunkCount int) <-chan int {

	chunkSize := fileSize / chunkCount
	Init()
	sortedResultsAddr := []string{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		// Seek 方法用于设置偏移量的，这样可以从某个特定位置开始操作数据流。听起来和 ReaderAt/WriteAt 接口有些类似，
		// 不过 Seeker 接口更灵活，可以更好的控制读写数据流的位置。
		_, _ = file.Seek(int64(i*chunkSize), 0)
		source := ReadSource(bufio.NewReader(file), chunkSize)
		addr := ":" + strconv.Itoa(7000+i)
		NetWorkSink(addr, InMemorySort(source))
		sortedResultsAddr = append(sortedResultsAddr, addr)
	}
	sortedResults := []<-chan int{}
	for _, addr := range sortedResultsAddr {
		source := NetWorkSource(addr)
		sortedResults = append(sortedResults, source)
	}
	return MergeN(sortedResults...)
}
