package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
)

type StudentDetail struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	Email              string     `json:"email"`
	SystemAccess       bool       `json:"systemAccess"`
	Phone              *string    `json:"phone"`
	Gender             *string    `json:"gender"`
	DOB                *time.Time `json:"dob"`
	Class              *string    `json:"class"`
	Section            *string    `json:"section"`
	Roll               *int       `json:"roll"`
	FatherName         *string    `json:"fatherName"`
	FatherPhone        *string    `json:"fatherPhone"`
	MotherName         *string    `json:"motherName"`
	MotherPhone        *string    `json:"motherPhone"`
	GuardianName       *string    `json:"guardianName"`
	GuardianPhone      *string    `json:"guardianPhone"`
	RelationOfGuardian *string    `json:"relationOfGuardian"`
	CurrentAddress     *string    `json:"currentAddress"`
	PermanentAddress   *string    `json:"permanentAddress"`
	AdmissionDate      *time.Time `json:"admissionDate"`
	ReporterName       *string    `json:"reporterName"`
}

var db *sql.DB

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:5432/school_mgmt?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	http.HandleFunc("/api/v1/students/", handleStudentReport)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	log.Printf("Go PDF service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleStudentReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract student ID from path: /api/v1/students/{id}/report
	path := r.URL.Path
	// Expected: /api/v1/students/{id}/report
	parts := splitPath(path)
	if len(parts) < 5 || parts[4] != "report" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	studentID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := getStudentDetail(studentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	pdf := generatePDF(student)

	filename := fmt.Sprintf("student_%d_report.pdf", studentID)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := pdf.Output(w); err != nil {
		log.Printf("Error writing PDF: %v", err)
		return
	}

	log.Printf("Generated PDF report for student %d (%s)", studentID, student.Name)
}

func splitPath(path string) []string {
	var parts []string
	for _, p := range splitString(path, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	var result []string
	current := ""
	for _, c := range s {
		if string(c) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	result = append(result, current)
	return result
}

func getStudentDetail(id int) (*StudentDetail, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			u.is_active,
			p.phone,
			p.gender,
			p.dob,
			p.class_name,
			p.section_name,
			p.roll,
			p.father_name,
			p.father_phone,
			p.mother_name,
			p.mother_phone,
			p.guardian_name,
			p.guardian_phone,
			p.relation_of_guardian,
			p.current_address,
			p.permanent_address,
			p.admission_dt,
			r.name
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		LEFT JOIN users r ON u.reporter_id = r.id
		WHERE u.id = $1 AND u.role_id = 3`

	row := db.QueryRow(query, id)

	var s StudentDetail
	err := row.Scan(
		&s.ID, &s.Name, &s.Email, &s.SystemAccess,
		&s.Phone, &s.Gender, &s.DOB, &s.Class, &s.Section, &s.Roll,
		&s.FatherName, &s.FatherPhone, &s.MotherName, &s.MotherPhone,
		&s.GuardianName, &s.GuardianPhone, &s.RelationOfGuardian,
		&s.CurrentAddress, &s.PermanentAddress, &s.AdmissionDate,
		&s.ReporterName,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func generatePDF(s *StudentDetail) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, 15)

	// Header
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(33, 37, 41)
	pdf.CellFormat(0, 12, "Student Report", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(108, 117, 125)
	pdf.CellFormat(0, 6, fmt.Sprintf("Generated on %s", time.Now().Format("January 2, 2006")), "", 1, "C", false, 0, "")

	// Divider
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY()+3, 200, pdf.GetY()+3)
	pdf.Ln(8)

	// Student Name (large)
	pdf.SetFont("Helvetica", "B", 16)
	pdf.SetTextColor(33, 37, 41)
	pdf.CellFormat(0, 10, s.Name, "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(108, 117, 125)
	pdf.CellFormat(0, 6, fmt.Sprintf("ID: %d | %s", s.ID, s.Email), "", 1, "L", false, 0, "")
	pdf.Ln(6)

	// Section: Basic Information
	addSectionHeader(pdf, "Basic Information")
	addField(pdf, "Name", s.Name)
	addField(pdf, "Email", s.Email)
	addField(pdf, "Gender", ptrStr(s.Gender))
	addField(pdf, "Phone", ptrStr(s.Phone))
	addField(pdf, "Date of Birth", ptrTime(s.DOB))
	addField(pdf, "System Access", boolStr(s.SystemAccess))

	// Section: Academic Information
	addSectionHeader(pdf, "Academic Information")
	addField(pdf, "Class", ptrStr(s.Class))
	addField(pdf, "Section", ptrStr(s.Section))
	addField(pdf, "Roll Number", ptrInt(s.Roll))
	addField(pdf, "Admission Date", ptrTime(s.AdmissionDate))

	// Section: Parent & Guardian Information
	addSectionHeader(pdf, "Parent & Guardian Information")
	addField(pdf, "Father's Name", ptrStr(s.FatherName))
	addField(pdf, "Father's Phone", ptrStr(s.FatherPhone))
	addField(pdf, "Mother's Name", ptrStr(s.MotherName))
	addField(pdf, "Mother's Phone", ptrStr(s.MotherPhone))
	addField(pdf, "Guardian's Name", ptrStr(s.GuardianName))
	addField(pdf, "Guardian's Phone", ptrStr(s.GuardianPhone))
	addField(pdf, "Relation of Guardian", ptrStr(s.RelationOfGuardian))

	// Section: Address Information
	addSectionHeader(pdf, "Address Information")
	addField(pdf, "Current Address", ptrStr(s.CurrentAddress))
	addField(pdf, "Permanent Address", ptrStr(s.PermanentAddress))

	// Section: Administrative
	addSectionHeader(pdf, "Administrative")
	addField(pdf, "Reports To", ptrStr(s.ReporterName))

	// Footer
	pdf.Ln(10)
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(0, 5, "This report was auto-generated by the School Management System.", "", 1, "C", false, 0, "")

	return pdf
}

func addSectionHeader(pdf *gofpdf.Fpdf, title string) {
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(25, 118, 210) // Blue
	pdf.CellFormat(0, 8, title, "", 1, "L", false, 0, "")
	pdf.SetDrawColor(25, 118, 210)
	pdf.Line(10, pdf.GetY(), 80, pdf.GetY())
	pdf.Ln(3)
}

func addField(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(60, 7, label+":", "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(0, 7, value, "", 1, "L", false, 0, "")
}

func ptrStr(s *string) string {
	if s == nil {
		return "N/A"
	}
	return *s
}

func ptrInt(i *int) string {
	if i == nil {
		return "N/A"
	}
	return strconv.Itoa(*i)
}

func ptrTime(t *time.Time) string {
	if t == nil {
		return "N/A"
	}
	return t.Format("January 2, 2006")
}

func boolStr(b bool) string {
	if b {
		return "Active"
	}
	return "Inactive"
}
