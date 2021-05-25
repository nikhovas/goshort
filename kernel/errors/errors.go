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

type GenericError struct {
	Name string
}

func (e *GenericError) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["name"] = e.Name
	data["type"] = "error"
	return data
}

func (e *GenericError) Error() string {
	return "Error " + e.Name
}

var UrlNotFoundError = &GenericError{Name: "Common.UrlNotFound"}
var KeyNotFoundError = &GenericError{Name: "Common.KeyNotFoundError"}
var ValueAlreadyExistsError = &GenericError{Name: "Common.ValueAlreadyExistsError"}
var GenericKeysAreNotSupported = &GenericError{Name: "Common.GenericKeysAreNotSupported"}
