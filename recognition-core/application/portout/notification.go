package recognitionapplicationportout

import "go-intelligent-monitoring-system/domain"

//NotificationPort ...
type NotificationPort interface {
	NotifyTopic(notification domain.Notification) error
	NotifySMS(notification domain.SMSNotification) error
	NotifyEmail(notification domain.EmailNotification) error
}
