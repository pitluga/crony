package crony

import (
	"io"
	"io/ioutil"
)

type ReaderConsumer interface {
	ConsumeReader(reader io.Reader) error
}

type Sink struct {
	Content string
}

func NewSink() *Sink {
	return &Sink{}
}

func (sink *Sink) ConsumeReader(reader io.Reader) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	sink.Content = string(bytes)
	return nil
}
