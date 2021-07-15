package logModules

//import (
//	"encoding/json"
//	"fmt"
//	"goshort/kernel"
//	"goshort/kernel/utils"
//	"goshort/types"
//)
//
//type FileDescriptorLogger struct {
//	types.LoggerBase
//	fd     int
//	Name   string
//	Kernel *kernel.Kernel
//}
//
//func CreateFileDescriptorLogger(kernel *kernel.Kernel) types.LoggerInterface {
//	return &FileDescriptorLogger{Kernel: kernel}
//}
//
//func (controller *FileDescriptorLogger) Init(config map[string]interface{}) error {
//	_ = controller.LoggerBase.Init(config)
//	controller.Name = config["name"].(string)
//	controller.fd = utils.UnwrapFieldOrDefault( config, "file_descriptor", 1).(int)
//	return nil
//}
//
//func (controller *FileDescriptorLogger) Run() error {
//	return nil
//}
//
//func (controller *FileDescriptorLogger) Send(le types.Log) error {
//	b, _ := json.Marshal(le.ToMap())
//	fmt.Fprintln()
//	println(string(b))
//	return nil
//}
//
//func (controller *FileDescriptorLogger) GetName() string {
//	return controller.name
//}
//
//func (controller *FileDescriptorLogger) GetType() string {
//	return "FileDescriptorLogger"
//}
