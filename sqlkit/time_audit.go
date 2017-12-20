package sqlkit

import "time"

func TimeAuditCols() (cols []string) {
	cols = []string{"created_at", "updated_at"}
	return
}

type TimeAudit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *TimeAudit) InitTimeAudit(isCreate bool) {
	if isCreate {
		a.CreatedAt = time.Now()
	}
}

func UpdateTimeAuditCols() (clause string) {
	clause = ", updated_at = NOW()"
	return
}

func InsertTimeAuditValues() (clause string) {
	clause = "NOW(), NOW()"
	return
}
