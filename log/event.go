package log


type AppendEvent struct {
	Id	uint64
}

func NewAppendEvent(id uint64) *AppendEvent {
	return &AppendEvent{
		Id: id,
	}
}
