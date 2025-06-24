package human

import (
	"errors"
	"time"

	"github.com/corani/mcp-human-go/internal/config"
	"github.com/corani/mcp-human-go/internal/memory"
)

type Ask struct {
	conf   *config.Config
	memory *memory.MemoryDB
}

func NewAsk(conf *config.Config, mem *memory.MemoryDB) *Ask {
	return &Ask{
		conf:   conf,
		memory: mem,
	}
}

func (a *Ask) Ask(question, context string) (string, error) {
	if question == "" {
		return "", errors.New("question cannot be empty")
	}

	answerCh := make(chan string, 1)

	_, err := a.memory.Add(question, context, func(answer string) {
		answerCh <- answer
	})
	if err != nil {
		return "", err
	}

	select {
	case answer := <-answerCh:
		return answer, nil
	case <-time.After(time.Duration(a.conf.MaxWait) * time.Second):
		return "", errors.New("no human was available to answer the question within the timeout period")
	}
}
