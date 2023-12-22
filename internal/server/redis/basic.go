package redis

type Session struct {
	Id     string
	Values map[any]any
}
