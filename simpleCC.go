package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	// "github.com/hyperledger/fabric-chaincode-go/shim"
	// pb "github.com/hyperledger/fabric-protos-go/peer"
	// "github.com/golang/protobuf/proto"
	// "github.com/hyperledger/fabric-protos-go/msp"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Simple struct {
}

type User struct {
	Name           string         `json:"name"`
	Id             string         `json:"id"`
	Account        int            `json:"account_balance"`
	Reputation     int            `json:"reputation"`
	Info           string         `json:"detail_information"`
	Skills         []string       `json:"skills"`
	Profession     []string       `json:"profession"`
	PostTasks      []string       `json:"post_tasks"`
	OngoingTasks   []string       `json:"ongoing_tasks"`
	CompleteTasks  []string       `json:"complete_tasks"`
	Otherplatforms map[string]int `json:"other_platforms_workyears"`
}

type Task struct {
	Name          string            `json:"name"`
	Id            string            `json:"id"`
	Type          string            `json:"task_type"` //1.competition 2.one2one 3.private
	Detail        string            `json:"detail"`
	Reward        int               `json:"reward"`
	State         string            `json:"state"`
	ReceiveTime   time.Time         `json:"receive_time"`
	Deadline      time.Time         `json:"deadline"`
	RequesterId   string            `json:"requester"`
	Candidate     map[string]string `json:"candidate_worker_and_solution"`
	FinalWorker   string            `json:"final_worker"`
	FinalSolution string            `json:"final_solution"`
	Requirement   []string          `json:"worker_requirement"` //reputation, complete num, skill, profession
}

type WorkRecord struct {
	TaskId    string    `json:"task_id"`
	Requester string    `json:"requester"`
	Worker    string    `json:"worker"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
}

func (t *Simple) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Simple Init")
	return shim.Success(nil)
}

func (t *Simple) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Simple Invoke")
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "userRegister":
		return userRegister(stub, args)
	case "userAddSkill":
		return userAddSkill(stub, args)
	case "userAddPro":
		return userAddPro(stub, args)
	case "userImport":
		return userImport(stub, args)
	case "taskPost":
		return taskPost(stub, args)
	case "taskReceive":
		return taskReceive(stub, args)
	case "taskCommit":
		return taskCommit(stub, args)
	case "rewardAllocate":
		return rewardAllocate(stub, args)
	case "recharge":
		return recharge(stub, args)
	case "userQuery":
		return userQuery(stub, args)
	case "taskQuery":
		return taskQuery(stub, args)
	case "recordQuery":
		return recordQuery(stub, args)
	case "alluserQuery":
		return alluserQuery(stub, args)
	case "alltaskQuery":
		return alltaskQuery(stub, args)
	case "taskUpdate":
		return taskUpdate(stub, args)
	default:
		return shim.Error("Invalid invoke function name.")
	}
}

//func constructUserKey(userId string) string {
//	return fmt.Sprintf("user_%s", userId)
//}

//func constructTaskKey(taskId string) string {
//	return fmt.Sprintf("task_%s", taskId)
//}

//func constructRecordKey(taskId string, requester string, worker string) string {
//	return fmt.Sprintf("record_%s_%s_%s", taskId, requester, worker)
//}

//用户输入name account 创建一个账户 参数：Name info
func userRegister(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//检查参数个数 Name info
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//验证参数的正确性
	name := args[0]
	info := args[1]
	if name == "" {
		return shim.Error("Invalid args")
	}
	//获取当前用户证书
	creatorByte, err := stub.GetCreator() //获取的是msp里面的signcerts
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	id := string(si.GetIdBytes())
	//验证数据是否存在 应该存在or不应该存在
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		id,
	})
	userCheck, err := stub.GetState(userKey)
	if err == nil && len(userCheck) != 0 {
		return shim.Error("User already exist!")
	}
	//写入状态
	user := &User{
		Name:           name,
		Id:             id,
		Account:        0,
		Reputation:     60,
		Info:           info,
		Skills:         make([]string, 0),
		Profession:     make([]string, 0),
		PostTasks:      make([]string, 0),
		OngoingTasks:   make([]string, 0),
		CompleteTasks:  make([]string, 0),
		Otherplatforms: make(map[string]int),
	}
	//序列化对象
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put user error %s", err))
	}
	return shim.Success(nil)
}

func userAddSkill(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	userid := string(si.GetIdBytes())
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		userid,
	})
	userCheck, err := stub.GetState(userKey)
	if err != nil || len(userCheck) == 0 {
		return shim.Error("User not found!")
	}
	user := new(User)
	err = json.Unmarshal(userCheck, user)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error %s", err))
	}
	for _, skill := range args {
		user.Skills = append(user.Skills, skill)
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update user error %s", err))
	}
	return shim.Success(nil)
}

func userAddPro(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	userid := string(si.GetIdBytes())
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		userid,
	})
	userCheck, err := stub.GetState(userKey)
	if err != nil || len(userCheck) == 0 {
		return shim.Error("User not found!")
	}
	user := new(User)
	err = json.Unmarshal(userCheck, user)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error %s", err))
	}
	for _, profession := range args {
		user.Profession = append(user.Profession, profession)
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update user error %s", err))
	}
	return shim.Success(nil)
}

func userImport(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	userid := string(si.GetIdBytes())
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		userid,
	})
	userCheck, err := stub.GetState(userKey)
	if err != nil || len(userCheck) == 0 {
		return shim.Error("User not found!")
	}
	user := new(User)
	err = json.Unmarshal(userCheck, user)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error %s", err))
	}
	length := len(args)
	if length%2 != 0 {
		return shim.Error("Incorrect number of arguments.")
	}
	for i := 0; i < length; i += 2 {
		user.Otherplatforms[args[i]], err = strconv.Atoi(args[i+1])
		if err != nil {
			return shim.Error("Invalid year.")
		}
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update user error %s", err))
	}
	return shim.Success(nil)
}

//用户发布一个任务输入name detail，id根据用户以及任务序号自动生成，初始状态为未接收 参数：Name id type detail reward posttime receivetime deadline requirement[]
func taskPost(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//参数 Name detail
	if len(args) < 8 {
		return shim.Error("Incorrect number of arguments. At least 8")
	}
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	userid := string(si.GetIdBytes())
	name := args[0]
	taskid := args[1]
	tasktype := args[2]
	detail := args[3]
	reward, _ := strconv.Atoi(args[4])
	local, _ := time.LoadLocation("Local")
	posttime, _ := time.ParseInLocation("2006-01-02 15:04:05", args[5], local)
	receivetime, _ := time.ParseInLocation("2006-01-02 15:04:05", args[6], local)
	deadline, _ := time.ParseInLocation("2006-01-02 15:04:05", args[7], local)
	if userid == "" || name == "" || detail == "" || taskid == "" || reward < 0 {
		return shim.Error("Invalid args")
	}
	if tasktype != "competition" && tasktype != "one2one" && tasktype != "private" {
		return shim.Error("Task type set wrong!")
	}
	beforeOrAfter := receivetime.After(posttime)
	if !beforeOrAfter {
		return shim.Error("ReceiveTime set wrong! ")
	}
	beforeOrAfter = deadline.After(posttime)
	if !beforeOrAfter {
		return shim.Error("Deadline set wrong!")
	}
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		userid,
	})
	userCheck, err := stub.GetState(userKey)
	if err != nil || len(userCheck) == 0 {
		return shim.Error("User not found!")
	}
	taskKey, err := stub.CreateCompositeKey("Task", []string{
		"task",
		taskid,
	})
	taskCheck, err := stub.GetState(taskKey)
	if err == nil && len(taskCheck) != 0 {
		return shim.Error("Taskid already exist!")
	}
	//1.写入任务 2.更新用户 3.写入任务提交记录
	user := new(User)
	err = json.Unmarshal(userCheck, user)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error %s", err))
	}
	if user.Account < reward {
		return shim.Error("Balance is not enough!")
	}
	state := "posted"
	var requirement []string
	if len(args) > 8 {
		requirement = args[8:]
	} else {
		requirement = make([]string, 0)
	}
	task := &Task{
		Name:          name,
		Id:            taskid,
		Type:          tasktype,
		Detail:        detail,
		Reward:        reward,
		State:         state,
		ReceiveTime:   receivetime,
		Deadline:      deadline,
		RequesterId:   userid,
		Candidate:     make(map[string]string),
		FinalWorker:   "",
		FinalSolution: "",
		Requirement:   requirement,
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal task error %s", err))
	}
	err = stub.PutState(taskKey, taskBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put task error %s", err))
	}
	user.PostTasks = append(user.PostTasks, taskid)
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update user error %s", err))
	}
	//t := time.Now()
	record := &WorkRecord{
		TaskId:    taskid,
		Requester: user.Id,
		Worker:    "no worker",
		Type:      "post",
		Time:      posttime,
	}
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal record error %s", err))
	}
	recordKey, err := stub.CreateCompositeKey("WorkRecord", []string{
		"record",
		taskid,
		user.Id,
		"no worker",
		"post",
		posttime.Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	}
	//err = stub.PutState(constructRecordKey(taskid, user.Id, ""), recordBytes)
	err = stub.PutState(recordKey, recordBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put record error %s", err))
	}
	return shim.Success(nil)
}

//更新任务的数据
func taskUpdate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	taskid := args[0]
	hash := args[1]
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	userid := string(si.GetIdBytes())
	if userid == "" || taskid == "" || hash == "" {
		return shim.Error("Invalid args")
	}
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		userid,
	})
	userCheck, err := stub.GetState(userKey)
	if err != nil || len(userCheck) == 0 {
		// return shim.Error("request User not found!")
		return shim.Error(fmt.Sprintf("User not found! %s", userid))
	}
	taskKey, err := stub.CreateCompositeKey("Task", []string{
		"task",
		taskid,
	})
	taskCheck, err := stub.GetState(taskKey)
	if err != nil || len(taskCheck) == 0 {
		return shim.Error("Task not found!")
	}
	task := new(Task)
	err = json.Unmarshal(taskCheck, task)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal task error %s", err))
	}
	task.Detail = task.Detail + "\n" + hash
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal task error %s", err))
	}
	err = stub.PutState(taskKey, taskBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update task error %s", err))
	}
	return shim.Success(nil)
}

//用户接收任务，需要根据任务id检查该任务是否存在，状态是否为未接收 参数：TaskID time
func taskReceive(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	taskid := args[0]
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", args[1], local)
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	workerid := string(si.GetIdBytes())
	if workerid == "" || taskid == "" {
		return shim.Error("Invalid args")
	}
	workerKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		workerid,
	})
	workerCheck, err := stub.GetState(workerKey)
	if err != nil || len(workerCheck) == 0 {
		return shim.Error("User not found!")
	}
	taskKey, err := stub.CreateCompositeKey("Task", []string{
		"task",
		taskid,
	})
	taskCheck, err := stub.GetState(taskKey)
	if err != nil || len(taskCheck) == 0 {
		return shim.Error("Task not found!")
	}
	task := new(Task)
	err = json.Unmarshal(taskCheck, task)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal task error %s", err))
	}
	if task.Type != "competition" {
		if task.State != "posted" {
			return shim.Error("Task state error!")
		}
	}
	beforeOrAfter := task.ReceiveTime.After(t)
	if !beforeOrAfter {
		return shim.Error("Task Receive is overtime!")
	}
	worker := new(User)
	err = json.Unmarshal(workerCheck, worker)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal worker error %s", err))
	}
	//加一步worker资质检查
	reputationRe, err := strconv.Atoi(task.Requirement[0])
	if worker.Reputation < reputationRe {
		return shim.Error("Worker do not meet the requirements!")
	}
	//1.修改任务状态 2.加入worker的worktasks 3.加入接收记录
	task.State = "received"
	if task.Type == "competition" {
		task.Candidate[workerid] = ""
	} else {
		task.FinalWorker = workerid
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal task error %s", err))
	}
	err = stub.PutState(taskKey, taskBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update task error %s", err))
	}
	worker.OngoingTasks = append(worker.OngoingTasks, taskid)
	workerBytes, err := json.Marshal(worker)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal worker error %s", err))
	}
	err = stub.PutState(workerKey, workerBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update worker error %s", err))
	}
	//t := time.Now()
	record := &WorkRecord{
		TaskId:    taskid,
		Requester: task.RequesterId,
		Worker:    workerid,
		Type:      "receive",
		Time:      t,
	}
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal record error %s", err))
	}
	recordKey, err := stub.CreateCompositeKey("WorkRecord", []string{
		"record",
		taskid,
		task.RequesterId,
		workerid,
		"receive",
		t.Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	}
	//err = stub.PutState(constructRecordKey(taskid, requesterid, workerid), recordBytes)
	err = stub.PutState(recordKey, recordBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put record error %s", err))
	}
	return shim.Success(nil)
}

//参数：TaskID solution time    worker上传解决方案 需要修改：任务状态 任务解决方案
func taskCommit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	taskid := args[0]
	solution := args[1]
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", args[2], local)
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	workerid := string(si.GetIdBytes())
	if workerid == "" || taskid == "" {
		return shim.Error("Invalid args")
	}
	workerKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		workerid,
	})
	workerCheck, err := stub.GetState(workerKey)
	if err != nil || len(workerCheck) == 0 {
		return shim.Error("User not found!")
	}
	taskKey, err := stub.CreateCompositeKey("Task", []string{
		"task",
		taskid,
	})
	taskCheck, err := stub.GetState(taskKey)
	if err != nil || len(taskCheck) == 0 {
		return shim.Error("Task not found!")
	}
	worker := new(User)
	err = json.Unmarshal(workerCheck, worker)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal worker error %s", err))
	}
	task := new(Task)
	err = json.Unmarshal(taskCheck, task)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal task error %s", err))
	}
	if task.Type == "competition" {
		_, exist := task.Candidate[workerid]
		if !exist {
			return shim.Error("Task candidate error!")
		}
	} else {
		if task.FinalWorker != workerid {
			return shim.Error("Task and worker not match!")
		}
	}
	exist := false
	//	var workerindex int
	for _, workertaskid := range worker.OngoingTasks {
		if workertaskid == taskid {
			exist = true
			//			workerindex = i
			break
		}
	}
	if !exist {
		return shim.Error("Ongoing tasks error!")
	}
	if task.State != "received" {
		return shim.Error("Task state error!")
	}
	//1.修改任务状态 2.修改requester的Posttasks指针 3.修改worker的worktasks指针 4.加入提交记录
	if solution == "" {
		return shim.Error("Invalid solution!")
	}
	if task.Type == "competition" {
		task.Candidate[workerid] = solution
	} else {
		task.FinalSolution = solution
	}
	beforeOrAfter := task.Deadline.After(t)
	if !beforeOrAfter {
		worker.Reputation = worker.Reputation - 10
		task.State = "overtime"
	} else {
		task.State = "committed"
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal task error %s", err))
	}
	err = stub.PutState(taskKey, taskBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update task error %s", err))
	}
	workerBytes, err := json.Marshal(worker)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal worker error %s", err))
	}
	err = stub.PutState(workerKey, workerBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update worker error %s", err))
	}
	//t := time.Now()
	record := &WorkRecord{
		TaskId:    taskid,
		Requester: task.RequesterId,
		Worker:    workerid,
		Type:      "commit",
		Time:      t,
	}
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal record error %s", err))
	}
	recordKey, err := stub.CreateCompositeKey("WorkRecord", []string{
		"record",
		taskid,
		task.RequesterId,
		workerid,
		"commit",
		t.Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	}
	//err = stub.PutState(constructRecordKey(taskid, requesterid, workerid), recordBytes)
	err = stub.PutState(recordKey, recordBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put record error %s", err))
	}
	if !beforeOrAfter {
		return shim.Success([]byte("commit overtime"))
	}
	return shim.Success(nil)
}

//参数：taskID workerid reward比例 time
func rewardAllocate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	taskid := args[0]
	workerid := args[1]
	rate, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid rate!")
	}
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", args[3], local)
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	requesterid := string(si.GetIdBytes())
	if requesterid == "" || workerid == "" || taskid == "" {
		return shim.Error("Invalid args")
	}
	requesterKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		requesterid,
	})
	requesterCheck, err := stub.GetState(requesterKey)
	if err != nil || len(requesterCheck) == 0 {
		// return shim.Error("request User not found!")
		return shim.Error(fmt.Sprintf("request User not found! %s", requesterid))

	}
	workerKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		workerid,
	})
	workerCheck, err := stub.GetState(workerKey)
	if err != nil || len(workerCheck) == 0 {
		// return shim.Error("worker User not found!")
		return shim.Error(fmt.Sprintf("worker User not found! %s", workerid))
	}
	taskKey, err := stub.CreateCompositeKey("Task", []string{
		"task",
		taskid,
	})
	taskCheck, err := stub.GetState(taskKey)
	if err != nil || len(taskCheck) == 0 {
		return shim.Error("Task not found!")
	}
	requester := new(User)
	err = json.Unmarshal(requesterCheck, requester)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal requester error %s", err))
	}
	worker := new(User)
	err = json.Unmarshal(workerCheck, worker)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal worker error %s", err))
	}
	task := new(Task)
	err = json.Unmarshal(taskCheck, task)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal task error %s", err))
	}
	if task.State != "committed" {
		return shim.Error("Task state error!")
	}
	exist := false
	for _, postT := range requester.PostTasks {
		if postT == taskid {
			exist = true
		}
	}
	if !exist {
		return shim.Error("Task and requester not match!")
	}
	exist = false
	var workerindex int
	for i, T := range worker.OngoingTasks {
		if T == taskid {
			exist = true
			workerindex = i
		}
	}
	if !exist {
		return shim.Error("Task and worker not match!")
	}
	if task.Type == "competition" {
		s, exist := task.Candidate[workerid]
		if !exist {
			return shim.Error("Task candidate error!")
		}
		if s == "" {
			return shim.Error("Invalid solution!")
		}
	} else {
		if task.FinalWorker != workerid {
			return shim.Error("Task and worker not match!")
		}
		if task.FinalSolution == "" {
			return shim.Error("Invalid solution!")
		}
	}
	reward := task.Reward * rate / 100
	if requester.Account < reward {
		return shim.Error("Requester account balance is not enough!")
	}
	//1.修改任务状态，填入奖励分配时间 2.修改requester的Posttasks指针，修改账户余额
	//3.修改worker的worktasks指针，修改账户余额 4.加入奖励分配记录
	task.State = "finished"
	if task.Type == "competition" {
		task.FinalWorker = workerid
		task.FinalSolution, _ = task.Candidate[workerid]
	}
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal task error %s", err))
	}
	err = stub.PutState(taskKey, taskBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update task error %s", err))
	}
	worker.Account = worker.Account + reward
	worker.OngoingTasks = append(worker.OngoingTasks[:workerindex], worker.OngoingTasks[workerindex+1:]...)
	worker.CompleteTasks = append(worker.CompleteTasks, taskid)
	if rate < 60 {
		worker.Reputation = worker.Reputation - (100-rate)/10
	} else {
		worker.Reputation = worker.Reputation + (rate-60)/10
	}
	workerBytes, err := json.Marshal(worker)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal worker error %s", err))
	}
	err = stub.PutState(workerKey, workerBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update worker error %s", err))
	}
	requester.Account = requester.Account - reward
	requesterBytes, err := json.Marshal(requester)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal requester error %s", err))
	}
	err = stub.PutState(requesterKey, requesterBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update requester error %s", err))
	}
	//t := time.Now()
	record := &WorkRecord{
		TaskId:    taskid,
		Requester: requesterid,
		Worker:    workerid,
		Type:      "reward " + strconv.Itoa(reward),
		Time:      t,
	}
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal record error %s", err))
	}
	recordKey, err := stub.CreateCompositeKey("WorkRecord", []string{
		"record",
		taskid,
		requesterid,
		workerid,
		"reward " + strconv.Itoa(reward),
		t.Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return shim.Error(fmt.Sprintf("create key error %s", err))
	}
	//err = stub.PutState(constructRecordKey(taskid, requesterid, workerid), recordBytes)
	err = stub.PutState(recordKey, recordBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put record error %s", err))
	}
	return shim.Success(nil)
}

//充值 参数：金额
func recharge(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	amountstr := args[0]
	creatorByte, err := stub.GetCreator()
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	userid := string(si.GetIdBytes())
	if userid == "" || amountstr == "" {
		return shim.Error("Invalid args")
	}
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		userid,
	})
	userCheck, err := stub.GetState(userKey)
	if err != nil || len(userCheck) == 0 {
		return shim.Error("User not found!")
	}
	amount, err := strconv.Atoi(amountstr)
	user := new(User)
	err = json.Unmarshal(userCheck, user)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error %s", err))
	}
	user.Account = user.Account + amount
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}
	err = stub.PutState(userKey, userBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("update user error %s", err))
	}
	return shim.Success(nil)
}

//参数：用户ID
func userQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//获取当前用户证书
	creatorByte, err := stub.GetCreator() //获取的是msp里面的signcerts
	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creatorByte, si)
	if err != nil {
		return shim.Error(fmt.Sprintf("get id error %s", err))
	}
	id := string(si.GetIdBytes())
	// if len(args) != 1 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 1")
	// }
	// userid := args[0]
	userid := id
	if userid == "" {
		return shim.Error("Invalid args")
	}
	userKey, err := stub.CreateCompositeKey("User", []string{
		"user",
		userid,
	})
	userCheck, err := stub.GetState(userKey)
	if err != nil || len(userCheck) == 0 {
		return shim.Error("User not found!")
	}
	return shim.Success(userCheck)
}

//参数：任务ID
func taskQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	taskid := args[0]
	if taskid == "" {
		return shim.Error("Invalid args")
	}
	taskKey, err := stub.CreateCompositeKey("Task", []string{
		"task",
		taskid,
	})
	taskCheck, err := stub.GetState(taskKey)
	if err != nil || len(taskCheck) == 0 {
		return shim.Error("Task not found!")
	}
	return shim.Success(taskCheck)
}

//参数：taskid
func recordQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	taskid := args[0]
	if taskid == "" {
		return shim.Error("Invalid args")
	}
	taskKey, err := stub.CreateCompositeKey("Task", []string{
		"task",
		taskid,
	})
	taskCheck, err := stub.GetState(taskKey)
	if err != nil || len(taskCheck) == 0 {
		return shim.Error("Task not found!")
	}
	keys := make([]string, 0)
	keys = append(keys, "record")
	keys = append(keys, taskid)
	result, err := stub.GetStateByPartialCompositeKey("WorkRecord", keys)
	if err != nil {
		return shim.Error(fmt.Sprintf("query record error %s", err))
	}
	defer result.Close()
	records := make([]*WorkRecord, 0)
	for result.HasNext() {
		recordVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query error %s", err))
		}
		record := new(WorkRecord)
		err = json.Unmarshal(recordVal.GetValue(), record)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error %s", err))
		}
		records = append(records, record)
	}
	recordsBytes, err := json.Marshal(records)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error %s", err))
	}

	return shim.Success(recordsBytes)
}

func alluserQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	result, err := stub.GetStateByPartialCompositeKey("User", []string{"user"})
	if err != nil {
		return shim.Error(fmt.Sprintf("query user error %s", err))
	}
	defer result.Close()
	users := make([]*User, 0)
	for result.HasNext() {
		userVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query error %s", err))
		}
		user := new(User)
		err = json.Unmarshal(userVal.GetValue(), user)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error %s", err))
		}
		users = append(users, user)
	}
	usersBytes, err := json.Marshal(users)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error %s", err))
	}
	return shim.Success(usersBytes)
}

func alltaskQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	result, err := stub.GetStateByPartialCompositeKey("Task", []string{"task"})
	if err != nil {
		return shim.Error(fmt.Sprintf("query task error %s", err))
	}
	defer result.Close()
	tasks := make([]*Task, 0)
	for result.HasNext() {
		taskVal, err := result.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("query error %s", err))
		}
		task := new(Task)
		err = json.Unmarshal(taskVal.GetValue(), task)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal error %s", err))
		}
		tasks = append(tasks, task)
	}
	tasksBytes, err := json.Marshal(tasks)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal error %s", err))
	}
	return shim.Success(tasksBytes)
}

func main() {
	err := shim.Start(new(Simple))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
