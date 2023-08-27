package packer

type Request struct {
	Header int
	Data   map[string]interface{}
}

type Response struct {
	Header int
	Data   map[string]interface{}
}
