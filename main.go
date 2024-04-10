package main

import (
	"context"
	"fmt"

	"backend/config"
	"backend/constants"
	"backend/controllers"
	"backend/routes"
	"backend/services"
	"log"

	"github.com/gin-contrib/cors"

	//	"rest-api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoclient *mongo.Client
	ctx         context.Context
	server      *gin.Engine
)

func initApp(mongoClient *mongo.Client) {
	//Customer Collection
	ctx = context.TODO()
	AdminCollection := mongoClient.Database(constants.DatabaseName).Collection("admin")
	DoctorCollection := mongoClient.Database(constants.DatabaseName).Collection("doctor")
	PatientCollection := mongoClient.Database(constants.DatabaseName).Collection("patient")
	MedicineCollection := mongoClient.Database(constants.DatabaseName).Collection("medicine")

	HMSService := services.PoultryServiceInit(AdminCollection, DoctorCollection, PatientCollection, MedicineCollection, ctx)
	PoultryController := controllers.InitPoultryController(HMSService)
	routes.PoultryRoute(server, PoultryController)

}

// https://poultry-front.vercel.app
// http://localhost:3000/
func main() {

	server = gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Allow any origin
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	mongoclient, err := config.ConnectDataBase()
	defer mongoclient.Disconnect(ctx)
	if err != nil {
		panic(err)
	}

	initApp(mongoclient)
	fmt.Println("server running on port", constants.Port)
	log.Fatal(server.Run(constants.Port))
	// log.Fatal(server.RunTLS(":4000", sslCertCrtPath, sslCertKeyPath))

}
