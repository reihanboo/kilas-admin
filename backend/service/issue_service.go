package service

import (
	"strings"

	"github.com/reihanboo/kilas-admin/dto"
	"github.com/reihanboo/kilas-admin/model"
	"gorm.io/gorm"
)

type IssueService struct {
	db *gorm.DB
}

type DashboardSummary struct {
	TotalIssues int64 `json:"total_issues"`
	OpenIssues  int64 `json:"open_issues"`
	InReview    int64 `json:"in_review"`
	Resolved    int64 `json:"resolved"`
	Rejected    int64 `json:"rejected"`
}

func NewIssueService(db *gorm.DB) *IssueService {
	return &IssueService{db: db}
}

func (s *IssueService) CreateIssue(input dto.CreateIssueRequest) (*model.Issue, error) {
	priority := strings.TrimSpace(strings.ToLower(input.Priority))
	if priority == "" {
		priority = string(model.IssuePriorityMedium)
	}

	issue := model.Issue{
		ReporterName:  input.ReporterName,
		ReporterEmail: input.ReporterEmail,
		TransactionID: input.TransactionID,
		Category:      input.Category,
		Title:         input.Title,
		Description:   input.Description,
		Priority:      model.IssuePriority(priority),
		Status:        model.IssueStatusOpen,
	}

	if err := s.db.Create(&issue).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func (s *IssueService) ListIssues(status string) ([]model.Issue, error) {
	var issues []model.Issue
	query := s.db.Order("created_at DESC")
	if status != "" {
		query = query.Where("status = ?", strings.ToLower(status))
	}
	if err := query.Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

func (s *IssueService) GetIssueByID(id uint) (*model.Issue, error) {
	var issue model.Issue
	if err := s.db.First(&issue, id).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func (s *IssueService) UpdateIssue(id uint, input dto.UpdateIssueRequest) (*model.Issue, error) {
	issue, err := s.GetIssueByID(id)
	if err != nil {
		return nil, err
	}

	if input.Status != "" {
		issue.Status = model.IssueStatus(strings.ToLower(input.Status))
	}
	if input.Priority != "" {
		issue.Priority = model.IssuePriority(strings.ToLower(input.Priority))
	}
	if input.AdminNotes != "" {
		issue.AdminNotes = input.AdminNotes
	}

	if err := s.db.Save(issue).Error; err != nil {
		return nil, err
	}

	return issue, nil
}

func (s *IssueService) DashboardSummary() (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	if err := s.db.Model(&model.Issue{}).Count(&summary.TotalIssues).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusOpen).Count(&summary.OpenIssues).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusInReview).Count(&summary.InReview).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusResolved).Count(&summary.Resolved).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusRejected).Count(&summary.Rejected).Error; err != nil {
		return nil, err
	}

	return summary, nil
}
