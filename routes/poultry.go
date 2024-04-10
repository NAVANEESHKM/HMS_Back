package routes

import (
	"backend/controllers"
	"crypto/tls"
	"log"
	"net/http"
	"net/smtp"

	//  "net/http"
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)
var check = true

func PoultryRoute(router *gin.Engine, controller controllers.HMSController) {
	router.POST("/post/doctor", controller.CreateDoctor)
	router.GET("/doctor/get", controller.GetDoctor)
	router.POST("patient/details", controller.PatientDetails)
	router.POST("/doctor/login", controller.DoctorLogin)
	router.POST("/doctor/clients", controller.DoctorClient)
	router.POST("/doctor/availability", controller.DoctorAvailability)
	router.POST("/patient/detailsupdate", controller.PatientDetailsUpdated)
	router.GET("/patientlist", controller.ListPatient)
	router.POST("/addmedicine", controller.AddMedicine)
	router.GET("/medicines", controller.GetMedicine)
	router.POST("/prescription", controller.AddPatinetPrescription)
	router.POST("/sendmail", SendMail)
}

func SendMail(ctx *gin.Context) {
	type Details struct {
		Emailid     string `json:"emailid" bson:"emailid"`
		PatientName string `json:"patientname" bson:"patientname"`
		DoctorName  string `json:"doctorname" bson:"doctorname"`
		Time        string `json:"time" bson:"time"`
		Date        string `json:"date" bson:"date"`
	}
	var details Details
	if err := ctx.ShouldBindJSON(&details); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sender := "navaneeshmuthusamy@gmail.com"
	password := "//need app password"
	recipient := details.Emailid
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(
		"To: " + recipient + "\r\n" +
			"Subject: Aquara Hospital!\r\n" +
			"\r\n" +
			"Dear " + details.PatientName + ",\r\n" +
			"\r\n" +
			"Your Appointmnet for Treatment from " + details.DoctorName + " has been Verified." + ",\r\n" +
			"\r\n" +
			"Time - " + details.Time + ",\r\n" +
			"\r\n" +
			"Date - " + details.Date + ",\r\n")

	auth := smtp.PlainAuth("", sender, password, smtpHost)

	client, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to connect to SMTP server"})
		log.Fatal(err)
		return
	}
	defer client.Close()

	if err := client.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
		log.Fatal(err)
	}

	if err := client.Auth(auth); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to authenticate"})
		log.Fatal(err)
		return
	}

	if err := client.Mail(sender); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to set sender"})
		log.Fatal(err)
		return
	}
	if err := client.Rcpt(recipient); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to set recipient"})
		log.Fatal(err)
		return
	}

	// Send the email content
	w, err := client.Data()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to send email content"})
		log.Fatal(err)
		return
	}
	_, err = w.Write(message)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to write message"})
		log.Fatal(err)
		return
	}
	err = w.Close()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to close email content"})
		log.Fatal(err)
		return
	}

	log.Println("Email sent successfully!")
	ctx.JSON(200, gin.H{"message": "Email sent successfully!"})
}
