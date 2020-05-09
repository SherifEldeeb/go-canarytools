package canarytools

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileForwader forwards incidents to file, it uses lumberjack
// https://github.com/natefinch/lumberjack
// Logger opens or creates the logfile on first Write. If the file exists and is
// less than MaxSize megabytes, lumberjack will open and append to that file. If
// the file exists and its size is >= MaxSize megabytes, the file is renamed by
// putting the current time in a timestamp in the name immediately before the
// file's extension (or the end of the filename if there's no extension). A new
// log file is then created using original filename.
type FileForwader struct {
	Filename string
	l        *log.Logger
	lj       *lumberjack.Logger
}

// NewFileForwader creates a new FileForwarder
func NewFileForwader(filename string, maxsize, maxbackups, maxage int, compress bool, l *log.Logger) (fileforwader *FileForwader, err error) {
	fileforwader = &FileForwader{}
	fileforwader.l = l
	if filename == "" {
		return nil, errors.New("filename can't be empty")
	}
	fileforwader.lj = &lumberjack.Logger{}
	fileforwader.lj.Filename = filename
	fileforwader.lj.Compress = compress
	fileforwader.lj.MaxAge = maxage
	fileforwader.lj.MaxBackups = maxbackups
	fileforwader.lj.MaxSize = maxsize

	return
}

// Forward writes incidents to file
func (f FileForwader) Forward(outChan <-chan []byte) {
	for v := range outChan {
		n, err := f.lj.Write(v)
		if err != nil {
			f.l.WithFields(log.Fields{
				"source": "FileForwarder",
				"stage":  "forward",
				"err":    err,
			}).Error("FileOut error writing incident")
			continue
		}
		f.l.WithFields(log.Fields{
			"source": "FileForwarder",
			"stage":  "forward",
			"bytes":  n,
		}).Info("FileOut incident written")
	}
}
