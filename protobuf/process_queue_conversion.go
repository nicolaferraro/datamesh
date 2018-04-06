package protobuf


// Make it compatible with common.EventConsumer
type ProcessQueueConsumer struct {
	queue DataMesh_ProcessQueueServer
}

func (consumer *ProcessQueueConsumer) Consume(event *Event) error {
	return consumer.queue.Send(event)
}

func NewProcessQueueConsumer(queue DataMesh_ProcessQueueServer) *ProcessQueueConsumer {
	return &ProcessQueueConsumer{
		queue: queue,
	}
}
