package queue

type MockQueue struct{}

func NewMockQueue() *MockQueue {
	return &MockQueue{}
}

func (m *MockQueue) Send(message string) (string, error) {
	return "mock-message-id", nil
}

func (m *MockQueue) Receive(maxMessages int) ([]Message, error) {
	return []Message{}, nil
}

func (m *MockQueue) Delete(receiptHandle string) error {
	return nil
}
