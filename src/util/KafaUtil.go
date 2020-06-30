package util

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"log"
	"os"
)

type Response struct {
	client   string
	Status   int      `json:"status"`
	Msg      string   `json:"msg"`
	ExtraMsg []string `json:"extra_msg"`
}

// 发送消息
func SendMsg(topic string, content string, clusterAddress []string) Response {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	// jps clusterAddress := []string{"bjlt03-h1510.sy:9092", "bjlt03-h1511.sy:9092", "bjlt03-h1512.sy:9092"}
	client, err := sarama.NewSyncProducer(clusterAddress, config)
	if err != nil {
		log.Printf("producer close error :, %v \n", err)
		return Response{client: "producer", Status: 500, Msg: err.Error()}
	}
	defer client.Close()
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(content)
	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Printf("send msg failed , err : %v\n", err)
	}
	fmt.Println("pid : %v, offset : %v", partition, offset)
	//time.Sleep(time.Second)
	return Response{client: "producer", Status: 200, Msg: "sendSuccess"}
}

// 消费数据
func ConsumerMsg(topic string) Response {
	var broker = "bjlt03-h1510.sy:9092"
	allMsg := make([]string, 0)
	consumer, err := sarama.NewConsumer([]string{broker}, nil)
	if err != nil {
		log.Printf("kafka conected error : %v", err)
	}
	defer consumer.Close()
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Printf("get partition error : %v", err)
	}

	//for _, p := range partitions {
	//	fmt.Printf("partitions p is : %v", p)
	//}
	// 一般申请kafka都会选择多分区，在生产者存储数据的时候会选择分区存储，所以在消费的时候我们需要循环遍历所有分区才能获取到全部元数据
	for _, p := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			log.Printf("get consumer partition failed : %v\n", err)
			continue
		}
		defer partitionConsumer.Close()
		for msg := range partitionConsumer.Messages() {
			allMsg = append(allMsg, string(msg.Value))
			fmt.Printf("msg  is %v, key is %v, offset is %v, \n", string(msg.Value), string(msg.Key), msg.Offset)
		}
	}
	return Response{client: "consumer", Status: 200, Msg: "consumerSuccess", ExtraMsg: allMsg}
}

func ConsumeGroup(topic string) {
	// init config
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	// init consumer
	brokers := []string{"localhost:9092"}
	topics := []string{"test"}
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				_, _ = fmt.Fprintf(os.Stdout, "%v/%v/%v\t%v\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark msg as processed
			}

		}
	}
}
