package elasticsearch

import (
	"fmt"
	"github.com/cch123/elasticsql"
	"log"
)

// 将常规的sql语句翻译成dsl语句方便使用
func GetDslBySqlDesc(sql string) string {
	dsl, esType, err := elasticsql.Convert(sql)
	if err != nil {
		fmt.Printf("err msg is : %v", err)
	}
	log.Printf("dsl is:  %s, esType is:  %s\n", dsl, esType)
	return dsl
}
