package protobuf


// Make it compatible with common.EventConsumer
type ProcessQueueConsumer struct {
	queue	 	DataMesh_ConnectServer
	Closed		chan bool
}

func (consumer *ProcessQueueConsumer) Consume(event *Event) error {
	if err := consumer.queue.Send(event); err != nil {
		consumer.Closed <- true
		return err
	}
	return nil
}

func NewProcessQueueConsumer(queue DataMesh_ConnectServer) *ProcessQueueConsumer {
	return &ProcessQueueConsumer{
		queue: queue,
	}
}
