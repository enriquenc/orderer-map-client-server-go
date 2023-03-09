package shared

const (
	AddItem    string = "add"
	RemoveItem string = "remove"
	GetItem    string = "get"
	GetAll     string = "getAll"
)

type Request struct {
	Action string
	Key    string
	Value  string
}

type TestDataAction struct {
	RequestData      Request
	ExpectedResponse interface{}
}
