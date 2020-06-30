package util

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// 读取小文件，可以使用ioUtil包进行
// 大文件 按行进行读取
func readFileToList(fileName string) []string {
	res := make([]string, 0)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
		return res
	}
	defer file.Close()
	// 缓冲 buffer
	buf := bufio.NewReader(file)
	for {
		line, _, err := buf.ReadLine()
		line_content := strings.TrimSpace(string(line))
		// handle(line)
		if err != nil {
			if err == io.EOF {
				//fmt.Println(err.Error())
				log.Println(err.Error())
				break
			}
		}
		res = append(res, line_content)
	}
	return res
}

// 大文件分片进行读取
// 第二个方案就是分片处理，当读取的是二进制文件，没有换行符的时候，使用下面的方案一样处理大文件
func ReadBigFile(fileName string, handle func([]byte)) error {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("can't opened this file")
		return err
	}
	defer f.Close()
	s := make([]byte, 4096)
	for {
		switch nr, err := f.Read(s[:]); true {
		case nr < 0:
			fmt.Fprintf(os.Stderr, "cat: error reading: %s\n", err.Error())
			os.Exit(1)
		case nr == 0: // EOF
			return nil
		case nr > 0:
			handle(s[0:nr])
		}
	}
	return nil
}

func TestReadBigFile(fileName string) []string {
	return readFileToList(fileName)
}
