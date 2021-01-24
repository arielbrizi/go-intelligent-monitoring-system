package configurationadapterin

import (
	configurationadapterout "go-intelligent-monitoring-system/configuration-core/adapterout"
	configurationapplication "go-intelligent-monitoring-system/configuration-core/application"
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

func TestAddAuthorizedFace(t *testing.T) {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapterTest()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(rekoAdapter)

	//"Adapter In": APIAdapter sets the authorized face from the request
	confAPIAdapter := NewAPIAdapter(faceIndexerService)

	router := gin.Default()
	router.POST("/configuration-core/authorized-face", confAPIAdapter.AddAuthorizedFaceHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/configuration-core/authorized-face", strings.NewReader("{\"name\": \"silvia3.jpg\",\"bucket\": \"camarasilvia\", \"collection\": \"camarasilvia\"}"))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"authorizedFace\":{\"name\":\"silvia3.jpg\",\"bucket\":\"camarasilvia\",\"collection\":\"camarasilvia\",\"id\":\"123\",\"Bytes\":null}}", w.Body.String())

	//Without name
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/configuration-core/authorized-face", strings.NewReader("{\"bucket\": \"camarasilvia\", \"collection\": \"camarasilvia\"}"))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	//Without bucket
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/configuration-core/authorized-face", strings.NewReader("{\"collection\": \"camarasilvia\", \"name\": \"silvia3.jpg\"}"))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	//Without collection
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/configuration-core/authorized-face", strings.NewReader("{\"name\": \"silvia3.jpg\", \"bucket\": \"camarasilvia\"}"))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestDeleteAddAuthorizedFace(t *testing.T) {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapterTest()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(rekoAdapter)

	//"Adapter In": APIAdapter sets the authorized face from the request
	confAPIAdapter := NewAPIAdapter(faceIndexerService)

	router := gin.Default()
	router.DELETE("/configuration-core/authorized-face", confAPIAdapter.DeleteAuthorizedFaceHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/configuration-core/authorized-face", strings.NewReader("{\"id\": \"95ec2bd7-b891-494f-b0cc-48e20c11ee5f\",\"collection\": \"camarasilvia\"}"))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"authorizedFace\":\"deleted\"}", w.Body.String())

	//without id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/configuration-core/authorized-face", strings.NewReader("{\"collection\": \"camarasilvia\"}"))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestGetAddAuthorizedFaces(t *testing.T) {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapterTest()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(rekoAdapter)

	//"Adapter In": APIAdapter sets the authorized face from the request
	confAPIAdapter := NewAPIAdapter(faceIndexerService)

	router := gin.Default()
	router.GET("/configuration-core/authorized-face", confAPIAdapter.GetAuthorizedFacesHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/configuration-core/authorized-face?collectionName=camarasilvia", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"authorizedFaces\":[{\"name\":\"silvia1.jpg\",\"bucket\":\"camarasilvia\",\"collection\":\"camarasilvia\",\"id\":\"2659022e-4ad2-4be8-81b9-1d4b1953ff90\",\"Bytes\":null}]}", w.Body.String())

	//without collectionName parameter
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/configuration-core/authorized-face", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}
