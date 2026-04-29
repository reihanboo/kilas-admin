package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        CurrentUser `json:"user"`
}

type CurrentUser struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateIssueRequest struct {
	ReporterName  string `json:"reporter_name" binding:"required"`
	ReporterEmail string `json:"reporter_email" binding:"required,email"`
	TransactionID string `json:"transaction_id"`
	Category      string `json:"category" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Priority      string `json:"priority"`
}

type UpdateIssueRequest struct {
	Status     string `json:"status"`
	Priority   string `json:"priority"`
	AdminNotes string `json:"admin_notes"`
}
