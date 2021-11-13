package saga

type roleSrv struct {

}

func (r *roleSrv) OnEnter() {
	panic("implement me")
}

func (r *roleSrv) Do() {
	panic("implement me")
}

func (r *roleSrv) OnLeave() {
	panic("implement me")
}

func (r *roleSrv) String() {
	panic("implement me")
}

type friendSrv struct {

}

func (f *friendSrv) String() {
	panic("implement me")
}

func tt(){
	OFactory := &OrchestratorFactory{}

	r1 := &roleSrv{}
	f1 := &friendSrv{}
	r2 := &roleSrv{}

	table :=
		`{
			"role": { "ok": "friend", "fail":"end" },
			"friend": { "ok": "end", "fail":"role-rollback" },
			"role-rollback": { "ok": "end", "fail":"end" },
		}`
}

