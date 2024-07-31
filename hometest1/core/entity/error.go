package entity

const (
	ProductNotFound string = "product not found"
	EmptyQuantity   string = "empty quantity"
)

type Err struct {
	message string
	code    int
}

func (e Err) Error() string {
	return e.message
}

func (e Err) GetCode() int {
	return e.code
}

func (e Err) GetMessage() string {
	return e.message
}

func NewError(msg string, code int) Err {
	return Err{message: msg, code: code}
}
