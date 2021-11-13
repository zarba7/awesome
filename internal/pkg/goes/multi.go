package goes

type Multi struct {
	Core
	tasks []Task
}

func (the *Multi) OnBlock() {
	for _, t := range the.tasks {
		t.Do()
	}
	the.tasks = the.tasks[:]
}

func (the *Multi) DoTask(t Task) {
	the.Core.DoTask(&tAppend{the,t })
}

type tAppend struct {
	*Multi
	t Task
}

func (the *tAppend) Do() {
	the.tasks = append(the.tasks, the.t)
}
