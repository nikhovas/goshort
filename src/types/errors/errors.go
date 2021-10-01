package errors

//type MiddlewareError struct {
//	MiddlewareName string
//	NextError      error
//}
//
//func (err *MiddlewareError) Error() string {
//	return "Error occurred while executing middleware: " + err.NextError.Error()
//}
//
//func (err *MiddlewareError) AdvanceCheck() {}
//
//func (err *MiddlewareError) ToJson() map[string]interface{} {
//	data := make(map[string]interface{})
//	data["errorType"] = "middlewareError"
//	data["middlewareName"] = err.MiddlewareName
//	data["nextError"] = kernel.GetErrorDescription(err.NextError)
//	return data
//}

type SimpleErrorWrapper struct {
	Err error
}

func (e *SimpleErrorWrapper) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["name"] = "SimpleError"
	data["type"] = "error"
	data["details"] = e.Err.Error()
	return data
}

func (e *SimpleErrorWrapper) Error() string {
	return "Error SimpleError " + e.Err.Error()
}

type GenericLog struct {
	Name    string
	IsError bool
}

func (e *GenericLog) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["name"] = e.Name
	if e.IsError {
		data["type"] = "error"
	} else {
		data["type"] = "log"
	}

	return data
}

func (e *GenericLog) Error() string {
	return "Error " + e.Name
}

var ValueAlreadyExistsError = &GenericLog{Name: "Common.ValueAlreadyExistsError", IsError: true}
var GenericKeysAreNotSupported = &GenericLog{Name: "Common.GenericKeysAreNotSupported", IsError: true}
var NotFoundError = &GenericLog{Name: "Common.NotFoundError", IsError: true}
var NotImplementedError = &GenericLog{Name: "Common.NotImplementedError", IsError: true}

type BadConnectionError struct {
	Host      string
	Port      int
	Protocol  string
	Retryable bool
}

func (bce *BadConnectionError) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":     "Common.BadConnectionError",
		"type":     "error",
		"host":     bce.Host,
		"port":     bce.Port,
		"protocol": bce.Protocol,
	}
}

func (bce BadConnectionError) Error() string {
	return "Error Common.BadConnectionError"
}
