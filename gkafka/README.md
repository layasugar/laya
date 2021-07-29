## go操作kafka

#### 引入包

```
import "github.com/layasugar/laya/gkafka"
```

#### 初始化

```
	var kc = &gkafka.KafkaConfig{
		Brokers:      Brokers,
		Topic:        Topic,
		Group:        Group,
		User:         User,
		Pwd:          Pwd,
		CertFile:     CertFile,
		KeyFile:      KeyFile,
		CaFile:       CaFile,
		KafkaVersion: KafkaVersion,
		VerifySsl:    VerifySsl,
	}
	Kafka = gkafka.InitKafka(kc)
```

#### consume 消费者 produce 生产者

```
    var DataChan = make(chan []byte)
    
    // 先激活消费
	go work()

	// 再初始化消费者
	go dao.Kafka.InitConsumer(DataChan)

	// 生产数据
	partition, offset, err := dao.Kafka.SendMsg("layatest", "asdddddddasdadasd")
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Printf("Message partion: %d, Message offset: %d.", partition, offset)
	}
	
	func work() {
	    for item := range DataChan {
		    log.Printf(string(item))
		}
	}
}
```