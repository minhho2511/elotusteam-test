package transforms

type FileReq struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FilePath     string `json:"file_path"`
	ContentType  string `json:"content_type"`
	FileSize     int64  `json:"file_size"`
	UserAgent    string `json:"user_agent"`
	IPAddress    string `json:"ip_address"`
	Referer      string `json:"referer"`
}
