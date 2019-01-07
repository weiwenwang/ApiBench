package common

type Request struct {
	Protocol int8
	Domain   string
	Port     int32
	Method   int8
	Path     string
	File     string
	Param    string
	C        int16
	Ps       int16
}
