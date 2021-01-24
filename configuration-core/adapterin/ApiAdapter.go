package configurationadapterin

import (
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	"go-intelligent-monitoring-system/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//APIAdapter represents the adapter wich adds to the monitoring system  all authorized faces using the API
type APIAdapter struct {
	faceIndexerService configurationapplicationportin.ConfigurationPort
}

//AddAuthorizedFaceHandler ...
func (da *APIAdapter) AddAuthorizedFaceHandler(c *gin.Context) {

	var imageRequest domain.Image

	if errJSON := c.ShouldBindJSON(&imageRequest); errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errJSON.Error()})
		log.WithFields(log.Fields{"imageRequest.Name": imageRequest.Name, "imageRequest.Bucket": imageRequest.Bucket}).WithError(errJSON).Error("Error binding request JSON")
		return
	}

	log.WithFields(log.Fields{"imageRequest.Name": imageRequest.Name, "imageRequest.Bucket": imageRequest.Bucket}).Info("Adding face")

	//The first parameter is nil, so the image will be found on bucket.
	authorizedFace, err := da.faceIndexerService.AddAuthorizedFace(nil, imageRequest.Name, imageRequest.Bucket, imageRequest.CollectionName)
	if err != nil {
		//TODO: Analize err and return known errors codes. Ex: 404 - imageRequest.Name doesn't exist
		log.WithFields(log.Fields{"imageRequest.Name": imageRequest.Name, "imageRequest.Bucket": imageRequest.Bucket}).WithError(err).Error("Adding authorized face")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authorizedFace": authorizedFace})

}

//DeleteAuthorizedFaceHandler ...
func (da *APIAdapter) DeleteAuthorizedFaceHandler(c *gin.Context) {

	var authorizedFaceRequest domain.AuthorizedFace

	if errJSON := c.ShouldBindJSON(&authorizedFaceRequest); errJSON != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errJSON.Error()})
		log.WithFields(log.Fields{"authorizedFaceRequest.ID": authorizedFaceRequest.ID}).WithError(errJSON).Error("Error binding request JSON")
		return
	}

	log.WithFields(log.Fields{"authorizedFaceRequest.ID": authorizedFaceRequest.ID}).Info("Deleting Authorized Face")

	err := da.faceIndexerService.DeleteAuthorizedFace(authorizedFaceRequest)
	if err != nil {
		//TODO: Analize err and return known errors codes. Ex: 404 - imageRequest.Name doesn't exist
		log.WithFields(log.Fields{"authorizedFaceRequest.ID": authorizedFaceRequest.ID}).WithError(err).Error("Deleting authorized face")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authorizedFace": "deleted"})

}

//GetAuthorizedFacesHandler ...
func (da *APIAdapter) GetAuthorizedFacesHandler(c *gin.Context) {

	valuesMap := c.Request.URL.Query()

	if len(valuesMap["collectionName"]) < 1 {
		log.Error("Missing parameter collectionName")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter collectionName"})
		return
	}

	authorizedFaces, err := da.faceIndexerService.GetAuthorizedFaces(valuesMap["collectionName"][0])
	if err != nil {
		//TODO: Analize err and return known errors codes. Ex: 404 - imageRequest.Name doesn't exist
		log.WithFields(log.Fields{"collectionName": valuesMap["collectionName"][0]}).WithError(err).Error("Error getting authorized faces")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authorizedFaces": authorizedFaces})
}

//NewAPIAdapter initializes an NewApiAdapter object.
func NewAPIAdapter(faceIndexerService configurationapplicationportin.ConfigurationPort) *APIAdapter {
	return &APIAdapter{
		faceIndexerService: faceIndexerService,
	}
}