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

func NewFileForwader(filename string, maxsize, maxbackups, maxage int, compress bool, l *log.Logger) (fileforwader *FileForwader, err error) {
	fileforwader = &FileForwader{}
	fileforwader.l = l

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
