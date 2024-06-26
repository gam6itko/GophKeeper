package serialize

// Text просто заглушка для преобразования текста.
type Text struct{}

func (ths Text) Serialize(text string) ([]byte, error) {
	return []byte(text), nil
}

func (ths Text) Deserialize(data []byte) (string, error) {
	return string(data), nil
}
