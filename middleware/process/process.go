package process

import (
	"SuperArch/middleware/redis"
	"SuperArch/middleware/taskcontrol"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"time"
)

func CronProcess(){
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("00:30").Do(PersistenceTaskType)
	s.Every(1).Day().At("01:00").Do(CleanExtremeTimeoutTask)

	s.StartAsync()
}

func PersistenceTaskType(){
	kv := redis.GetALLFromHset(taskcontrol.TaskTypeHkey)
	taskcontrol.SchedulerTaskControl.UpdateTaskTypeToDB(kv)
}

func CleanExtremeTimeoutTask(){
	// Clean the task which run more than 1 day
	kv := redis.GetALLFromHset(taskcontrol.TaskControlHkey)
	for k, v := range kv{
		taskCtl := taskcontrol.TaskControl{}
		err := json.Unmarshal([]byte(v), &taskCtl)
		if err != nil {
			continue
		}
		if taskCtl.InsertTime.Before(time.Now().AddDate(0, 0, -1)){
			redis.DeleteHsetByKey(taskcontrol.TaskControlHkey, k)
		}
	}
}