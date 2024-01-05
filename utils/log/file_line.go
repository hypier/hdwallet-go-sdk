package log

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
)

type fileHook struct{}

func newFileHook() *fileHook {
	return &fileHook{}
}

func (f *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *fileHook) Fire(entry *logrus.Entry) error {
	var s string
	_, b, c, _ := runtime.Caller(8)
	i := strings.LastIndex(b, "/")
	if i != -1 {
		s = b[i+1:len(b)] + ":" + IntToString(c)
	}
	entry.Data["FilePath"] = s
	return nil
}
func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
