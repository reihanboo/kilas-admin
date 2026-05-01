package model

import "time"

type User struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Email             string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Username          string     `gorm:"size:255;not null" json:"username"`
	Password          string     `gorm:"type:text" json:"-"`
	Provider          string     `gorm:"size:50;default:'local'" json:"provider"`
	AvatarURL         string     `json:"avatar_url"`
	Tokens            int        `gorm:"default:500" json:"tokens"`
	LastLoginDate     *time.Time `json:"last_login_date"`
	LoginStreak       int        `gorm:"default:0" json:"login_streak"`
	SubscriptionUntil *time.Time `json:"subscription_until"`
	Language          string     `gorm:"size:10;default:'id'" json:"language"`
	Role              string     `gorm:"size:30;default:'user';not null" json:"role"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;uniqueIndex;not null" json:"name"`
	Price       int       `gorm:"not null" json:"price"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	Type        string    `gorm:"size:50;not null;default:'currency'" json:"type"`
	IsListed    bool      `gorm:"default:true" json:"is_listed"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Transaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	ProductID  uint      `gorm:"not null;index" json:"product_id"`
	Amount     int       `gorm:"not null" json:"amount"`
	Tokens     int       `gorm:"not null" json:"tokens"`
	Status     string    `gorm:"default:'pending'" json:"status"`
	PaymentURL string    `json:"payment_url"`
	CreatedAt  time.Time `json:"created_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type Deck struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	Tags          string    `json:"tags"`
	CloneCount    int       `gorm:"default:0" json:"clone_count"`
	IsAIGenerated bool      `gorm:"default:false" json:"is_ai_generated"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Cards         []Card    `gorm:"foreignKey:DeckID" json:"cards,omitempty"`
}

type Card struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeckID        uint      `gorm:"not null;index" json:"deck_id"`
	Front         string    `gorm:"type:text;not null" json:"front"`
	Back          string    `gorm:"type:text;not null" json:"back"`
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	Interval      int       `gorm:"default:0" json:"interval"`
	Repetitions   int       `gorm:"default:0" json:"repetitions"`
	EaseFactor    float64   `gorm:"default:2.5" json:"ease_factor"`
	Stability     float64   `gorm:"default:0" json:"stability"`
	Difficulty    float64   `gorm:"default:0" json:"difficulty"`
	DueDate       time.Time `json:"due_date"`
	IsAICreated   bool      `gorm:"default:false" json:"is_ai_created"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AIGenerationHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CardCount int       `gorm:"not null" json:"card_count"`
	CreatedAt time.Time `json:"created_at"`
}

type IssueStatus string

const (
	IssueStatusOpen     IssueStatus = "open"
	IssueStatusInReview IssueStatus = "in_review"
	IssueStatusResolved IssueStatus = "resolved"
	IssueStatusRejected IssueStatus = "rejected"
)

type IssuePriority string

const (
	IssuePriorityLow    IssuePriority = "low"
	IssuePriorityMedium IssuePriority = "medium"
	IssuePriorityHigh   IssuePriority = "high"
)

type Issue struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	ReporterName  string        `gorm:"size:120;not null" json:"reporter_name"`
	ReporterEmail string        `gorm:"size:255;not null" json:"reporter_email"`
	TransactionID string        `gorm:"size:120;index" json:"transaction_id"`
	Category      string        `gorm:"size:120;not null" json:"category"`
	Title         string        `gorm:"size:255;not null" json:"title"`
	Description   string        `gorm:"type:text;not null" json:"description"`
	Status        IssueStatus   `gorm:"size:30;default:'open';index;not null" json:"status"`
	Priority      IssuePriority `gorm:"size:30;default:'medium';index;not null" json:"priority"`
	AdminNotes    string        `gorm:"type:text" json:"admin_notes"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
