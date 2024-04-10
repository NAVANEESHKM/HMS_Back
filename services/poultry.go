package services

import (
	"backend/interfaces"
	"backend/models"
	"context"
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HMSService struct {
	AdminCollection    *mongo.Collection
	DoctorCollection   *mongo.Collection
	PatientCollection  *mongo.Collection
	MedicineCollection *mongo.Collection
	ctx                context.Context
}

func PoultryServiceInit(admin *mongo.Collection, doctor *mongo.Collection, patient *mongo.Collection, Medicine *mongo.Collection, ctx context.Context) interfaces.IHMS {
	return &HMSService{admin, doctor, patient, Medicine, ctx}
}

func (p *HMSService) CreateDoctor(customer *models.DoctorDetails) (string, error) {
	// Execute the update operation
	_, err := p.AdminCollection.InsertOne(context.Background(), customer)
	if err != nil {
		fmt.Println("Transaction not updated")
		return "Transaction not updated", err
	}

	return "greater", nil
}

func (p *HMSService) GetDoctor(ctx context.Context) ([]models.DoctorDetails, error) {
	var values []models.DoctorDetails
	find := bson.M{}
	// Execute the update operation
	cursor, err := p.AdminCollection.Find(context.Background(), find)
	if err != nil {
		fmt.Println("Transaction not updated")
		return nil, err
	}
	for cursor.Next(ctx) {
		var admin models.DoctorDetails
		if err := cursor.Decode(&admin); err != nil {
			fmt.Println("Error decoding document:", err)
			return nil, err
		}
		values = append(values, admin)
	}

	return values, nil
}

func (p *HMSService) PatientDetails(customer *models.Patient) error {
	var Doctorname models.DoctorName
	lenstr := strconv.Itoa(int(customer.Number))

	Doctorname.DoctorName = customer.Doctor
	Doctorname.Purpose = customer.Purpose
	Doctorname.Precaution = []string{}
	Doctorname.Date = customer.Date

	log.Println("Patient id .........", customer.PatientID)
	// Execute the update operation
	if customer.PatientID == "NEW" {
		customer.PatientID = string(customer.Name[0]) + lenstr[0:1] + string(customer.Name[len(customer.Name)-1]) + lenstr[len(lenstr)-1:]

		log.Println("inside new")
		var data *models.PatientData
		if data == nil {
			data = &models.PatientData{} // Replace YourStructType with the actual type of data

		}
		data.DoctorNames = []models.DoctorName{}
		data.Email = customer.Email
		data.Name = customer.Name
		data.Number = customer.Number
		data.PatientID = customer.PatientID
		data.DoctorNames = append(data.DoctorNames, Doctorname)
		_, err := p.PatientCollection.InsertOne(context.Background(), data)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("PATIENT ID ALREADY EXIST")
		find := bson.M{"patientid": customer.PatientID}
		update := bson.M{
			"$push": bson.M{
				"doctornames": Doctorname,
			},
		}
		_, err := p.PatientCollection.UpdateOne(context.Background(), find, update)
		if err != nil {
			fmt.Println("Inside")
			return err
		}
	}
	customer.Precaution = []string{}
	find := bson.M{"name": customer.Doctor}
	update := bson.M{
		"$push": bson.M{
			"patients": customer,
		},
	}
	_, err := p.AdminCollection.UpdateOne(context.Background(), find, update)
	if err != nil {
		return err
	}
	return nil
}

func (p *HMSService) DoctorLogin(value *models.Login) (string, string, error) {
	find := bson.M{"emailid": value.Email}
	err := p.AdminCollection.FindOne(context.Background(), find)

	if err.Err() != nil {
		if err.Err() == mongo.ErrNoDocuments {
			return "", "Not Found", nil
		} else {
			return "", "Error", err.Err()
		}

	}
	var doctor models.DoctorDetails
	if err1 := err.Decode(&doctor); err1 != nil {
		return "", "Error", err1

	}

	return doctor.PhoneNumber, "Found", nil

}

func (p *HMSService) DoctorClient(value *models.Login) (*models.DoctorDetails, error) {
	var store *models.DoctorDetails
	find := bson.M{"emailid": value.Email}
	cursor := p.AdminCollection.FindOne(context.Background(), find)
	if err := cursor.Decode(&store); err != nil {
		fmt.Println("Error decoding document:", err)
		return nil, err
	}
	return store, nil
}

func (p *HMSService) DoctorAvailability(value *models.UpdateStatus) (string, error) {
	find := bson.M{"emailid": value.Email}
	update := bson.M{
		"$set": bson.M{
			"availability": value.Availability,
		},
	}
	_, err := p.AdminCollection.UpdateOne(context.Background(), find, update)
	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (p *HMSService) PatientDetailsUpdated(value *models.PatientDetails) (*models.PatientData, error) {
	find := bson.M{"patientid": value.PatientID}
	var store *models.PatientData
	cursor := p.PatientCollection.FindOne(context.Background(), find)
	if err := cursor.Decode(&store); err != nil {
		fmt.Println("Error decoding document:", err)
		return nil, err
	}
	return store, nil

}

func (p *HMSService) ListPatient(ctx context.Context) ([]models.PatientData, error) {
	var values []models.PatientData
	find := bson.M{}
	// Execute the update operation
	cursor, err := p.PatientCollection.Find(context.Background(), find)
	if err != nil {
		fmt.Println("Transaction not updated")
		return nil, err
	}
	for cursor.Next(ctx) {
		var admin models.PatientData
		if err := cursor.Decode(&admin); err != nil {
			fmt.Println("Error decoding document:", err)
			return nil, err
		}
		values = append(values, admin)
	}

	return values, nil
}

func (p *HMSService) AddMedicine(customer *models.MedicineDetails) (string, error) {
	// Execute the update operation
	_, err := p.MedicineCollection.InsertOne(context.Background(), customer)
	if err != nil {
		fmt.Println("Transaction not updated")
		return "Transaction not updated", err
	}

	return "greater", nil
}

func (p *HMSService) GetMedicine(ctx context.Context) ([]models.MedicineDetails, error) {
	var values []models.MedicineDetails
	find := bson.M{}
	// Execute the update operation
	cursor, err := p.MedicineCollection.Find(context.Background(), find)
	if err != nil {
		fmt.Println("Transaction not updated")
		return nil, err
	}
	for cursor.Next(ctx) {
		var admin models.MedicineDetails
		if err := cursor.Decode(&admin); err != nil {
			fmt.Println("Error decoding document:", err)
			return nil, err
		}
		values = append(values, admin)
	}

	return values, nil
}

func (p *HMSService) AddPatinetPrescription(value *models.PrecautionUpdate) (string, error) {

	find := bson.M{
		"emailid":            value.DoctorEmail,
		"patients.patientid": value.PatientID,
	}

	// Define the update to push the patient data
	update := bson.M{
		"$push": bson.M{
			"patients.$.precaution": bson.M{"$each": value.Precaution},
		},
	}
	_, err := p.AdminCollection.UpdateOne(context.Background(), find, update)
	if err != nil {
		return "Failed", err
	}

	find1 := bson.M{
		"patientid":              value.PatientID,
		"doctornames.doctorname": value.DoctorName,
	}

	// Define the update to push the patient data
	update1 := bson.M{
		"$push": bson.M{
			"doctornames.$.precaution": bson.M{"$each": value.Precaution},
		},
	}
	_, err1 := p.PatientCollection.UpdateOne(context.Background(), find1, update1)
	if err != nil {
		return "Failed", err1
	}

	return "Success", nil
}
