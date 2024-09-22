package shell

type Scan struct {
	ScanFolder string `json:"scan_folder"`
	TaskId     string `json:"task_id"`
	Scanner    string `json:"scanner"`
}

func CreateNewScan(scanFolder string, taskId string, scanner string) Scan {
	return Scan{
		ScanFolder: scanFolder,
		TaskId:     taskId,
		Scanner:    scanner,
	}
}
