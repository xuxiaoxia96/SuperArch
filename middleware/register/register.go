package register

import (
	"SuperArch/conf"
	"SuperArch/core/ai/painter"
	"SuperArch/middleware/taskcontrol"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type Register struct {
	ModulePool map[string]reflect.Type
}

var SchedulerRegister Register

func (register *Register)Init()  {
	register.ModulePool = make(map[string]reflect.Type)
	register.ModulePool["ai.painter.txt2image"] = reflect.TypeOf(painter.Txt2Image{})
	register.ModulePool["ai.painter.image2image"] = reflect.TypeOf(painter.Image2Image{})
	register.ModulePool["ai.painter.imagescore"] = reflect.TypeOf(painter.ImageScore{})
	register.ModulePool["ai.painter.image2txt"] = reflect.TypeOf(painter.Image2Txt{})
	register.ModulePool["ai.painter.imagevqa"] = reflect.TypeOf(painter.ImageVQA{})

	taskcontrol.SchedulerTaskControl.InitTaskType(register.ModulePool)
}

func (register *Register)SendToMQ(topic string, message []byte)  {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		conf.Cfg.RabbitMQ.Username, conf.Cfg.RabbitMQ.Password, conf.Cfg.RabbitMQ.Host, conf.Cfg.RabbitMQ.Port))
	if err != nil{
		logrus.Errorf("[SendToMQ][RabbitMQ Dial] %s", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil{
		logrus.Errorf("[SendToMQ][RabbitMQ Channel] %s", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		conf.Cfg.RabbitMQ.Exchange,		// name
		"topic",		// type
		true,		// durable
		false,	// auto-deleted
		false,		// internal
		false,		// no-wait
		nil,			// arguments
	)
	if err != nil{
		logrus.Errorf("[SendToMQ][Declare Exchange] %s", err)
		return
	}

	mqctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(mqctx,
		conf.Cfg.RabbitMQ.Exchange,		// exchange
		topic,	// routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: 	"text/plain",
			Body:     		message   ,
		})
	if err != nil{
		logrus.Errorf("[SendToMQ][Publish MQ MSG] %s", err)
		return
	}

	//logrus.Infof(" [x] Sent %s", bodyBytes)
	logrus.Info(" [x] Sent success")
}

func (register *Register)GetNoWaitingMsgCntByTopic(topic string) int64{
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		conf.Cfg.RabbitMQ.Username, conf.Cfg.RabbitMQ.Password, conf.Cfg.RabbitMQ.Host, conf.Cfg.RabbitMQ.Port))
	if err != nil{
		logrus.Errorf("[GetNoWaitingMsgCntByTopic][RabbitMQ Dial] %s", err)
		return 0
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil{
		logrus.Errorf("[GetNoWaitingMsgCntByTopic][RabbitMQ Channel] %s", err)
		return 0
	}
	defer ch.Close()

	dok, err := ch.QueueDeclare(topic,true, false, false, true, nil)
	if err != nil{
		logrus.Errorf("[GetNoWaitingMsgCntByTopic][QueueDeclare] %s", err)
		return 0
	}
	return int64(dok.Messages)
}

