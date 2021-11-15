package enum

//PostStatus is the status of a given post
type PostStatus int

var (
	//PostOpen is the default status
	PostOpen PostStatus
	//PostStarted is used when the post has been accepted and work is in progress
	PostStarted PostStatus = 1
	//PostCompleted is used when the post has been accepted and already implemented
	PostCompleted PostStatus = 2
	//PostDeclined is used when organizers decide to decline an post
	PostDeclined PostStatus = 3
	//PostPlanned is used when organizers have accepted an post and it's on the roadmap
	PostPlanned PostStatus = 4
	//PostDuplicate is used when the post has already been posted before
	PostDuplicate PostStatus = 5
	//PostDeleted is used when the post is completely removed from the site and should never be shown again
	PostDeleted PostStatus = 6
)
var postStatusIDs = map[PostStatus]string{
	PostOpen:      "open",
	PostStarted:   "started",
	PostCompleted: "completed",
	PostDeclined:  "declined",
	PostPlanned:   "planned",
	PostDuplicate: "duplicate",
	PostDeleted:   "deleted",
}

var postStatusNames = map[string]PostStatus{
	"open":      PostOpen,
	"started":   PostStarted,
	"completed": PostCompleted,
	"declined":  PostDeclined,
	"planned":   PostPlanned,
	"duplicate": PostDuplicate,
	"deleted":   PostDeleted,
}

// MarshalText returns the Text version of the post status
func (status PostStatus) MarshalText() ([]byte, error) {
	return []byte(postStatusIDs[status]), nil
}

// UnmarshalText parse string into a post status
func (status *PostStatus) UnmarshalText(text []byte) error {
	*status = postStatusNames[string(text)]
	return nil
}

// Name returns the name of a post status
func (status PostStatus) Name() string {
	name, ok := postStatusIDs[status]
	if ok {
		return name
	}
	return "unknown"
}
