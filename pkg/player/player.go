package player

type Player interface {
	Play(fileName string) error
	Shutdown()
}
