package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/reihanboo/kilas-admin/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminCRUDService struct {
	db *gorm.DB
}

type ListOptions struct {
	Q      string
	Limit  int
	Offset int
}

type DashboardOverview struct {
	Users               int64 `json:"users"`
	Transactions        int64 `json:"transactions"`
	Products            int64 `json:"products"`
	Decks               int64 `json:"decks"`
	Cards               int64 `json:"cards"`
	AIGenerationHistory int64 `json:"ai_generation_history"`
	Issues              int64 `json:"issues"`
	OpenIssues          int64 `json:"open_issues"`
	InReviewIssues      int64 `json:"in_review_issues"`
	ResolvedIssues      int64 `json:"resolved_issues"`
	RejectedIssues      int64 `json:"rejected_issues"`
}

func NewAdminCRUDService(db *gorm.DB) *AdminCRUDService {
	return &AdminCRUDService{db: db}
}

func (s *AdminCRUDService) Summary() (*DashboardOverview, error) {
	o := &DashboardOverview{}
	if err := s.db.Model(&model.User{}).Count(&o.Users).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Transaction{}).Count(&o.Transactions).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Product{}).Count(&o.Products).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Deck{}).Count(&o.Decks).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Card{}).Count(&o.Cards).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.AIGenerationHistory{}).Count(&o.AIGenerationHistory).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Count(&o.Issues).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusOpen).Count(&o.OpenIssues).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusInReview).Count(&o.InReviewIssues).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusResolved).Count(&o.ResolvedIssues).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Issue{}).Where("status = ?", model.IssueStatusRejected).Count(&o.RejectedIssues).Error; err != nil {
		return nil, err
	}
	return o, nil
}

func (s *AdminCRUDService) List(entity string, options ListOptions) (interface{}, error) {
	cfg, err := getEntityConfig(entity)
	if err != nil {
		return nil, err
	}

	list := cfg.newSlice()
	query := s.db.Order("id DESC")
	for _, relation := range cfg.preloads {
		query = query.Preload(relation)
	}

	if options.Q != "" && len(cfg.searchableColumns) > 0 {
		like := "%" + strings.ToLower(strings.TrimSpace(options.Q)) + "%"
		for i, col := range cfg.searchableColumns {
			clause := "LOWER(" + col + ") LIKE ?"
			if i == 0 {
				query = query.Where(clause, like)
			} else {
				query = query.Or(clause, like)
			}
		}
	}

	limit := options.Limit
	if limit <= 0 {
		limit = 0
	}
	if limit > 200 {
		limit = 200
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	if err := query.Find(list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (s *AdminCRUDService) Get(entity string, id uint) (interface{}, error) {
	cfg, err := getEntityConfig(entity)
	if err != nil {
		return nil, err
	}

	obj := cfg.newModel()
	query := s.db
	for _, relation := range cfg.preloads {
		query = query.Preload(relation)
	}
	if err := query.First(obj, id).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *AdminCRUDService) Create(entity string, payload map[string]interface{}) (interface{}, error) {
	cfg, err := getEntityConfig(entity)
	if err != nil {
		return nil, err
	}

	cleanPayload, err := sanitizePayload(payload)
	if err != nil {
		return nil, err
	}
	if err := cfg.beforeWrite(cleanPayload); err != nil {
		return nil, err
	}
	if err := cfg.validateRelations(s.db, cleanPayload); err != nil {
		return nil, err
	}

	obj := cfg.newModel()
	if err := s.db.Model(obj).Create(cleanPayload).Error; err != nil {
		return nil, err
	}

	return cleanPayload, nil
}

func (s *AdminCRUDService) Update(entity string, id uint, payload map[string]interface{}) (interface{}, error) {
	cfg, err := getEntityConfig(entity)
	if err != nil {
		return nil, err
	}

	obj := cfg.newModel()
	if err := s.db.First(obj, id).Error; err != nil {
		return nil, err
	}

	cleanPayload, err := sanitizePayload(payload)
	if err != nil {
		return nil, err
	}
	if err := cfg.beforeWrite(cleanPayload); err != nil {
		return nil, err
	}
	if err := cfg.validateRelations(s.db, cleanPayload); err != nil {
		return nil, err
	}

	if len(cleanPayload) == 0 {
		return s.Get(entity, id)
	}
	if err := s.db.Model(obj).Updates(cleanPayload).Error; err != nil {
		return nil, err
	}
	return s.Get(entity, id)
}

func (s *AdminCRUDService) Delete(entity string, id uint) error {
	cfg, err := getEntityConfig(entity)
	if err != nil {
		return err
	}
	obj := cfg.newModel()
	return s.db.Delete(obj, id).Error
}

type entityConfig struct {
	newModel          func() interface{}
	newSlice          func() interface{}
	preloads          []string
	searchableColumns []string
	beforeWrite       func(payload map[string]interface{}) error
	validateRelations func(db *gorm.DB, payload map[string]interface{}) error
}

func getEntityConfig(entity string) (*entityConfig, error) {
	switch strings.ToLower(entity) {
	case "users":
		return &entityConfig{
			newModel:          func() interface{} { return &model.User{} },
			newSlice:          func() interface{} { return &[]model.User{} },
			searchableColumns: []string{"email", "username", "language", "role"},
			beforeWrite:       compose(transformUserPassword, normalizeDateFields("last_login_date", "subscription_until", "created_at", "updated_at")),
			validateRelations: noRelationValidation,
		}, nil
	case "transactions":
		return &entityConfig{
			newModel:          func() interface{} { return &model.Transaction{} },
			newSlice:          func() interface{} { return &[]model.Transaction{} },
			preloads:          []string{"User", "Product"},
			searchableColumns: []string{"status", "payment_url"},
			beforeWrite:       normalizeDateFields("created_at"),
			validateRelations: composeRelationValidators(requireExistingID("user_id", &model.User{}), requireExistingID("product_id", &model.Product{})),
		}, nil
	case "products":
		return &entityConfig{
			newModel:          func() interface{} { return &model.Product{} },
			newSlice:          func() interface{} { return &[]model.Product{} },
			searchableColumns: []string{"name", "type", "description"},
			beforeWrite:       normalizeDateFields("created_at", "updated_at"),
			validateRelations: noRelationValidation,
		}, nil
	case "decks":
		return &entityConfig{
			newModel:          func() interface{} { return &model.Deck{} },
			newSlice:          func() interface{} { return &[]model.Deck{} },
			preloads:          []string{"Cards"},
			searchableColumns: []string{"title", "description", "tags"},
			beforeWrite:       normalizeDateFields("created_at", "updated_at"),
			validateRelations: requireExistingID("user_id", &model.User{}),
		}, nil
	case "cards":
		return &entityConfig{
			newModel:          func() interface{} { return &model.Card{} },
			newSlice:          func() interface{} { return &[]model.Card{} },
			searchableColumns: []string{"front", "back", "front_image_url", "back_image_url"},
			beforeWrite:       normalizeDateFields("due_date", "created_at", "updated_at"),
			validateRelations: requireExistingID("deck_id", &model.Deck{}),
		}, nil
	case "ai_generation_history":
		return &entityConfig{
			newModel:          func() interface{} { return &model.AIGenerationHistory{} },
			newSlice:          func() interface{} { return &[]model.AIGenerationHistory{} },
			searchableColumns: []string{"text"},
			beforeWrite:       normalizeDateFields("created_at"),
			validateRelations: requireExistingID("user_id", &model.User{}),
		}, nil
	case "issues":
		return &entityConfig{
			newModel:          func() interface{} { return &model.Issue{} },
			newSlice:          func() interface{} { return &[]model.Issue{} },
			searchableColumns: []string{"reporter_name", "reporter_email", "transaction_id", "category", "title", "description", "status", "priority", "admin_notes"},
			beforeWrite:       normalizeDateFields("created_at", "updated_at"),
			validateRelations: noRelationValidation,
		}, nil
	default:
		return nil, errors.New("unsupported entity")
	}
}

func sanitizePayload(payload map[string]interface{}) (map[string]interface{}, error) {
	clean := map[string]interface{}{}
	for k, v := range payload {
		key := strings.ToLower(strings.TrimSpace(k))
		if key == "" {
			continue
		}
		switch key {
		case "id":
			continue
		default:
			clean[key] = v
		}
	}
	return clean, nil
}

func passthrough(_ map[string]interface{}) error {
	return nil
}

func compose(fns ...func(map[string]interface{}) error) func(map[string]interface{}) error {
	return func(payload map[string]interface{}) error {
		for _, fn := range fns {
			if err := fn(payload); err != nil {
				return err
			}
		}
		return nil
	}
}

func transformUserPassword(payload map[string]interface{}) error {
	rawPassword, hasPassword := payload["password"]
	if !hasPassword {
		return nil
	}
	password, ok := rawPassword.(string)
	if !ok || strings.TrimSpace(password) == "" {
		return errors.New("password must be a non-empty string")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	payload["password"] = string(hashed)
	return nil
}

func normalizeDateFields(fields ...string) func(map[string]interface{}) error {
	return func(payload map[string]interface{}) error {
		for _, field := range fields {
			raw, ok := payload[field]
			if !ok || raw == nil {
				continue
			}

			switch v := raw.(type) {
			case string:
				if strings.TrimSpace(v) == "" {
					delete(payload, field)
					continue
				}
				parsed, err := parseFlexibleTime(v)
				if err != nil {
					return fmt.Errorf("invalid %s format", field)
				}
				payload[field] = parsed.Format("2006-01-02 15:04:05")
			case time.Time:
				payload[field] = v.Format("2006-01-02 15:04:05")
			default:
				return fmt.Errorf("invalid %s format", field)
			}
		}
		return nil
	}
}

func parseFlexibleTime(input string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, input); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("invalid time")
}

func noRelationValidation(_ *gorm.DB, _ map[string]interface{}) error {
	return nil
}

func composeRelationValidators(validators ...func(*gorm.DB, map[string]interface{}) error) func(*gorm.DB, map[string]interface{}) error {
	return func(db *gorm.DB, payload map[string]interface{}) error {
		for _, validator := range validators {
			if err := validator(db, payload); err != nil {
				return err
			}
		}
		return nil
	}
}

func requireExistingID(field string, modelRef interface{}) func(*gorm.DB, map[string]interface{}) error {
	return func(db *gorm.DB, payload map[string]interface{}) error {
		raw, ok := payload[field]
		if !ok {
			return nil
		}
		id, err := coerceUint(raw)
		if err != nil || id == 0 {
			return fmt.Errorf("%s must be a valid id", field)
		}

		var count int64
		if err := db.Model(modelRef).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("referenced %s does not exist", field)
		}
		payload[field] = id
		return nil
	}
}

func coerceUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case float64:
		return uint(v), nil
	case float32:
		return uint(v), nil
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case uint:
		return v, nil
	case uint64:
		return uint(v), nil
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return 0, errors.New("empty id")
		}
		var out uint
		_, err := fmt.Sscanf(trimmed, "%d", &out)
		if err != nil {
			return 0, err
		}
		return out, nil
	default:
		return 0, errors.New("invalid id type")
	}
}
