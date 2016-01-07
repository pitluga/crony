package crony

import (
	"bufio"
	"fmt"
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

type StdoutPipe struct {
	prefix string
}

func NewStdoutPipe(prefix string) *StdoutPipe {
	return &StdoutPipe{prefix}
}

func (consumer *StdoutPipe) ConsumeReader(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("%v - %v\n", consumer.prefix, scanner.Text())
	}
	return scanner.Err()
}
