package recognitionapplicationportin

//QueueImagePort ...
type QueueImagePort interface {
	ReceiveImagesFromQueue() error
}
