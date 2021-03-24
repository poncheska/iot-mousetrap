package utils

type StringLogger struct {
	Logs string
}

func NewStringLogger() *StringLogger {
	return &StringLogger{}
}

func (sl *StringLogger) Write(bs []byte) (int, error) {
	sl.Logs += string(bs)
	return len(bs), nil
}

func (sl *StringLogger) Clear() {
	sl.Logs = ""
}
