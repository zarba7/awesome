package saga

type Response struct {
	TrxID int64
	Step string
	NextArgs []byte
	Code int32
	Message string
}

type Request struct {
	TrxID int64
	Step string
	Args []byte
}
