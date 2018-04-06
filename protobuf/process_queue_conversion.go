package protobuf


// Make it compatible with common.EventConsumer
type ProcessQueueConsumer struct {
	queue		DataMesh_ProcessQueueServer
	isClosed	bool
	Closed		chan bool
}

func (consumer *ProcessQueueConsumer) Consume(event *Event) error {
	return consumer.queue.Send(event)
}

func (consumer *ProcessQueueConsumer) Close() {
	consumer.isClosed = true
	consumer.Closed <- true
}

func (consumer *ProcessQueueConsumer) IsClosed() bool {
	return consumer.isClosed
}

func NewProcessQueueConsumer(queue DataMesh_ProcessQueueServer) *ProcessQueueConsumer {
	return &ProcessQueueConsumer{
		queue: queue,
	}
}
