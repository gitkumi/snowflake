package queue

type Queue interface {
	Send(message string) (string, error)

	Receive(maxMessages int) ([]Message, error)

	Delete(receiptHandle string) error
}

type Message struct {
	ID            string
	Body          string
	ReceiptHandle string
}
