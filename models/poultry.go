package models

type DoctorDetails struct {
	EmailID       string `json:"emailid" bson:"emailid"`
	Name          string `json:"name" bson:"name"`
	Qualification string `json:"qualification" bson:"qualification"`
	Specialist    string `json:"specialist" bson:"specialist"`
	Availability  bool   `json:"availability" bson:"availability"`
    Image        string `json:"image"` 
	PhoneNumber string `json:"phonenumber" bson:"phonenumber"`
	Patients      []Patient
}

type Patient struct {
	PatientID string `json:"patientid" bson:"patientid"`
	Name      string `json:"name" bson:"name"`
	Email     string `json:"email" bson:"email"`
	Purpose   string `json:"purpose" bson:"purpose"`
	Precaution []string
	Doctor    string `json:"doctor" bson:"doctor"`
	Number    int64  `json:"number" bson:"number"`
	Date      string `json:"date" bson:"date"`
	Time      string `json:"time" bson:"time"`
	TimeNow   string `json:"timenow" bson:"timenow"`
}
type Precaution struct{
	Medicines string `json:"medicines" bson:"medicines"`
}

type PatientData struct {
	PatientID   string `json:"patientid" bson:"patientid"`
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	Number      int64  `json:"number" bson:"number"`
	DoctorNames []DoctorName
}

type DoctorName struct {
	DoctorName string `json:"doctorname" bson:"doctorname"`
	Purpose   string `json:"purpose" bson:"purpose"`
	Date string `json:"date" bson:"date"`
	Precaution []string 
}
type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UpdateStatus struct {
	Email        string `json:"email" bson:"email"`
	Availability bool   `json:"availability" bson:"availability"`
}
type PatientDetails struct {
	PatientID string `json:"patientid" bson:"patientid"`
}

type MedicineDetails struct {
	MedicineName string `json:"medicinename" bson:"medicinename"`
	Age string `json:"age" bson:"age"`
	Price  int `json:"price" bson:"price"`
    Description string `json:"description" bson:"description"` 
}
type PrecautionUpdate struct{
	DoctorName string `json:"doctorname" bson:"doctorname"`
	DoctorEmail string `json:"doctoremail" bson:"doctoremail"`
	PatientID string `json:"patientid" bson:"patientid"`
	Name      string `json:"name" bson:"name"`
	Email     string `json:"email" bson:"email"`
	Precaution []string

}