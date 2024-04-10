package interfaces

import (
	"backend/models"
	"context"
)

type IHMS interface {
	CreateDoctor(customer *models.DoctorDetails) (string, error)
	GetDoctor(ctx context.Context) ([]models.DoctorDetails, error)
	PatientDetails(customer *models.Patient) (error) 
	DoctorLogin(value *models.Login) (string,string, error)
	DoctorClient(value *models.Login)(*models.DoctorDetails,error)
	DoctorAvailability(value *models.UpdateStatus) (string, error)
	PatientDetailsUpdated(value*models.PatientDetails)(*models.PatientData,error)
	ListPatient(ctx context.Context) ([]models.PatientData, error)
	AddMedicine(customer *models.MedicineDetails) (string, error)
	GetMedicine(ctx context.Context) ([]models.MedicineDetails, error)
	AddPatinetPrescription(value *models.PrecautionUpdate) (string, error) }
