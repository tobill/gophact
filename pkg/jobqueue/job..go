package jobqueue

//Job model 
type Job interface {
	Execute() int
	SetStatus(s int)
	GetStatus() (int) 
	GetObjectID() (uint64)
	GetJobName() (string)
}
