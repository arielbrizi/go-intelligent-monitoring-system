package domain

//Notification is the domain struct to represent users notifications
type Notification struct {
	Type string

	Topic string

	Image Image

	Message string
}
