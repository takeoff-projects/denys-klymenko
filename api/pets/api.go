package pets

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"takeoff-projects/denys-klymenko/core/pets"
	_ "takeoff-projects/denys-klymenko/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @query.collection.format multi

// @x-extension-openapi {"example": "value on a json format"}

func handleRequests(dal *PetsDAL) {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	h := Controller{dal: dal}

	r.GET("/pets", h.getPets)
	r.POST("/pets", h.createPet)
	r.GET("/pets/:id", h.getPet)
	r.DELETE("/pets/:id", h.deletePet)
	r.PUT("/pets/:id", h.updatePet)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Listening on port %s", port)
	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("fail to start http server %s", err)
	}
}

type Controller struct {
	dal *PetsDAL
}

// swagger godoc
// @Id swagger
// @Summary swagger
// @Description swagger console
// @Produce html
// @Success 200 {file} file
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @x-google-backend  {"address": "${backend}/swagger/index.html"}
// @Router /swagger/index.html [get]
func (h *Controller) swagger(c *gin.Context) {
	//stub for swagger
}

// listPets godoc
// @Id listPet
// @Summary list all pets
// @Description list all pets ordered by likes
// @Produce  json
// @Success 200 {object} []pets.Pet
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @x-google-backend  {"address": "${backend}/pets"}
// @Router /pets [get]
func (h *Controller) getPets(c *gin.Context) {
	ps, err := h.dal.List(context.TODO())
	if err != nil {
		c.String(http.StatusBadRequest, "fail to list pets %s", err)
	}

	c.JSON(http.StatusOK, struct {
		Pets []pets.Pet `json:"pets"`
	}{ps})
}

// getPet godoc
// @Id getPet
// @Summary get pet by id
// @Description get pet by id
// @Param id path string true "pet id"
// @Produce  json
// @Success 200 {object} pets.Pet
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @x-google-backend  {"address": "${backend}/pets"}
// @Router /pets/{id} [get]
func (h *Controller) getPet(c *gin.Context) {
	id := c.Param("id")

	pet, err := h.dal.Get(context.TODO(), id)
	if err != nil {
		c.String(http.StatusBadRequest, "Hello %s", err)
	}

	c.JSON(http.StatusOK, pet)
}

// deletePet godoc
// @Id deletePet
// @Summary delete pet by id
// @Description delete pet by id
// @Param id path string true "pet id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @x-google-backend  {"address": "${backend}/pets"}
// @Router /pets/{id} [delete]
func (h *Controller) deletePet(c *gin.Context) {
	id := c.Param("id")
	err := h.dal.Delete(context.TODO(), id)
	if err != nil {
		c.String(http.StatusBadRequest, "Hello %s", err)
	}
	c.String(http.StatusOK, "ok")
}

// updatePet godoc
// @Id updatePet
// @Summary update pet by id
// @Description update pet by id
// @Accept json
// @Param pet body pets.Update true "pet to update"
// @Param id path string true "pet id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @x-google-backend  {"address": "${backend}/pets"}
// @Router /pets/{id} [put]
func (h *Controller) updatePet(c *gin.Context) {
	id := c.Param("id")

	var update pets.Update
	err := c.BindJSON(&update)
	if err != nil {
		c.String(http.StatusBadRequest, "wrong format %s", err)
	}

	err = h.dal.Update(context.TODO(), id, update)
	if err != nil {
		c.String(http.StatusBadRequest, "Hello %s", err)
	}

	c.String(http.StatusOK, "ok")
}

// createPet godoc
// @Id createPet
// @Summary create a pet
// @Description create a pet
// @Accept json
// @Param pet body pets.Create true "pet to create"
// @Success 200 {string} string "ok"
// @Failure 400,404 {string} string
// @Failure 500 {string} string
// @Failure default {string} string
// @x-google-backend  {"address": "${backend}/pets"}
// @Router /pets [post]
func (h *Controller) createPet(c *gin.Context) {
	var pet pets.Create
	err := c.BindJSON(&pet)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid request body", err)
	}

	err = h.dal.Create(context.TODO(), pet)
	if err != nil {
		c.String(http.StatusBadRequest, "fail to create pet", err)
	}
	c.String(http.StatusCreated, "ok")
}

func Run() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID == "" {
		log.Fatalf("GOOGLE_CLOUD_PROJECT must be set")
		return
	}

	firestoreClient, err := firestore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	defer func() {
		_ = firestoreClient.Close()
	}()

	dal := PetsDAL{
		Client:    firestoreClient,
		ProjectID: projectID,
	}

	handleRequests(&dal)
}
