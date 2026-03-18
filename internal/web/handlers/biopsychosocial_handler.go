package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"arandu/internal/application/services"
	"arandu/internal/domain/patient"
	"arandu/internal/infrastructure/repository/sqlite"

	patientComponents "arandu/web/components/patient"
)

type BiopsychosocialServiceInterface interface {
	GetContext(patientID string) (*services.BiopsychosocialContext, error)
	AddMedication(patientID, name, dosage, frequency, prescriber string, startedAt time.Time) (*patient.Medication, error)
	GetMedications(patientID string) ([]*patient.Medication, error)
	SuspendMedication(medicationID string) (*patient.Medication, error)
	ActivateMedication(medicationID string) (*patient.Medication, error)
	RecordVitals(patientID string, date time.Time, sleepHours *float64, appetiteLevel *int, weight *float64, physicalActivity int, notes string) (*patient.Vitals, error)
	GetLatestVitals(patientID string) (*patient.Vitals, error)
	GetAverageVitals(patientID string, days int) (*sqlite.VitalsAverage, error)
}

type BiopsychosocialHandler struct {
	biopsychosocialService BiopsychosocialServiceInterface
}

func NewBiopsychosocialHandler(biopsychosocialService BiopsychosocialServiceInterface) *BiopsychosocialHandler {
	return &BiopsychosocialHandler{
		biopsychosocialService: biopsychosocialService,
	}
}

func (h *BiopsychosocialHandler) GetContextPanel(w http.ResponseWriter, r *http.Request) {
	patientID := extractPatientIDFromPath(r.URL.Path)
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	context, err := h.biopsychosocialService.GetContext(patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewModel := patientComponents.BiopsychosocialPanelViewModel{
		PatientID:     context.PatientID,
		Medications:   toMedicationViewModels(context.AllMedications),
		LatestVitals:  toVitalsViewModel(context.LatestVitals),
		VitalsAverage: toVitalsAverageViewModel(context.VitalsAverage),
	}

	err = patientComponents.BiopsychosocialPanel(viewModel).Render(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BiopsychosocialHandler) AddMedication(w http.ResponseWriter, r *http.Request) {
	patientID := extractPatientIDFromPath(r.URL.Path)
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	dosage := strings.TrimSpace(r.FormValue("dosage"))
	frequency := strings.TrimSpace(r.FormValue("frequency"))
	prescriber := strings.TrimSpace(r.FormValue("prescriber"))
	startedAtStr := r.FormValue("started_at")

	if name == "" {
		h.renderMedicationList(w, r, patientID, "Nome do medicamento é obrigatório")
		return
	}

	var startedAt time.Time
	if startedAtStr != "" {
		var err error
		startedAt, err = time.Parse("2006-01-02", startedAtStr)
		if err != nil {
			startedAt = time.Now()
		}
	} else {
		startedAt = time.Now()
	}

	_, err := h.biopsychosocialService.AddMedication(patientID, name, dosage, frequency, prescriber, startedAt)
	if err != nil {
		h.renderMedicationList(w, r, patientID, err.Error())
		return
	}

	h.renderMedicationList(w, r, patientID, "")
}

func (h *BiopsychosocialHandler) UpdateMedicationStatus(w http.ResponseWriter, r *http.Request) {
	patientID := extractPatientIDFromPath(r.URL.Path)
	medicationID := extractMedicationIDFromPath(r.URL.Path)
	newStatus := r.URL.Query().Get("status")

	if patientID == "" || medicationID == "" {
		http.Error(w, "Patient ID and Medication ID are required", http.StatusBadRequest)
		return
	}

	var err error
	switch newStatus {
	case "suspend":
		_, err = h.biopsychosocialService.SuspendMedication(medicationID)
	case "activate":
		_, err = h.biopsychosocialService.ActivateMedication(medicationID)
	default:
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.renderMedicationList(w, r, patientID, "")
}

func (h *BiopsychosocialHandler) RecordVitals(w http.ResponseWriter, r *http.Request) {
	patientID := extractPatientIDFromPath(r.URL.Path)
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	var sleepHours *float64
	if sh := r.FormValue("sleep_hours"); sh != "" {
		shVal, err := strconv.ParseFloat(sh, 64)
		if err == nil {
			sleepHours = &shVal
		}
	}

	var appetiteLevel *int
	if al := r.FormValue("appetite_level"); al != "" {
		alVal, err := strconv.Atoi(al)
		if err == nil && alVal >= 1 && alVal <= 10 {
			appetiteLevel = &alVal
		}
	}

	var weight *float64
	if wt := r.FormValue("weight"); wt != "" {
		wtVal, err := strconv.ParseFloat(wt, 64)
		if err == nil {
			weight = &wtVal
		}
	}

	physicalActivity := 0
	if pa := r.FormValue("physical_activity"); pa != "" {
		paVal, err := strconv.Atoi(pa)
		if err == nil {
			physicalActivity = paVal
		}
	}

	notes := strings.TrimSpace(r.FormValue("notes"))

	dateStr := r.FormValue("date")
	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			date = time.Now()
		}
	} else {
		date = time.Now()
	}

	_, err := h.biopsychosocialService.RecordVitals(patientID, date, sleepHours, appetiteLevel, weight, physicalActivity, notes)
	if err != nil {
		h.renderVitalsWidget(w, r, patientID, err.Error())
		return
	}

	h.renderVitalsWidget(w, r, patientID, "")
}

func (h *BiopsychosocialHandler) renderMedicationList(w http.ResponseWriter, r *http.Request, patientID, errorMsg string) {
	ctx := context.Background()
	medications, err := h.biopsychosocialService.GetMedications(patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewModel := patientComponents.MedicationListViewModel{
		PatientID:   patientID,
		Medications: toMedicationViewModels(medications),
		Error:       errorMsg,
	}

	err = patientComponents.MedicationList(viewModel).Render(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BiopsychosocialHandler) renderVitalsWidget(w http.ResponseWriter, r *http.Request, patientID, errorMsg string) {
	ctx := context.Background()

	vitals, _ := h.biopsychosocialService.GetLatestVitals(patientID)
	avg, _ := h.biopsychosocialService.GetAverageVitals(patientID, 30)

	viewModel := patientComponents.VitalsWidgetViewModel{
		PatientID:     patientID,
		LatestVitals:  toVitalsViewModel(vitals),
		VitalsAverage: toVitalsAverageViewModel(avg),
		Error:         errorMsg,
	}

	err := patientComponents.VitalsWidget(viewModel).Render(ctx, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func toMedicationViewModels(medications []*patient.Medication) []patientComponents.MedicationListItemViewModel {
	var result []patientComponents.MedicationListItemViewModel
	for _, m := range medications {
		result = append(result, patientComponents.MedicationListItemViewModel{
			ID:          m.ID,
			Name:        m.Name,
			Dosage:      m.Dosage,
			Frequency:   m.Frequency,
			Prescriber:  m.Prescriber,
			Status:      string(m.Status),
			StatusKey:   string(m.Status),
			StatusLabel: getStatusLabel(m.Status),
			StartedAt:   m.StartedAt.Format("02/01/2006"),
			IsActive:    m.Status == patient.MedicationStatusActive,
			IsSuspended: m.Status == patient.MedicationStatusSuspended,
			IsFinished:  m.Status == patient.MedicationStatusFinished,
		})
	}
	return result
}

func getStatusLabel(status patient.MedicationStatus) string {
	switch status {
	case patient.MedicationStatusActive:
		return "Ativo"
	case patient.MedicationStatusSuspended:
		return "Suspenso"
	case patient.MedicationStatusFinished:
		return "Finalizado"
	default:
		return string(status)
	}
}

func toVitalsViewModel(v *patient.Vitals) *patientComponents.VitalsItemViewModel {
	if v == nil {
		return &patientComponents.VitalsItemViewModel{HasData: false}
	}
	vm := &patientComponents.VitalsItemViewModel{
		ID:      v.ID,
		Date:    v.Date.Format("02/01/2006"),
		HasData: true,
		Notes:   v.Notes,
	}
	if v.SleepHours != nil {
		vm.SleepHours = strconv.FormatFloat(*v.SleepHours, 'f', 1, 64)
	}
	if v.AppetiteLevel != nil {
		vm.AppetiteLevel = strconv.Itoa(*v.AppetiteLevel)
	}
	if v.Weight != nil {
		vm.Weight = strconv.FormatFloat(*v.Weight, 'f', 1, 64)
	}
	vm.PhysicalActivity = strconv.Itoa(v.PhysicalActivity)
	return vm
}

func toVitalsAverageViewModel(avg *sqlite.VitalsAverage) *patientComponents.VitalsAverageItemViewModel {
	if avg == nil || avg.Count == 0 {
		return &patientComponents.VitalsAverageItemViewModel{HasData: false}
	}
	vm := &patientComponents.VitalsAverageItemViewModel{
		RecordCount: avg.Count,
		HasData:     true,
	}
	if avg.AverageSleepHours != nil {
		vm.AvgSleepHours = strconv.FormatFloat(*avg.AverageSleepHours, 'f', 1, 64)
	}
	if avg.AverageAppetiteLevel != nil {
		vm.AvgAppetiteLevel = strconv.FormatFloat(*avg.AverageAppetiteLevel, 'f', 1, 64)
	}
	if avg.AverageWeight != nil {
		vm.AvgWeight = strconv.FormatFloat(*avg.AverageWeight, 'f', 1, 64)
	}
	if avg.AveragePhysicalActivity != nil {
		vm.AvgPhysicalActivity = strconv.FormatFloat(*avg.AveragePhysicalActivity, 'f', 0, 64)
	}
	return vm
}

func extractMedicationIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "medications" && i+1 < len(parts) {
			nextPart := parts[i+1]
			if nextPart != "medications" {
				return nextPart
			}
		}
	}
	return ""
}
