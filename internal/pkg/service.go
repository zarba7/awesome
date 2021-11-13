package base


type Service interface {
	OnInit(opt Option) error
	OnQuit()
}
