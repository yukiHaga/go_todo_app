package entity

import "time"

type TaskID int64
type TaskStatus string

// iotaでバリューを数値として定義しても良いかも。その場合は、TaskStatusはStringインターフェースを満たした方が良い。
const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

type Task struct {
	ID      TaskID     `json:"id"`
	Title   string     `json:"title"`
	Status  TaskStatus `json:"status"`
	Created time.Time  `json:"created"`
}

type Tasks []*Task
