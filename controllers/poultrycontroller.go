package controllers

import (
	"backend/interfaces"
	"backend/models"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
)

type HMSController struct {
	HMSService interfaces.IHMS
}

func InitPoultryController(HMSService interfaces.IHMS) HMSController {
	return HMSController{HMSService} //DI(dependency injection) pattern
}

func (pc *HMSController) CreateDoctor(ctx *gin.Context) {

	// Get the file from the form data

	file, _, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form"})
		return
	}
	defer file.Close()

	// Read the content of the file into a byte slice
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while Reading file"})
		return
	}

	// Encode the file content to Base64
	base64Encoded := base64.StdEncoding.EncodeToString(fileContent)

	// Print or store the Base64-encoded string
	var customer models.DoctorDetails
	customer.Image = "data:image/png;base64," + base64Encoded
	file1 := ctx.Request.FormValue("name")
	customer.Name = file1
	file2 := ctx.Request.FormValue("emailid")
	customer.EmailID = file2
	file3 := ctx.Request.FormValue("qualification")
	customer.Qualification = file3
	file4 := ctx.Request.FormValue("specialist")
	customer.Specialist = file4
	file5 := ctx.Request.FormValue("phonenumber")
	customer.PhoneNumber = file5

	log.Println("Customer.....", customer)

	// if err := ctx.ShouldBindJSON(&customer); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, err.Error())
	// 	return
	// }
	if customer.Patients == nil {
		customer.Patients = []models.Patient{}
		customer.Availability = true
		patients := models.Patient{
			Precaution: []string{},
		}
		customer.Patients = append(customer.Patients, patients)
	}

	val, err := pc.HMSService.CreateDoctor(&customer)
	if val == "greater" {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "greater"})
		return
	} else if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

}

func (pc *HMSController) GetDoctor(ctx *gin.Context) {

	val, err := pc.HMSService.GetDoctor(ctx)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "message": val})
		return
	}

}

func (pc *HMSController) PatientDetails(ctx *gin.Context) {

	var customer *models.Patient
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	customer.TimeNow = time.Now().Format("2006-01-02")

	err := pc.HMSService.PatientDetails(customer)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

}

func (pc *HMSController) DoctorLogin(ctx *gin.Context) {
	var value *models.Login
	if err := ctx.ShouldBindJSON(&value); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	doctor, result, err := pc.HMSService.DoctorLogin(value)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Fail", "message": result})
		return
	} else {
		fmt.Println("NUMBER..........", doctor)
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		// Twilio Account SID and Auth Token
		// Read Twilio credentials from environment variables
		accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
		authToken := os.Getenv("TWILIO_AUTH_TOKEN")

		// Create a Twilio client
		client := gotwilio.NewTwilioClient(accountSid, authToken)

		// Sender's phone number (Twilio number)
		from := os.Getenv("TWILIO_FROM")

		// Recipient's phone number
		to := "+91" + doctor

		// Message to send
		opt_store := GenerateOTP()
		message := opt_store
		// Send the message
		resp, ex, err := client.SendSMS(from, to, message, "", "")
		if err != nil {
			log.Fatalf("Error sending SMS: %s", err)
		}
		if ex != nil {
			log.Fatalf("Exception sending SMS: %s", ex.Message)
		}
		log.Printf("SMS Sent! SID: %s, Status: %s", resp.Sid, resp.Status)
		ctx.JSON(http.StatusAccepted, gin.H{"status": "Success", "message": result, "otp": opt_store})
		return
	}
}
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return fmt.Sprintf("%06d", rand.Intn(max-min+1)+min)
}

func (pc *HMSController) DoctorClient(ctx *gin.Context) {
	var value *models.Login
	if err := ctx.ShouldBindJSON(&value); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := pc.HMSService.DoctorClient(value)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Fail", "message": result})
		return
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"status": "Success", "message": result})
		return
	}
}

func (pc *HMSController) DoctorAvailability(ctx *gin.Context) {
	var value *models.UpdateStatus

	if err := ctx.ShouldBindJSON(&value); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := pc.HMSService.DoctorAvailability(value)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": result, "Error": err})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": result, "Error": "None"})

}

func (pc *HMSController) PatientDetailsUpdated(ctx *gin.Context) {
	var value models.PatientDetails
	if err := ctx.ShouldBindJSON(&value); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	val, err := pc.HMSService.PatientDetailsUpdated(&value)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "message": val})
		return
	}
}

func (pc *HMSController) ListPatient(ctx *gin.Context) {
	val, err := pc.HMSService.ListPatient(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "message": val})
		return
	}

}

func (pc *HMSController) AddMedicine(ctx *gin.Context) {
	var customer *models.MedicineDetails
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err := pc.HMSService.AddMedicine(customer)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
}

func (pc *HMSController) GetMedicine(ctx *gin.Context) {

	val, err := pc.HMSService.GetMedicine(ctx)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "message": val})
		return
	}

}
func (pc *HMSController) AddPatinetPrescription(ctx *gin.Context) {
	var Precaution *models.PrecautionUpdate
	if err := ctx.ShouldBindJSON(&Precaution); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err := pc.HMSService.AddPatinetPrescription(Precaution)
	if err != nil {

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"status": "success"})
		return
	}

}
