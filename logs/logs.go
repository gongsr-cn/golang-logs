package Logs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	LOGS  = " [logs] "
	DEBUG = " [debug]"
	INFO  = " [info] "
	WARN  = " [warn] "
	ERROR = " [error]"
	PANIC = " [panic]"
)
var (
	logsF *os.File
	FileName = "storage"
	SerialNumber uint8 = 1
	MaxSize uint32  = 100 * (1 << 20)
)

type Logs interface {
	Debug(msg string) error
	Info(msg string) error
	Warn(msg string) error
	Error(msg string) error
}

type logs struct {
	LogPath string   `json:"log_path"` // log file save path
	osFile  *os.File `json:"os_file"`  // system interface logs's log file
	file    *file `json:"file"`
}

type file struct {
	mu           sync.Mutex 			//
	size     	 uint32	`json:"size"`	// file size, more than 100M to create a log file
	config 		 *Config `json:"config"`	// config pointer
	osFile   	 *os.File `json:"os_file"`	// system interface
}

func NewLogs(lPath string, c *Config) (*logs, error) {
	serialNumber, size := checkDirectory(lPath)

	// create log file at log self
	var err error
	logsF, err = os.OpenFile(lPath+"/logs.log", os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		logsF, _ = createFile(lPath + "/logs")
	}

	c.verify()
	c.LogsPath = lPath
	c.SerialNumber = serialNumber

	f, err := newFile(size, c)
	if err != nil {
		output(ERROR, err.Error(), logsF)
		return nil, err
	}
	return &logs{
		osFile:  logsF,
		LogPath: lPath,
		file:    f,
	}, nil
}

func (l *logs) Debug(msg string) error {
	if err := l.output(DEBUG, msg); err != nil {
		return err
	}
	return nil
}

func (l *logs) Info(msg string) error {
	if err := l.output(INFO, msg); err != nil {
		return err
	}
	return nil
}

func (l *logs) Warn(msg string) error {
	if err := l.output(WARN, msg); err != nil {
		return err
	}
	return nil
}

func (l *logs) Error(msg string) error {
	if err := l.output(ERROR, msg); err != nil {
		return err
	}
	return nil
}

// output
// ----------------------------------------------------------------
// fprintf information to log file
// ----------------------------------------------------------------
func (l *logs) output(types, str string) error {
	size := uint32(len(logFormat(types, str)))
	l.file.mu.Lock()
	defer l.file.mu.Unlock()
	if l.file.checkFileSize(size) {
		output(types, str, l.file.osFile)
	}
	return errors.New(str)
}

func output(types, str string, f *os.File) {
	msg := logFormat(types, str)
	fmt.Fprintf(f, msg)
}

// newFile
// ----------------------------------------------------------------
// file struct, return this pointer
// ----------------------------------------------------------------
func newFile(size uint32, c *Config) (*file, error) {
	// Combined file name, read or open this file
	fileName := c.LogsPath + "/" + c.FileName + strconv.Itoa(int(c.SerialNumber))
	f, err := os.OpenFile(fileName+".log", os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		f, err = createFile(fileName)
		if err != nil {
			output(ERROR, err.Error(), logsF)
			return nil, err
		}
	}
	return &file{
		mu:		sync.Mutex{},
		size:	size,
		config: c,
		osFile:	f,
	}, nil
}

// checkFileSize
// ----------------------------------------------------------------
// check the file size, if more than this max size
// create other file
// ----------------------------------------------------------------
func (f *file) checkFileSize(n uint32) bool {
	var err error
	// more than 100M
	if f.size+n > f.config.MaxSize {
		f.config.SerialNumber++
		fileName := f.config.LogsPath + "/" + f.config.FileName + strconv.Itoa(int(f.config.SerialNumber))
		f.osFile, err = createFile(fileName)
		if err != nil {
			f.config.SerialNumber--
			output(ERROR, err.Error(), logsF)
			return false
		}
		f.size = n
		return true
	}
	//
	f.size = f.size + n
	return true
}

// checkDirectory
// ----------------------------------------------------------------
// check the directory path, check log file list
// ----------------------------------------------------------------
func checkDirectory(dir string) (uint8, uint32) {
	var serialNumber = uint8(1)
	var size = uint32(0)
	// read directory, return the directory's file list
	fileList, err := os.ReadDir(dir)
	if err != nil {
		err = createDir(dir)
		if err != nil {
			output(ERROR, err.Error(), logsF)
		}
		return serialNumber, size
	}
	for k, v := range fileList {
		// select current file serial number, print to this file
		if strings.Contains(v.Name(), FileName) {
			serialNumber++
		}
		// compare the size of two files
		if k == len(fileList)-1 {
			f, _ := os.Stat(dir + "/" + v.Name())
			size = uint32(f.Size())
			if size < MaxSize {
				serialNumber--
			}
		}
	}
	return serialNumber, size
}

// createDir
// ----------------------------------------------------------------
// if the directory is not exists, create the directory
// ----------------------------------------------------------------
func createDir(dir string) error {
	if !strings.Contains(dir, "/") {
		return os.Mkdir(dir, 0777)
	}
	return os.MkdirAll(dir, 0777)
	// dirList := strings.Split(dir, "/")
	// dirTmp := ""
	// for _, v := range dirList {
	// 	if v == "." {
	// 		dirTmp = dirTmp + "."
	// 		continue
	// 	}
	// 	dirTmp = dirTmp + "/" + v
	// 	if _, err := os.ReadDir(dirTmp); err != nil {
	// 		os.Mkdir(dirTmp, 0777)
	// 	}
	// }
	// err := os.Mkdir(dirTmp, 0777)
	// return err
}

// createFile
// ----------------------------------------------------------------
// if the file is not exists, create the file
// ----------------------------------------------------------------
func createFile(filename string) (*os.File, error) {
	output(INFO, fmt.Sprintf("create %s success\n", filename + ".log"), logsF)
	return os.Create(filename + ".log")
}

// logFormat
// ----------------------------------------------------------------
// log information formatting
// ----------------------------------------------------------------
func logFormat(lType, str string) string {
	msg := time.Now().Format("2006-01-02 15:04:05.000000")
	msg = msg +  lType + " " + str + "\n"
	return msg
}
