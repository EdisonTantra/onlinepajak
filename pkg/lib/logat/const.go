package logat

import "time"

const (
	fieldCorrelationID = "_cid"
	fieldEvent         = "event"
	fieldData          = "data"
	fieldContext       = "context"
	fieldCaller        = "caller"
)

const (
	contextKeyUserID = "USER_ID"
)

const (
	callerSkipOffset = 2
)

const (
	EncoderJSON    = Encoder("json")
	EncoderConsole = Encoder("console")
)

const (
	LevelDebug = Level("debug")
	LevelInfo  = Level("info")
	LevelWarn  = Level("warn")
	LevelError = Level("error")
	LevelFatal = Level("fatal")
)

const (
	TimeEpoch       = Time("epoch")
	TimeEpochMilli  = Time("epoch_milli")
	TimeEpochNano   = Time("epoch_nano")
	TimeRFC3339     = Time(time.RFC3339)
	TimeRFC3339Nano = Time(time.RFC3339Nano)
	TimeISO8601     = Time("2006-01-02T15:04:05.000Z0700")
)
