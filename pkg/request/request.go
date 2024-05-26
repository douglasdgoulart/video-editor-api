package request

type Input struct {
	FileURL          string `json:"file_url,omitempty"`
	UploadedFilePath string `json:"uploaded_file_path,omitempty"`
}

type Output struct {
	FilePattern string `json:"file_pattern,omitempty" required:"true"`
	WebhookURL  string `json:"webhook_url,omitempty"`
}

type EditorRequest struct {
	Input        Input             `json:"input,omitempty"`
	Output       Output            `json:"output" required:"true"`
	Codec        string            `json:"codec,omitempty"`
	Bitrate      string            `json:"bitrate,omitempty"`
	Resolution   string            `json:"resolution,omitempty"`
	AudioCodec   string            `json:"audio_codec,omitempty"`
	AudioBitrate string            `json:"audio_bitrate,omitempty"`
	Filters      map[string]string `json:"filters,omitempty"`
	ExtraOptions string            `json:"extra_options,omitempty"`
	StartTime    string            `json:"start_time,omitempty"`
	Frames       string            `json:"frames,omitempty"`
}
