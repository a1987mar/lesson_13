package commands

const (
	PutCommandName    string = "put"
	GetCommandName    string = "get"
	DeleteCommandName string = "delete"
)

type NewCollection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Doc  Doc    `json:"Doc"`
}

type Doc struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
