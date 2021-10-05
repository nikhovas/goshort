package logModules

import (
	"goshort/src/kernel/utils/other"
	"goshort/src/types"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func stringToConsoleLogString(data string) string {
	brackets := false
	if data == "" {
		brackets = true
	}

	for _, symbol := range []string{`\`, ` `, `"`, `[`, `]`, `(`, `)`} {
		if strings.Contains(data, symbol) {
			data = strings.ReplaceAll(data, symbol, `\`+symbol)
			brackets = true
		}
	}

	if brackets {
		data = `"` + data + `"`
	}

	return data
}

func toConsoleLogString(data interface{}) string {
	switch typedData := data.(type) {
	case map[string]interface{}:
		return "[" + mapToConsoleLogStringInternal(typedData) + "]"
	case string:
		return stringToConsoleLogString(typedData)
	case int:
		return strconv.Itoa(typedData)
	case []interface{}:
		res := "("
		for counter := range typedData {
			res += toConsoleLogString(typedData[counter])
			if counter != len(typedData)-1 {
				res += " "
			}
		}
		res += ")"
		return res
	default:
		return ""
	}
}

func mapToConsoleLogStringInternal(data map[string]interface{}) string {
	res := ""
	for k, v := range data {
		res += toConsoleLogString(k) + "=" + toConsoleLogString(v) + " "
	}
	res = res[:len(res)-1]
	return res
}

func logToConsoleString(le types.Log) string {
	return mapToConsoleLogStringInternal(le.ToMap())
}

type Console struct {
	types.LoggerBase
	name        string
	Kernel      types.KernelInterface
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func CreateConsole(kernel types.KernelInterface) types.LoggerInterface {
	return &Console{Kernel: kernel}
}

func (controller *Console) Init(config map[string]interface{}) error {
	_ = controller.LoggerBase.Init(config)
	controller.name = config["name"].(string)
	controller.infoLogger = log.New(os.Stdout, "INF: ", log.Ldate|log.Ltime)
	controller.errorLogger = log.New(os.Stderr, "ERR: ", log.Ldate|log.Ltime)
	return nil
}

func (controller *Console) Run(wg *sync.WaitGroup) error {
	wg.Done()
	controller.IsAvailableVal = 1
	return nil
}

func (controller *Console) Stop() error {
	return nil
}

func (controller *Console) Send(le types.Log) error {
	data := logToConsoleString(le)
	if le.IsError() {
		controller.errorLogger.Println(data)
	} else {
		controller.infoLogger.Println(data)
	}

	return nil
}

func (controller *Console) SendError(err error) error {
	return controller.Send(other.InterfaceToLogWrapper(err))
}

func (controller *Console) SendBatch(batch *types.LoggingQueueNode) error {
	for batch != nil {
		err := controller.Send(batch.Log)
		if err != nil {
			return err
		}
		batch = batch.Next
	}
	return nil
}

func (controller *Console) GetName() string {
	return controller.name
}

func (controller *Console) GetType() string {
	return "ConsoleLogger"
}
