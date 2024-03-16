package recorder

type Recorder interface {
	Record(filename string) error
	Shutdown()
}
