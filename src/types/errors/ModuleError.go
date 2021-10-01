package errors

//type ModuleErrorInterface interface {
//	kernel.AdvancedError
//}
//
//type ModuleError struct {
//	ModuleName string
//	ModuleType string
//	NextError  error
//}
//
//func (err ModuleError) Error() string {
//	return "Error occurred while executing middleware: " + err.NextError.Error()
//}
//
//func (err * ModuleError) ToJson() map[string]interface{} {
//	data := make(map[string]interface{})
//	data["moduleName"] = err.ModuleName
//	data["moduleType"] = err.ModuleType
//	data["nextError"] = kernel.GetErrorDescription(err.NextError)
//	return data
//}
//
//func GenerateModuleError(module kernel.ModuleInterface, err error) *ModuleError {
//	return &ModuleError{
//		ModuleName: module.GetName(),
//		ModuleType: module.GetType(),
//		NextError:  err,
//	}
//}
