package taskcontrol

import (
	"SuperArch/middleware/postgres"
	"SuperArch/middleware/redis"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

const (
	TaskControlHkey = "TaskControl"
	TaskTypeHkey = "TaskType"
)

const (
	TaskWaiting = 0
	TaskRunning = 1
	TaskFinish = 2
	TaskException = 9
)

type TaskControl struct {
	Status 			int64
	InsertTime		time.Time
}

var SchedulerTaskControl TaskControl

type TaskType struct {
	AverCostTime 	float64
	MostCostTime	float64
}

func (taskControl TaskControl)InitTaskControl() string {
	taskControl.Status = 0
	taskControl.InsertTime = time.Now()
	taskControlBytes, _ := json.Marshal(taskControl)
	kv := make(map[string]interface{})
	taskUuid := uuid.NewString()
	kv[taskUuid] = string(taskControlBytes)
	redis.InsertHset(TaskControlHkey, kv)

	return taskUuid
}

func (taskControl TaskControl)InitTaskType(ModulePool map[string]reflect.Type) {
	var taskType TaskType
	taskTypeSettingsExtra := make(map[string]string)
	taskTypeSettings := GetTaskTypeSettings()
	for module,_ := range ModulePool{
		moduleExtra := false
		for storedModule, _ := range taskTypeSettings{
			if storedModule == module{
				moduleExtra = true
				break
			}
		}
		if moduleExtra == true{
			err := json.Unmarshal([]byte(taskTypeSettings[module]), &taskType)
			if err != nil {
				continue
			}
		}else{
			taskType = TaskType{
				AverCostTime: 1,
				MostCostTime: 10,
			}
			tmpTaskTypeSetting,err := json.Marshal(taskType)
			if err != nil {
				continue
			}
			taskTypeSettingsExtra[module] = string(tmpTaskTypeSetting)
		}

		taskTypeBytes,_ := json.Marshal(taskType)
		kv := make(map[string]interface{})
		kv[module] = string(taskTypeBytes)
		redis.InsertHset(TaskTypeHkey, kv)

	}

	// Persistence Extra TaskType Config
	taskControl.SaveTaskTypeToDB(taskTypeSettingsExtra)
}

func GetTaskTypeSettings() map[string]string{
	// Get All TaskType Settings From DB
	taskTypeSettings := make(map[string]string)
	pdb := postgres.GetPostgresClient()
	if pdb == nil{
		return taskTypeSettings
	}
	sql := `SELECT "module","settings" from "TaskType"`
	rows, err := pdb.Query(sql)
	if err != nil {
		logrus.Errorf("[TaskControl Module][SaveResultToDB][Get Result] %s", err)
		return taskTypeSettings
	}
	for rows.Next(){
		tmpModule := ""
		tmpSettings := ""
		err = rows.Scan(&tmpModule, &tmpSettings)
		if err != nil {
			continue
		}
		taskTypeSettings[tmpModule] = tmpSettings
	}

	pdb.Close()

	return taskTypeSettings
}

func (taskControl TaskControl) SaveTaskTypeToDB(tasTypeSettings map[string]string){
	pdb := postgres.GetPostgresClient()
	if pdb == nil{
		return
	}
	sql := `INSERT INTO "TaskType" ("module", settings) VALUES ($1, $2)`
	txn, _ := pdb.Begin()
	stmt, err := txn.Prepare(sql)
	defer stmt.Close()
	for module,settings := range tasTypeSettings{
		_, err = stmt.Exec(module, settings)
		if err != nil{
			continue
		}
	}

	err = txn.Commit()
	if err != nil {
		logrus.Errorf("[TaskControl Module][SaveTaskTypeToDB][Commit SQL] %s", err)
		return
	}
	pdb.Close()
	logrus.Infof("[TaskControl Module][SaveTaskTypeToDB] Insert Successfully! %s", tasTypeSettings)
}

func (taskControl TaskControl) UpdateTaskTypeToDB(tasTypeSettings map[string]string){
	pdb := postgres.GetPostgresClient()
	if pdb == nil{
		return
	}
	sql := `UPDATE "TaskType" SET settings=$1 WHERE module=$2`
	txn, err := pdb.Begin()
	stmt, err := txn.Prepare(sql)
	defer stmt.Close()
	for module,settings := range tasTypeSettings{
		_, err = stmt.Exec(settings, module)
		if err != nil{
			continue
		}
	}

	err = txn.Commit()
	if err != nil {
		logrus.Errorf("[TaskControl Module][UpdateTaskTypeToDB][Commit SQL] %s", err)
		return
	}
	pdb.Close()
	logrus.Infof("[TaskControl Module][UpdateTaskTypeToDB] Update Successfully! %s", tasTypeSettings)
}

func (taskControl *TaskControl) UpdateTaskControlStatus(taskId string, status int64) {
	taskControl.Status = status
	newTaskControlStr, _ := json.Marshal(taskControl)
	newTaskControl := make(map[string]interface{})
	newTaskControl[taskId] = string(newTaskControlStr)
	redis.UpdateHset(TaskControlHkey, newTaskControl)
}

func (taskControl TaskControl)FetchRetryAfter(module string) float64 {
	taskTypeStr := redis.GetValFromHsetBykey(TaskTypeHkey, module)
	taskType := TaskType{}
	err := json.Unmarshal([]byte(taskTypeStr), &taskType)
	if err != nil {
		return 0
	}

	return taskType.AverCostTime
}

func (taskControl TaskControl)FetchTaskControl(taskId string) *TaskControl {
	taskStatus := redis.GetValFromHsetBykey(TaskControlHkey, taskId)
	if taskStatus == "" {
		return nil
	}

	err := json.Unmarshal([]byte(taskStatus), &taskControl)
	if err != nil {
		return nil
	}

	return &taskControl
}

func (taskControl TaskControl) GetResultFromDB(requestId, module string) string{
	result := ""
	pdb := postgres.GetPostgresClient()
	if pdb == nil{
		return result
	}
	sql := `SELECT result from "TaskResult" WHERE request_id=$1 AND module=$2`
	err := pdb.QueryRow(sql, requestId, module).Scan(&result)
	if err != nil {
		logrus.Errorf("[TaskControl Module][SaveResultToDB][Get Result] %s", err)
		return result
	}
	logrus.Infof("[TaskControl Module][GetResultFromDB] Get Result Successfully! requestId: %s, module: %s", requestId, module)
	pdb.Close()

	return result
}

type TaskModuleCommon struct {
	TaskId		string
}

