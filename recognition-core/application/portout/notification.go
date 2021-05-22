package recognitionapplicationportout

import "go-intelligent-monitoring-system/domain"

//NotificationPort ...
type NotificationPort interface {
	NotifyUnauthorizedFace(notification domain.Notification) error
	NotifyInitializedSystem() error
}
