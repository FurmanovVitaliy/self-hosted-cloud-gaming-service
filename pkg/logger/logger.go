package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

const pathToLogDir = "/home/vitalii/Шаблоны/Rooms"
const fileRights = 0770

type Logger struct {
	*logrus.Entry
}

var instance Logger
var once sync.Once

func Init(level string) Logger {
	once.Do(func() {
		logrusLevel, err := logrus.ParseLevel(level)
		if err != nil {
			log.Fatalln(err)
		}

		l := logrus.New()
		l.SetReportCaller(true)
		l.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("[%s:%d]", filename, f.Line), fmt.Sprintf("%s()", f.Function)
			},
			DisableColors: false,
			FullTimestamp: true,
		}
		l.SetOutput(os.Stdout)
		l.SetLevel(logrusLevel)
		instance = Logger{logrus.NewEntry(l)}
	})
	return instance
}
