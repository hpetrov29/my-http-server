package http

type Request struct{
	Method string
	Path string
	Params []string
	Headers map[string]string
}