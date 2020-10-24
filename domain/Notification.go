package domain

//Notification is the domain struct to represent users notifications
type Notification struct {
	Type string

	Topic string

	Image Image

	Message string
}

//SMSNotification is the domain struct to represent SMS users notifications
type SMSNotification struct {
	Notification

	DestinationNumber string
}

//EmailNotification is the domain struct to represent e-Mail users notifications
type EmailNotification struct {
	Notification

	FromAdress string

	DestinationAdress string

	Subject string
}
