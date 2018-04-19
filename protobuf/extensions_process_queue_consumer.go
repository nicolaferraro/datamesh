package protobuf


// Make it compatible with common.EventConsumer
type ProcessQueueConsumer struct {
	queue		DataMesh_ProcessQueueServer
	Closed		chan bool
}

func (consumer *ProcessQueueConsumer) Consume(event *Event) error {
	if err := consumer.queue.Send(event); err != nil {
		consumer.Closed <- true
		return err
	}
	return nil
}

func NewProcessQueueConsumer(queue DataMesh_ProcessQueueServer) *ProcessQueueConsumer {
	return &ProcessQueueConsumer{
		queue: queue,
	}
}
