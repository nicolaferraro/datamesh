package eventlog

import (
	"github.com/nicolaferraro/datamesh/protobuf"
	"bytes"
	"errors"
	"strconv"
	"bufio"
	"os"
	"github.com/golang/protobuf/proto"
)

type EventLogReader struct {
	eventLog	*EventLog
	Top			uint64
	offset		int64
	current		uint64
	reader		*bufio.Reader
}

func newEventLogReader(eventLog *EventLog, top uint64) (*EventLogReader, error) {
	if reader, err := os.OpenFile(eventLog.file.Name(), os.O_RDONLY, 0644); err != nil {
		return nil, err
	} else {
		logReader := EventLogReader{
			eventLog: eventLog,
			Top: top,
			current: 1,
			reader: bufio.NewReader(reader),
		}
		return &logReader, nil
	}
}

func (r *EventLogReader) Next() (*protobuf.Event, error) {
	if r.current > r.Top {
		return nil, errors.New("Stream finished")
	}

	buffer, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	length, err := findNextRecordLength(buffer)
	if err != nil {
		return nil, err
	}

	record := make([]byte, length + 1)
	numRead, err := fillBuffer(r.reader, record)
	if err != nil {
		return nil, err
	}
	record = record[0:len(record)-1]
	if numRead - 1 != int(length) {
		return nil, errors.New("Wrong number of bytes. Expected " + strconv.Itoa(int(length)) + " got " + strconv.Itoa(numRead))
	}

	unescaped := unescape(record)

	var event protobuf.Event
	if err := proto.Unmarshal(unescaped, &event); err != nil {
		return nil, err
	}

	r.current++
	return &event, nil
}

func fillBuffer(reader *bufio.Reader, buffer []byte) (int, error) {
	return fillBufferAcc(reader, buffer, 0)
}

func fillBufferAcc(reader *bufio.Reader, buffer []byte, acc int) (int, error) {
	numRead, err := reader.Read(buffer)
	if err != nil {
		return acc + numRead, err
	}

	if numRead == len(buffer) {
		return acc + numRead, err
	} else {
		return fillBufferAcc(reader, buffer[numRead:], acc + numRead)
	}
}

func findNextRecordLength(buffer []byte) (uint64, error) {
	if buffer[0] != RecordSeparator {
		return 0, errors.New("Unexpected character found instead of " + string(RecordSeparator))
	}

	part := buffer[1:]
	comm := bytes.IndexByte(part, ',')
	if comm < 0 {
		return 0, errors.New("Wrong header format")
	}
	end := bytes.IndexByte(part, '\n')
	if end < 0 {
		return 0, errors.New("Wrong header format")
	}

	recordLenBuf := part[comm + 1:end]
	recordLen, err := strconv.ParseUint(string(recordLenBuf), 10, 64)
	if err != nil {
		return 0, errors.New("Wrong header format")
	}

	return recordLen, nil
}

func unescape(data []byte) []byte {
	if !bytes.Contains(data, []byte{RecordSeparator}) {
		return data
	}

	escaped := data
	escaped = bytes.Replace(escaped, []byte{RecordSeparator, RecordSeparator}, []byte{RecordSeparator}, -1)

	return escaped
}