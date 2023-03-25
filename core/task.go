package core

import (
	"SuperArch/middleware/redis"
	"SuperArch/middleware/register"
	"SuperArch/middleware/taskcontrol"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
)

type Task struct {
	Category 		string
	Module 			string
	Action 			string
	TaskInfo 		interface{}
}


func (task Task)AddTask(ctx iris.Context)  {
	type Resp struct {
		TaskId		string
		Topic		string
	}
	err := ctx.ReadJSON(&task)
	if err != nil {
		logrus.Errorf("[AddTask][Read Json] %s", err)
		ctx.StatusCode(500)
		return
	}

	topic := fmt.Sprintf("%s.%s.%s", task.Category, task.Module, task.Action)
	taskInfoType := register.SchedulerRegister.ModulePool[topic]
	taskInfo := reflect.New(taskInfoType).Interface()

	taskUuid := taskcontrol.SchedulerTaskControl.InitTaskControl()
	reflect.ValueOf(taskInfo).Elem().FieldByName("TaskId").Set(reflect.ValueOf(taskUuid))

	err = mapstructure.Decode(task.TaskInfo, &taskInfo)
	if err != nil{
		logrus.Errorf("[AddTask][Mapstructure] %s", err)
		ctx.StatusCode(500)
		return
	}

	bodyBytes, _ := json.Marshal(taskInfo)
	register.SchedulerRegister.SendToMQ(topic, bodyBytes)

	resp := Resp{
		TaskId: taskUuid,
		Topic: topic,
	}

	retryAfter := taskcontrol.SchedulerTaskControl.FetchRetryAfter(topic)
	ctx.Header("Retry-After", strconv.FormatInt(int64(retryAfter), 10))
	ctx.JSON(resp)

	ctx.Next()
}

func (task Task)TaskStatus(ctx iris.Context){
	type Req struct {
		TaskId		string
		Topic		string
	}
	type Resp struct {
		TaskStatus  string
		TaskMsg		string
		TaskResult  interface{}
	}

	req := Req{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		logrus.Errorf("[TaskStatus][Read Json] %s", err)
		ctx.JSON(Resp{
			TaskStatus: "fail",
			TaskMsg: "Input Form Error",
			TaskResult: nil,
		})
		return
	}
	if req.Topic == "" || req.TaskId == "" {
		logrus.Errorf("[TaskStatus][Input Param] Missing some params")
		ctx.JSON(Resp{
			TaskStatus: "fail",
			TaskMsg: "Input Form Error: Missing some params",
			TaskResult: nil,
		})
		return
	}

	taskControl := taskcontrol.SchedulerTaskControl.FetchTaskControl(req.TaskId)
	if taskControl == nil{
		logrus.Errorf("[TaskStatus][FetchTaskControl] Error")
		ctx.JSON(Resp{
			TaskStatus: "fail",
			TaskMsg: "Fetch TaskControl Error",
			TaskResult: nil,
		})
		return
	}

	switch taskControl.Status {
	case taskcontrol.TaskFinish:
		result := taskcontrol.SchedulerTaskControl.GetResultFromDB(req.TaskId, req.Topic)
		if result == ""{
			ctx.JSON(Resp{
				TaskStatus: "fail",
				TaskMsg: "Get Task Result Failed",
				TaskResult: nil,
			})
		}else {
			ctx.JSON(Resp{
				TaskStatus: "success",
				TaskMsg: "Get Task Result Successfully",
				TaskResult: result,
			})
			//redis.DeleteHsetByKey(taskcontrol.TaskControlHkey, req.TaskId)
		}

		break
	case taskcontrol.TaskRunning:
		logrus.Infof("Task %s is still running", req.TaskId)
		ctx.Header("Retry-After", "10")
		ctx.JSON(Resp{
			TaskStatus: "running",
			TaskMsg: fmt.Sprintf("Task %s is still running", req.TaskId),
			TaskResult: nil,
		})
		break
	case taskcontrol.TaskWaiting:
		logrus.Infof("Task %s is waiting", req.TaskId)
		retryAfter := taskcontrol.SchedulerTaskControl.FetchRetryAfter(req.Topic)
		noWaitingMsgCnt := register.SchedulerRegister.GetNoWaitingMsgCntByTopic(req.Topic)
		ctx.Header("Retry-After", strconv.FormatInt(int64(retryAfter)*(noWaitingMsgCnt+1), 10))
		ctx.JSON(Resp{
			TaskStatus: "waiting",
			TaskMsg: fmt.Sprintf("Task %s is still waiting", req.TaskId),
			TaskResult: fmt.Sprintf("{\"WaitingListSize\": %d}", noWaitingMsgCnt),
		})
		break
	default:
		ctx.JSON(Resp{
			TaskStatus: "fail",
			TaskMsg: "Get Task Result Failed",
			TaskResult: nil,
		})
		redis.DeleteHsetByKey(taskcontrol.TaskControlHkey, req.TaskId)
	}
	ctx.Next()
}
