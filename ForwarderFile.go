package canarytools

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// FileForwader
type FileForwader struct {
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int  //days
	Compress   bool // disabled by default
	l          *log.Logger
}

func NewFileForwader(filename, loglevel string, maxsize, maxbackups, maxage int, compress bool) (fileforwader *FileForwader, err error) {
	fileforwader = &FileForwader{}
	// logging config
	fileforwader.l = log.New()
	switch loglevel {
	case "info":
		fileforwader.l.SetLevel(log.InfoLevel)
	case "warning":
		fileforwader.l.SetLevel(log.WarnLevel)
	case "debug":
		fileforwader.l.SetLevel(log.DebugLevel)
	default:
		return nil, errors.New("unsupported log level (can be 'info', 'warning' or 'debug')")
	}

	if filename == "" {
		return nil, errors.New("filename can't be empty")
	}
	fileforwader.Compress = compress
	fileforwader.Filename = filename
	fileforwader.MaxAge = maxage
	fileforwader.MaxBackups = maxbackups
	fileforwader.MaxSize = maxsize
	return
}

func (f FileForwader) Forward(outChan <-chan Incident) {
	// connect := t.host + ":" + strconv.Itoa(t.port)
	// c, err := net.Dial("tcp", connect)
	// if err != nil {
	// 	return
	// }
	// defer c.Close()
	// enc := json.NewEncoder(c)
	// enc.SetEscapeHTML(false)

	// for i := range outChan {
	// 	err := enc.Encode(i)
	// 	if err != nil {
	// 		logrus.Errorf("encoding: %#v", i)
	// 	}
	// }
}
