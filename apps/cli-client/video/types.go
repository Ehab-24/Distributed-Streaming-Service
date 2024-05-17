package video

type Resolution struct {
	width  int
	height int
}

type Bitrate int32

type VideoQaulity struct {
	resolution Resolution
	bitrate    Bitrate
}

type Duration struct {
	Hours   int
	Minutes int
	Seconds int
}

type UploadVideoPayload struct {
	title string `json:"title"`
}
