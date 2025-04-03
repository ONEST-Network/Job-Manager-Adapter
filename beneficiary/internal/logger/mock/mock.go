package mock

type MockLogger struct {
	InfoCalls  []struct{ Msg string; Fields map[string]interface{} }
	ErrorCalls []struct{ Msg string; Err error; Fields map[string]interface{} }
	WarnCalls  []struct{ Msg string; Fields map[string]interface{} }
	DebugCalls []struct{ Msg string; Fields map[string]interface{} }
}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (m *MockLogger) Info(msg string, fields map[string]interface{}) {
	m.InfoCalls = append(m.InfoCalls, struct {
		Msg    string
		Fields map[string]interface{}
	}{msg, fields})
}

func (m *MockLogger) Error(msg string, err error, fields map[string]interface{}) {
	m.ErrorCalls = append(m.ErrorCalls, struct {
		Msg    string
		Err    error
		Fields map[string]interface{}
	}{msg, err, fields})
}

func (m *MockLogger) Debug(msg string, fields map[string]interface{}) {
	m.DebugCalls = append(m.DebugCalls, struct {
		Msg    string
		Fields map[string]interface{}
	}{msg, fields})
}

func (m *MockLogger) Warn(msg string, fields map[string]interface{}) {
	m.WarnCalls = append(m.WarnCalls, struct {
		Msg    string
		Fields map[string]interface{}
	}{msg, fields})
} 