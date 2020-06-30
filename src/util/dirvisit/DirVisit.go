package dirvisit

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			// is file
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

// ioutil.ReadDir函数会返回一个os.FileInfo类型的slice，os.FileInfo类型也是os.Stat这个函数的返回值。对每一个子目录而言，
// walkDir会递归地调用其自身，同时也在递归里获取每一个文件的信息。walkDir函数会向fileSizes这个channel发送一条消息。
// 这条消息包含了文件的字节大小。
//
//下面的主函数，用了两个goroutine。后台的goroutine调用walkDir来遍历命令行给出的
// 每一个路径并最终关闭fileSizes这个channel。主goroutine会对其从channel中接收到的文件大小进行累加，并输出其和。
var verbose = flag.Bool("v", false, "show verbose progress messages")

func DoVisitDirConcurrent() {
	// Determine the initial directories.

	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	// Traverse the file tree.
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	//go func() {
	//	for _, root := range roots {
	//		walkDir(root, fileSizes)
	//	}
	//	close(fileSizes)
	//}()
	//
	//// Print the results.
	//var nfiles, nbytes int64
	//for size := range fileSizes {
	//	nfiles++
	//	nbytes += size
	//}
	var nfiles, nbytes int64
	for _, root := range roots {
		n.Add(1)
		go walkDir2(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// sema is a counting semaphore for limiting concurrency in dirents.
// 由于这个程序在高峰期会创建成百上千的goroutine，我们需要修改dirents函数，
// 用计数信号量来阻止他同时打开太多的文件
var sema = make(chan struct{}, 20)

func walkDir2(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	defer n.Done()            // 释放许可
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir2(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}
