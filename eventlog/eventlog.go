package eventlog

import (
	"os"
	"bytes"
	"strconv"
	"path"
	"io"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/golang/protobuf/proto"
	"github.com/nicolaferraro/datamesh/common"
	"github.com/nicolaferraro/datamesh/notification"
	"context"
)

const (
	EventLogFileName	= "event.log"
	RecordSeparator		= '!'

	lookupWindowSize	= 4096
	lookupWindowShift	= 3072
)

type EventLog struct {
	directory		string
	file			*os.File
	entries			uint64
	serializer		*common.Serializer
	notificationBus	*notification.NotificationBus
}

func NewEventLog(ctx context.Context, directory string, bus *notification.NotificationBus) (*EventLog, error) {
	log := EventLog{
		directory: directory,
		notificationBus: bus,
	}
	log.serializer = common.NewSerializer(ctx, &log)

	if err := log.init(); err != nil {
		return nil, err
	}

	return &log, nil
}

func (log *EventLog) Consume(evt *protobuf.Event) error {
	// TODO need feedback
	log.serializer.Push(evt)
	return nil
}

func (log *EventLog) AppendRaw(data []byte) (uint64, error) {
	newSize := log.entries + 1
	escapedData := escape(data)
	record := make([]byte, 0)
	record = append(record, RecordSeparator)
	record = strconv.AppendUint(record, newSize, 10)
	record = append(record, ',')
	record = strconv.AppendUint(record, uint64(len(escapedData)), 10)
	record = append(record, '\n')
	record = append(record, escapedData...)
	record = append(record, '\n')
	if _, err := log.file.Write(record); err != nil {
		return log.entries, err
	}

	log.entries = newSize
	return newSize, nil
}

// implements common.Serializer callback
func (log *EventLog) ExecuteSerially(value interface{}) (bool, error) {
	if evt, ok := value.(*protobuf.Event); ok {
		msg, err := proto.Marshal(evt)
		if err != nil {
			return false, err
		}

		newSize, err := log.AppendRaw(msg)

		evtCopy := *evt
		evtCopy.Version = newSize

		// FIXME do it at every sync for all events
		log.notificationBus.Notify(notification.NewEventAppendedNotification(&evtCopy, false))

		return true, err
	}
	return false, nil
}

func (log *EventLog) Sync() error {
	return log.file.Sync()
}

func (log *EventLog) init() error {
	if err := os.MkdirAll(log.directory, 0755); err != nil {
		return err
	}

	name := path.Join(log.directory, EventLogFileName)
	f, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	fileSize := stat.Size()

	count, fileEnd, err := findLogStatus(f, fileSize)
	if err != nil {
		return err
	}

	if int64(fileEnd) < fileSize {
		f.Truncate(int64(fileEnd))
	}

	log.entries = count
	log.file = f
	return nil
}

func findLogStatus(file *os.File, size int64) (uint64, int, error) {
	offset := size - lookupWindowSize
	for offset > -lookupWindowSize {
		prg, fileEnd, found, err := findLogStatusFromOffset(file, offset, size)
		if err != nil {
			return 0, 0, err
		}
		if found {
			return prg, fileEnd, nil
		}

		offset = offset - lookupWindowShift
	}

	return 0, 0, nil
}

func findLogStatusFromOffset(file *os.File, offset int64, fileSize int64) (uint64, int, bool, error) {
	if offset < 0 {
		offset = 0
	}

	buffer := make([]byte, lookupWindowSize)
	length, err := file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		return 0, 0, false, err
	}

	buffer = buffer[0:length]
	limit := length

	for limit > 0 {
		pos := bytes.LastIndexByte(buffer[0:limit], RecordSeparator)
		if pos < 0 {
			return 0, 0, false, nil
		}

		prg, fileEnd, found := isControlRow(buffer, pos, fileSize)
		if found {
			return prg, fileEnd, true, nil
		}

		limit = pos
	}
	return 0, 0, false, nil
}

func isControlRow(buf []byte, offset int, fileSize int64) (uint64, int, bool) {
	if offset < 0 {
		return 0, 0, false
	}
	if offset > 0 && buf[offset-1] != '\n' {
		return 0, 0, false
	}

	if buf[offset] != RecordSeparator {
		return 0, 0, false
	}

	numStart := offset + 1
	part := buf[numStart:]
	comm := bytes.IndexByte(part, ',')
	if comm < 0 {
		return 0, 0, false
	}
	end := bytes.IndexByte(part, '\n')
	if end < 0 {
		return 0, 0, false
	}

	num := part[0:comm]
	prg, err := strconv.ParseUint(string(num), 10, 64)
	if err != nil {
		return 0, 0, false
	}

	recordLenBuf := part[comm + 1:end]
	recordLen, err := strconv.ParseUint(string(recordLenBuf), 10, 64)
	if err != nil {
		return 0, 0, false
	}

	recordEndOffset := int64(offset) + int64(end) + int64(recordLen) + 3
	if recordEndOffset > fileSize {
		return 0, 0, false
	}

	return prg, int(recordEndOffset), true
}

func escape(data []byte) []byte {
	if !bytes.Contains(data, []byte{RecordSeparator}) {
		return data
	}

	escaped := data
	escaped = bytes.Replace(data, []byte{RecordSeparator}, []byte{RecordSeparator, RecordSeparator}, -1)
	return escaped
}

func unescape(data []byte) []byte {
	if !bytes.Contains(data, []byte{RecordSeparator}) {
		return data
	}

	escaped := data
	escaped = bytes.Replace(escaped, []byte{RecordSeparator, RecordSeparator}, []byte{RecordSeparator}, -1)

	return escaped
}
