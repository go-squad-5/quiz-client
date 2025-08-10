package app

type STATUS string

const (
  STATUS_CREATED  STATUS = "created"
  STATUS_STARTED  STATUS = "started"
  STATUS_COMPLETED STATUS = "completed"
  STATUS_FAILED   STATUS = "failed"
)

type Session struct {
  ID        string
  UserID    string
  StartTime int64
  EndTime   int64
  Answers   map[string]string
  Score     int
  Status    STATUS
  CreatedAt int64
}
