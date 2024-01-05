package bot

type Usecase interface {
	Consume() error
}
