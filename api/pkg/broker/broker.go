package broker

import "github.com/bovinxx/code-processor/api/model"

type Broker interface {
	PublishMessage(task model.Task) error
}
