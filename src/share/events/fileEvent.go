package events

import eventCenter "KeyMouseSimulation/common/Event"

// FileScanNewFile 扫描到新文件
const FileScanNewFile eventCenter.Topic = "file_scan_new_file"

type FileScanNewFileData struct {
	NewFile  []string
	FileList []string
}
