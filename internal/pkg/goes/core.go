package goes

type Task interface {
	Do()
}

type Go interface {
	Go(g Go, backlog int)
	DoTask(t Task)
	OnBlock()
	OnClose()
	Quit()
}

type Core struct {
	chTask chan Task
	quit   chan struct{}
	g      Go
}

func (the *Core) Go(g Go, backlog int) {
	the.chTask = make(chan Task, backlog)
	the.quit = make(chan struct{})
	the.g = g
	go the.run(g)
}

func (the *Core) run(g Go) {
	defer close(the.quit)
	for t := range the.chTask {
	LOOP:
		t.Do()
		if the.g == nil {
			break
		}
		select {
		case t = <-the.chTask:
			goto LOOP
		default:
			g.OnBlock()
		}
	}
	g.OnBlock()
	g.OnClose()
}

func (the *Core) DoTask(t Task) {
	if the.g != nil {
		select {
		case the.chTask <- t:
		case <-the.quit:
		}
	}
}

func (the *Core) Quit() {
	the.DoTask(&closer{the})
	<-the.quit
}

type closer struct{ *Core }

func (the *closer) Do() { the.g = nil }
