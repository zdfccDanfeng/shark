package lsm

import (
	"fmt"
	"github.com/robaho/keydb"
)

// LSM Tree (Log-structured merge-tree)
// LSM的优势是能够高效率地写入数据，这种效率来源于充分利用磁盘和内存这两种存储介质的特点。
// 将大量的随机写入落到内存中，再批量地将这些随机写合并成连续写落地磁盘
// LSM是一种应付大数据量写入磁盘的数据结构模型，在NoSQL系统中非常常见，尤其是应对写多读少的场景非常有效.

//   写性能 LSM分层组织数据，随机写先落地内存，之后的 内存刷盘 和 合并操作 都是顺序io，保证了写性能的高效。
//   读性能 读取的数据可能落在两个部分：内存 和 磁盘。
//       对于内存，高效读的数据结构选择比较多，针对不同的场景有很多的选择，本项目使用的是 AVL树 。
//       对于磁盘，高效的数据结构往往都来自于索引和并发io。在本项目中，使用了 key/data文件（索引）、分段segment（并发）、快照读来（并发）技术。

func TestKeyDb() {
	// 从指定目录打开一个数据库实例，所有的数据最终会写入到此目录中
	db, err := keydb.Open("/Users/zhangdanfeng/GoPro/src/github.com/shark/templates/keydb", true)
	if err != nil {
		panic(err)
	}
	// 开启一个操作gogo表的事务，注意在keydb中单个事务只能读写同一个表
	// 单个表的数据最终会落地在以表名开头的一系列文件中
	tx, err := db.BeginTX("gogo")
	if err != nil {
		panic(err)
	}
	// 往此事务中写入一个数据，此时数据仍未持久化到磁盘，是内存操作，速度快
	err = tx.Put([]byte("k1"), []byte("v1"))
	if err != nil {
		panic(err)
	}
	// 按key读取数据，未命中内存的时候才会去磁盘中读取
	value, err := tx.Get([]byte("k1"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("value:%v", value)
	// 此次事务的所有修改一起落盘
	if err = tx.CommitSync(); err != nil {
		panic(err)
	}
	// 关闭数据库，并显示调用合并方法
	if err = db.CloseWithMerge(1); err != nil {
		panic(err)
	}
}
