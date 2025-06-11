package logger_hook

import (
	"bytes"

	"github.com/sirupsen/logrus"

	Library "mini-accounting/library"
)

type LogHook interface {
	Levels() []logrus.Level
	Fire(e *logrus.Entry) error
	GetBuffer() *bytes.Buffer
}

type LogHookImpl struct {
	// INIT BUFFER AS PUBLISHER
	buffer  *bytes.Buffer
	library Library.Library
}

func New(
	buffer *bytes.Buffer,
	library Library.Library,
) LogHook {
	return &LogHookImpl{
		buffer:  buffer,
		library: library,
	}
}

func (h *LogHookImpl) Levels() []logrus.Level {
	// MAKE THIS HOOK CAN BE CALLED BY ALL LEVEL LOG
	return logrus.AllLevels
}

func (h *LogHookImpl) Fire(e *logrus.Entry) error {
	// CONVERT "logrus.Entry" INTO JSON BYTE
	data, err := h.library.JsonMarshal(e.Data)
	// WHEN CONVERTION RETURNS ERROR
	if err != nil {
		return err
	}
	// UPDATE DATA THAT IS CONTAINED BY PUBLISHER
	h.buffer.Write(data)
	return nil
}

func (h *LogHookImpl) GetBuffer() *bytes.Buffer {
	// GET UPDATED DATA
	return h.buffer
}
