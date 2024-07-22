package pubsub

import (
	"context"
	"strings"
	"time"

	"github.com/jhq0113/gochat/lib/pogo"

	"github.com/Allenxuxu/toolkit/convert"
	"github.com/apache/rocketmq-clients/golang"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type CmdMsg struct {
	Channel string     `json:"c"`
	TraceId string     `json:"tid"`
	AppId   uint8      `json:"aid"`
	Id      int64      `json:"id"`
	Param   pogo.Param `json:"p"`
}

func (cm *CmdMsg) Marshal() []byte {
	data, _ := json.Marshal(cm)
	return data
}

func SubscribeRedis(red *redis.Client, ctx context.Context, channels ...string) <-chan CmdMsg {
	var (
		ps  = red.Subscribe(ctx, channels...)
		ch  = make(chan CmdMsg, 1)
		msg *redis.Message
		err error
	)

	go func() {
		defer ps.Close()

		for {
			msg, err = ps.ReceiveMessage(ctx)
			if err != nil {
				if strings.HasPrefix(err.Error(), "redis: unknown message") {
					continue
				}
				close(ch)
				return
			}

			var cmdMsg CmdMsg
			if err = json.Unmarshal(convert.StringToBytes(msg.Payload), &cmdMsg); err != nil {
				continue
			}

			cmdMsg.Channel = msg.Channel
			ch <- cmdMsg
		}
	}()

	return ch
}

func SubscribeKafkaConfig(configMap kafka.ConfigMap, cb kafka.RebalanceCb, channels ...string) (<-chan CmdMsg, error) {
	c, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, err
	}

	return SubscribeKafka(c, cb, channels...)
}

func SubscribeKafka(c *kafka.Consumer, cb kafka.RebalanceCb, channels ...string) (<-chan CmdMsg, error) {
	var (
		ch  = make(chan CmdMsg, 1)
		err = c.SubscribeTopics(channels, cb)
	)

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			ev, ok := <-c.Events()
			if !ok {
				return
			}

			switch e := ev.(type) {
			case kafka.AssignedPartitions:
			case kafka.RevokedPartitions:
			case kafka.PartitionEOF:
			case kafka.Error:
			case *kafka.Message:
				var cmdMsg CmdMsg
				if err = json.Unmarshal(e.Value, &cmdMsg); err != nil {
					continue
				}

				cmdMsg.Channel = *e.TopicPartition.Topic
				ch <- cmdMsg
			}
		}
	}()

	return ch, nil
}

func SubscribeRocket(c golang.SimpleConsumer, maxMessageNum int32, invisibleDuration time.Duration) (<-chan CmdMsg, error) {
	var (
		ch  = make(chan CmdMsg, 1)
		err = c.Start()
	)

	if err != nil {
		return nil, err
	}

	defer c.GracefulStop()

	go func() {
		for {
			mvs, _ := c.Receive(context.TODO(), maxMessageNum, invisibleDuration)
			if len(mvs) > 0 {
				for _, mv := range mvs {
					if er := c.Ack(context.TODO(), mv); er == nil {
						var cmdMsg CmdMsg
						if err = json.Unmarshal(mv.GetBody(), &cmdMsg); err != nil {
							continue
						}

						cmdMsg.Channel = mv.GetTopic()
						ch <- cmdMsg
					}
				}
			}
		}
	}()

	return ch, nil
}
