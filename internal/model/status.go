package model

// Status denotes the status of data
type Status int

// The following constant are the know status
const (
	Deleted  Status = -1
	Unknown  Status = 0
	Active   Status = 1
	Inactive Status = 2
)

// Value convert  into int
func (status Status) Value() int {
	return int(status)
}
