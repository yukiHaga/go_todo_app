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
	ID        TaskID     `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Status    TaskStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type Tasks []*Task
