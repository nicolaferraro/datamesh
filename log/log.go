package log

import (
	"os"
	"bytes"
	"strconv"
	"path"
	"io"
	"github.com/nicolaferraro/datamesh/protobuf"
	"github.com/golang/protobuf/proto"
	"github.com/nicolaferraro/datamesh/common"
)

const (
	EventLogFileName	= "event.log"
	RecordSeparator		= '!'

	lookupWindowSize	= 4096
	lookupWindowShift	= 3072
)

type Log struct {
	directory	string
	file		*os.File
	entries		uint64
	Cache		*LogCache
}

func NewLog(directory string) (*Log, error) {
	log := &Log{
		directory: directory,
		Cache: NewLogCache(),
	}

	if err := log.init(); err != nil {
		return nil, err
	}

	return log, nil
}

func (log *Log) Append(data []byte) (uint64, error) {
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

/*
 * Implements common.EventConsumer
 */
func (log *Log) Consume(evt *protobuf.Event) error {
	msg, err := proto.Marshal(evt)
	if err != nil {
		return err
	}

	if err = log.Cache.Accept(evt); err != nil {
		return err
	}

	_, err = log.Append(msg)
	return err
}

/*
 * Implements common.Observer
 */
func (log *Log) OnNotification(notification common.Notification) error {
	if notification.Type == common.NotificationTypeProjectionVersion {
		version := notification.Payload.(uint64)
		log.Cache.Prune(version)
	}
	return nil
}

func (log *Log) Sync() error {
	return log.file.Sync()
}

func (log *Log) init() error {
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
