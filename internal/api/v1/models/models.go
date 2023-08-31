package api

type AddSegmentRequest struct {
	Slug    string  `json:"slug" binding:"required"`
	Percent float32 `json:"percent" binding:"omitempty,gte=1,lte=100"`
}

type DeleteSegmentRequest struct {
	Slug string `json:"slug" binding:"required"`
}

type AddUserSegmentRequest struct {
	NewSegments []Segment `json:"new_segments" binding:"required"`
	OldSegments []Segment `json:"old_segments" binding:"required"`
}

type Segment struct {
	Name string `json:"name"`
	Ttl  int64  `json:"ttl"`
}
