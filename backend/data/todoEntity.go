package data

import "encoding/json"

func UnmarshalTodo(data []byte) (Todo, error) {
	var r Todo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Todo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Todo struct {
	UserID    int64  `json:"userId"`
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
