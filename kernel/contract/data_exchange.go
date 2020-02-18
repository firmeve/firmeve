package contract

type (
	DataExchanger interface {
		Encode(protocol Protocol,v interface{}) ([]byte,error)
		Decode(protocol Protocol,v interface{}) (interface{},error)
	}
)