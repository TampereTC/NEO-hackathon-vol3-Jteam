package message

type Message struct {
  MessageCounter int `json:"message_counter,omitempty"`
  MessageText string `json:"message_text,omitempty"`
  LastMessage bool `json:"last_message,omitempty"`
}
