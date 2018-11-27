package models

import (
	"strconv"
	"time"
)

//Post represents an post on a tenant board
type Post struct {
	ID            int           `json:"id"`
	Number        int           `json:"number"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Description   string        `json:"description"`
	CreatedAt     time.Time     `json:"createdAt"`
	User          *User         `json:"user"`
	HasVoted      bool          `json:"hasVoted"`
	VotesCount    int           `json:"votesCount"`
	CommentsCount int           `json:"commentsCount"`
	Status        PostStatus    `json:"status"`
	Response      *PostResponse `json:"response,omitempty"`
	Tags          []string      `json:"tags"`
}

//Vote represents a vote given by a user on a post
type Vote struct {
	User *struct {
		ID    int    `json:"id" db:"id"`
		Name  string `json:"name" db:"name"`
		Email string `json:"email,omitempty" db:"email"`
	} `json:"user" db:"user"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// CanBeVoted returns true if this post can have its vote changed
func (i *Post) CanBeVoted() bool {
	return i.Status != PostCompleted && i.Status != PostDeclined && i.Status != PostDuplicate
}

// NewPost represents a new post
type NewPost struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdatePost represents a request to edit an existing post
type UpdatePost struct {
	Number      int    `route:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// DeletePost represents a request to delete an existing post
type DeletePost struct {
	Number int    `route:"number"`
	Text   string `json:"text"`
}

// NewComment represents a new comment
type NewComment struct {
	Number  int    `route:"number"`
	Content string `json:"content"`
}

// EditComment represents a request to edit existing comment
type EditComment struct {
	PostNumber int    `route:"number"`
	ID         int    `route:"id"`
	Content    string `json:"content"`
}

// DeleteComment represents a request to delete an existing comment
type DeleteComment struct {
	PostNumber int `route:"number"`
	CommentID  int `route:"id"`
}

// SetResponse represents the action to update an post response
type SetResponse struct {
	Number         int        `route:"number"`
	Status         PostStatus `json:"status"`
	Text           string     `json:"text"`
	OriginalNumber int        `json:"originalNumber"`
}

//PostResponse is a staff response to a given post
type PostResponse struct {
	Text        string        `json:"text"`
	RespondedAt time.Time     `json:"respondedAt"`
	User        *User         `json:"user"`
	Original    *OriginalPost `json:"original"`
}

//OriginalPost holds details of the original post of a duplicate
type OriginalPost struct {
	Number int        `json:"number"`
	Title  string     `json:"title"`
	Slug   string     `json:"slug"`
	Status PostStatus `json:"status"`
}

//Comment represents an user comment on an post
type Comment struct {
	ID        int        `json:"id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	User      *User      `json:"user"`
	EditedAt  *time.Time `json:"editedAt,omitempty"`
	EditedBy  *User      `json:"editedBy,omitempty"`
}

//Tag represents a simple tag
type Tag struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Color    string `json:"color"`
	IsPublic bool   `json:"isPublic"`
}

//CreateEditTag is used to create a new tag or edit existing
type CreateEditTag struct {
	Slug     string `route:"slug"`
	Name     string `json:"name"`
	Color    string `json:"color" format:"upper"`
	IsPublic bool   `json:"isPublic"`
}

// DeleteTag is used to delete an existing tag
type DeleteTag struct {
	Slug string `route:"slug"`
}

// AssignUnassignTag is used to assign or remove a tag to/from an post
type AssignUnassignTag struct {
	Slug   string `route:"slug"`
	Number int    `route:"number"`
}

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
	return "Unknown"
}

var (
	//SubscriberInactive means that the user cancelled the subscription
	SubscriberInactive = 0
	//SubscriberActive means that the subscription is active
	SubscriberActive = 1
)

//NotificationChannel represents the medium that the notification is sent
type NotificationChannel int

var (
	//NotificationChannelWeb is a in-app notification
	NotificationChannelWeb NotificationChannel = 1
	//NotificationChannelEmail is an email notification
	NotificationChannelEmail NotificationChannel = 2
)

//NotificationEvent represents all possible notification events
type NotificationEvent struct {
	UserSettingsKeyName           string
	DefaultSettingValue           string
	RequiresSubscriptionUserRoles []Role
	DefaultEnabledUserRoles       []Role
	Validate                      func(string) bool
}

func notificationEventValidation(v string) bool {
	return v == "0" || v == "1" || v == "2" || v == "3"
}

var (
	//NotificationEventNewPost is triggered when a new post is posted
	NotificationEventNewPost = NotificationEvent{
		UserSettingsKeyName:           "event_notification_new_post",
		DefaultSettingValue:           strconv.Itoa(int(NotificationChannelWeb | NotificationChannelEmail)),
		RequiresSubscriptionUserRoles: []Role{},
		DefaultEnabledUserRoles: []Role{
			RoleAdministrator,
			RoleCollaborator,
		},
		Validate: notificationEventValidation,
	}
	//NotificationEventNewComment is triggered when a new comment is posted
	NotificationEventNewComment = NotificationEvent{
		UserSettingsKeyName: "event_notification_new_comment",
		DefaultSettingValue: strconv.Itoa(int(NotificationChannelWeb | NotificationChannelEmail)),
		RequiresSubscriptionUserRoles: []Role{
			RoleVisitor,
		},
		DefaultEnabledUserRoles: []Role{
			RoleAdministrator,
			RoleCollaborator,
			RoleVisitor,
		},
		Validate: notificationEventValidation,
	}
	//NotificationEventChangeStatus is triggered when a new post has its status changed
	NotificationEventChangeStatus = NotificationEvent{
		UserSettingsKeyName: "event_notification_change_status",
		DefaultSettingValue: strconv.Itoa(int(NotificationChannelWeb | NotificationChannelEmail)),
		RequiresSubscriptionUserRoles: []Role{
			RoleVisitor,
		},
		DefaultEnabledUserRoles: []Role{
			RoleAdministrator,
			RoleCollaborator,
			RoleVisitor,
		},
		Validate: notificationEventValidation,
	}
	//AllNotificationEvents contains all possible notification events
	AllNotificationEvents = []NotificationEvent{
		NotificationEventNewPost,
		NotificationEventNewComment,
		NotificationEventChangeStatus,
	}
)
