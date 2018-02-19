package models

import (
	"strconv"
	"time"
)

//Idea represents an idea on a tenant board
type Idea struct {
	ID              int           `json:"id"`
	Number          int           `json:"number"`
	Title           string        `json:"title"`
	Slug            string        `json:"slug"`
	Description     string        `json:"description"`
	CreatedOn       time.Time     `json:"createdOn"`
	User            *User         `json:"user"`
	ViewerSupported bool          `json:"viewerSupported"`
	TotalSupporters int           `json:"totalSupporters"`
	TotalComments   int           `json:"totalComments"`
	Status          int           `json:"status"`
	Response        *IdeaResponse `json:"response"`
	Tags            []int64       `json:"tags"`
	Ranking         float64       `json:"ranking"`
}

// CanBeSupported returns true if this idea can be Supported/UnSupported
func (i *Idea) CanBeSupported() bool {
	return i.Status != IdeaCompleted && i.Status != IdeaDeclined && i.Status != IdeaDuplicate
}

//BasicIdea is a subset of Idea with few fields
type BasicIdea struct {
	ID              int    `json:"id"`
	Number          int    `json:"number"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	TotalSupporters int    `json:"totalSupporters"`
	Status          int    `json:"status"`
}

// NewIdea represents a new idea
type NewIdea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateIdea represents a request to edit an existing idea
type UpdateIdea struct {
	Number      int    `route:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// NewComment represents a new comment
type NewComment struct {
	Number  int    `route:"number"`
	Content string `json:"content"`
}

// SetResponse represents the action to update an idea response
type SetResponse struct {
	Number         int    `route:"number"`
	Status         int    `json:"status"`
	Text           string `json:"text"`
	OriginalNumber int    `json:"originalNumber"`
}

//IdeaResponse is a staff response to a given idea
type IdeaResponse struct {
	Text        string        `json:"text"`
	RespondedOn time.Time     `json:"respondedOn"`
	User        *User         `json:"user"`
	Original    *OriginalIdea `json:"original"`
}

//OriginalIdea holds details of the original idea of a duplicate
type OriginalIdea struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Status int    `json:"status"`
}

//Comment represents an user comment on an idea
type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedOn time.Time `json:"createdOn"`
	User      *User     `json:"user"`
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

// AssignUnassignTag is used to assign or remove a tag to/from an idea
type AssignUnassignTag struct {
	Slug   string `route:"slug"`
	Number int    `route:"number"`
}

var (
	//IdeaOpen is the default status
	IdeaOpen = 0
	//IdeaStarted is used when the idea has been accepted and work is in progress
	IdeaStarted = 1
	//IdeaCompleted is used when the idea has been accepted and already implemented
	IdeaCompleted = 2
	//IdeaDeclined is used when organizers decide to decline an idea
	IdeaDeclined = 3
	//IdeaPlanned is used when organizers have accepted an idea and it's on the roadmap
	IdeaPlanned = 4
	//IdeaDuplicate is used when the idea has already been posted before
	IdeaDuplicate = 5
)

// GetIdeaStatusName returns the name of an idea status
func GetIdeaStatusName(status int) string {
	switch status {
	case 0:
		return "Open"
	case 1:
		return "Started"
	case 2:
		return "Completed"
	case 3:
		return "Declined"
	case 4:
		return "Planned"
	case 5:
		return "Duplicate"
	default:
		return "Unknown"
	}
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
	UserSettingsKeyName          string
	DefaultSettingValue          string
	RequiresSubscripionUserRoles []Role
	DefaultEnabledUserRoles      []Role
	Validate                     func(string) bool
}

func notificationEventValidation(v string) bool {
	return v == "0" || v == "1" || v == "2" || v == "3"
}

var (
	//NotificationEventNewIdea is triggered when a new idea is posted
	NotificationEventNewIdea = NotificationEvent{
		UserSettingsKeyName:          "event_notification_new_idea",
		DefaultSettingValue:          strconv.Itoa(int(NotificationChannelWeb | NotificationChannelEmail)),
		RequiresSubscripionUserRoles: []Role{},
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
		RequiresSubscripionUserRoles: []Role{
			RoleVisitor,
		},
		DefaultEnabledUserRoles: []Role{
			RoleAdministrator,
			RoleCollaborator,
			RoleVisitor,
		},
		Validate: notificationEventValidation,
	}
	//NotificationEventChangeStatus is triggered when a new idea has its status changed
	NotificationEventChangeStatus = NotificationEvent{
		UserSettingsKeyName: "event_notification_change_status",
		DefaultSettingValue: strconv.Itoa(int(NotificationChannelWeb | NotificationChannelEmail)),
		RequiresSubscripionUserRoles: []Role{
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
		NotificationEventNewIdea,
		NotificationEventNewComment,
		NotificationEventChangeStatus,
	}
)
