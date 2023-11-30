package kafka

//import (
//	"context"
//	"fmt"
//	//confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
//	"imsdk/internal/common/pkg/config"
//	"imsdk/pkg/app"
//	"imsdk/pkg/errno"
//	"imsdk/pkg/log"
//)
//
//type Conf struct {
//	Addr     string `toml:"addr"`
//	Username string `toml:"username"`
//	Password string `toml:"password"`
//}
//
//type Msg = confluentKafka.Message
//type Err = confluentKafka.Error
//
//var (
//	conf Conf
//)
//
//const (
//	BrokersDown = confluentKafka.ErrAllBrokersDown
//)
//
//func Start() {
//	app.Config().Bind("db", "kafka", &conf)
//}
//
//func ProduceMsg(ctx context.Context, topicTag string, msg []byte) error {
//	logCtx := log.WithFields(ctx, map[string]string{"action": "ProduceMsg"})
//	topic, _ := config.GetConfigTopic(topicTag)
//	log.Logger().Info(logCtx, "produce data: ", topicTag, topic, string(msg))
//	configMap := &confluentKafka.ConfigMap{
//		"bootstrap.servers": conf.Addr,
//		//"security.protocol": "SASL_SSL",
//		//"sasl.mechanism":    "AWS_MSK_IAM",
//		//"sasl.jaas.config":  "software.amazon.msk.auth.iam.IAMLoginModule required;",
//		//"sasl.client.callback.handler.class":"software.amazon.msk.auth.iam.IAMClientCallbackHandler",
//	}
//	//if conf.Username != "" && conf.Password != "" {
//	//	//err := configMap.SetKey("sasl.username", conf.Username)
//	//	//err = configMap.SetKey("sasl.password", conf.Password)
//	//	err := configMap.SetKey("security.protocol", "SASL_SSL")
//	//	err = configMap.SetKey("sasl.mechanism", "AWS_MSK_IAM")
//	//	if err != nil {
//	//		log.Logger().Error(logCtx, "failed to add kafka conf, err: ", err)
//	//		return errno.Add("conf error", errno.DefErr)
//	//	}
//	//}
//	p, err := confluentKafka.NewProducer(configMap)
//	if err != nil {
//		return err
//	}
//	defer p.Close()
//	deliveryChan := make(chan confluentKafka.Event)
//	err = p.Produce(&confluentKafka.Message{TopicPartition: confluentKafka.TopicPartition{Topic: &topic, Partition: confluentKafka.PartitionAny}, Value: msg}, deliveryChan)
//	e := <-deliveryChan
//	m := e.(*confluentKafka.Message)
//	if m.TopicPartition.Error != nil {
//		log.Logger().Error(logCtx, "failed to delivery msg, err: ", m.TopicPartition.Error)
//		return errno.Add("delivery failed", errno.DefErr)
//	} else {
//		str := fmt.Sprintf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
//		log.Logger().Info(logCtx, "produce data: ", str)
//	}
//	close(deliveryChan)
//	return nil
//}
//
//func GetConsumeReader(topic string) (*confluentKafka.Consumer, error) {
//	logCtx := log.WithFields(context.Background(), map[string]string{"action": "GetConsumeReader"})
//	configMap := &confluentKafka.ConfigMap{
//		"bootstrap.servers":     conf.Addr,
//		"broker.address.family": "v4",
//		"session.timeout.ms":    6000,
//		"auto.offset.reset":     "earliest",
//		"group.id":              "default",
//	}
//	if conf.Username != "" && conf.Password != "" {
//		//err := configMap.SetKey("sasl.username", conf.Username)
//		//err = configMap.SetKey("sasl.password", conf.Password)
//		//err := configMap.SetKey("security.protocol", "SASL_SSL")
//		//err = configMap.SetKey("sasl.mechanism", "AWS_MSK_IAM")
//		//if err != nil {
//		log.Logger().Error(logCtx, "failed to add kafka conf, err: ", nil)
//		//	return nil, errno.Add("conf error", errno.DefErr)
//		//}
//	}
//	c, err := confluentKafka.NewConsumer(configMap)
//	if err != nil {
//		return nil, err
//	}
//	c.SubscribeTopics([]string{topic}, nil)
//	return c, nil
//}
