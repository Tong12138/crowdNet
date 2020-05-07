//sturctinfo.go
package controller

import(
	// "encoding/json"
	"time"
)

type User struct {
	Name           string         `json:"name"`
	Id             string         `json:"id"`
	Account        int            `json:"account balance"`
	Reputation     int            `json:"reputation"`
	Info           string         `json:"detail information"`
	Skills         []string       `json:"skills"`
	Profession     []string       `json:"profession"`
	PostTasks      []string       `json:"post tasks"`
	OngoingTasks   []string       `json:"ongoing tasks"`
	CompleteTasks  []string       `json:"complete tasks"`
	Otherplatforms map[string]int `json:"other platforms workyears"`
}

type Task struct {
	Name          string            `json:"name"`
	Id            string            `json:"id"`
	Type          string            `json:"task type"` //1.competition 2.one2one 3.private
	Detail        string            `json:"detail"`
	Reward        int               `json:"reward"`
	State         string            `json:"state"`
	ReceiveTime   time.Time         `json:"receive time"`
	Deadline      time.Time         `json:"deadline"`
	RequesterId   string            `json:"requester"`
	Candidate     map[string]string `json:"candidate worker and solution"`
	FinalWorker   string            `json:"final worker"`
	FinalSolution string            `json:"final solution"`
	Requirement   []string          `json:"worker requirement"` //reputation, complete num, skill, profession
}

type WorkRecord struct {
	TaskId    string    `json:"task_id"`
	Requester string    `json:"requester"`
	Worker    string    `json:"worker"`
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
}
