package main

import (
	"context"
	"fmt"
	kafka2 "github.com/confluentinc/confluent-kafka-go/kafka"
	"imsdk/internal/common/model/forward"
	"imsdk/internal/common/model/message"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/app"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/mq/kafka"
	"imsdk/pkg/redis"
	"os"
	"sync"
)

var (
	ak = ""
)

func main() {
	log.Start()
	app.LoadConfig()
	ak, _ = config.GetConfigAk()
	ctx := context.Background()
	ctx = log.WithFields(ctx, map[string]string{"action": "consumer"})
	defer func() {
		if err := recover(); err != nil {
			log.Logger().Error(ctx, "panic err: ", funcs.PanicTrace(err))
		}
	}()

	runEnv, err := app.Config().GetChildConf("global", "system", "run_env")
	if err != nil {
		log.Logger().Fatal(ctx, "failed to get system config")
	}

	_ = os.Setenv("RUN_MODULE", "consumer")
	_ = os.Setenv("RUN_ENV", runEnv.(string))
	mongo.Start()
	kafka.Start()
	redis.Start()
	consume(ctx)
	select {}
}

func consume(ctx context.Context) {
	var topics map[string]int
	err := app.Config().Bind("mq", "topics", &topics)
	if err != nil {
		log.Logger().Fatal(ctx, "consume; err: ", err)
	}

	wg := &sync.WaitGroup{}
	//fmt.Println("topics:", topics)
	for topic, amount := range topics {
		for i := 0; i < amount; i++ {
			wg.Add(1)
			go func(msgTopic string) {
				defer wg.Done()
				consumeTopic(ctx, msgTopic)
			}(topic)
		}
	}

	wg.Wait()
}
func consumeTopic(ctx context.Context, topic string) {
	reader, err := kafka.GetConsumeReader(topic)
	run := true
	for run == true {
		ev := reader.Poll(100)
		if ev == nil {
			continue
		}
		run, err = consumer(ctx, ev, topic, reader)
		if err != nil {
			run, err = consumer(ctx, ev, topic, reader)
			if err != nil {
				run, err = consumer(ctx, ev, topic, reader)
				if err != nil {
					reader.CommitMessage(ev.(*kafka.Msg))
				}
			}
		}
	}
	log.Logger().Error(ctx, "consumeTopic", "Closing consumer\n")
	reader.Close()
}

func consumer(ctx context.Context, ev kafka2.Event, topic string, reader *kafka2.Consumer) (bool, error) {
	var err error
	run := true
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Logger().Error(ctx, "consumer panic ,err: ", pErr)
		}
	}()
	switch e := ev.(type) {
	case *kafka.Msg:
		log.Logger().Info(ctx, fmt.Sprintf("%% Msg on %s:\n%s\n", e.TopicPartition, string(e.Value)))
		switch topic {
		case getTopic(base.TopicMsgToWaitReProcess):
			err = message.ProcessMessage(ctx, e.Value)
			break
		case getTopic(base.TopicSocketMessagePush):
			err = forward.PushMessageToSocket(ctx, e.Value)
			break
		case getTopic(base.TopicMsgToOffline):
			err = forward.ProcessOfflineMsg(ctx, e.Value)
			break
		}
		if err == nil {
			reader.CommitMessage(e)
		}
	case kafka.Err:
		log.Logger().Error(ctx, fmt.Sprintf(" Err: %v: %v\n", e.Code(), e), e)
		if e.Code() == kafka.BrokersDown {
			run = false
		}
	}
	return run, err
}

func getTopic(topic string) string {
	res, _ := config.GetConfigTopic(topic)
	//fmt.Println("getTopic Res:", res, topic)
	//return ak + "-" + topic
	return res
}
