package Logs

type Config struct {
	SerialNumber uint8  `json:"serial_number"` // the number of the current file
	MaxSize      uint32 `json:"max_size"`      // 100M = 100*1024*1024 byte
	Size         uint32 `json:"size"`          // file size, more than 100M to create a log file
	FileName     string `json:"file_name"`     // the name of the current file, only the name has no suffix
	LogsPath     string `json:"logs_path"`     // log file save directory
}

func (c *Config) verify() {
	if c.SerialNumber == 0 {
		c.SerialNumber = SerialNumber
	}
	if c.MaxSize == 0 {
		c.MaxSize = MaxSize
	}
	if c.FileName == "" {
		c.FileName = FileName
	}
}
