package model

import "time"

type DashboardResponse struct {
	Accounts           []*Account           `json:"accounts"`
	MonthSummary       *MonthSummary        `json:"monthSummary"`
	RecentTransactions []*Transaction       `json:"recentTransactions"`
}

type MonthSummary struct {
	Month        int   `json:"month"`
	Year         int   `json:"year"`
	TotalIncome  int64 `json:"totalIncome"`
	TotalExpense int64 `json:"totalExpense"`
	Net          int64 `json:"net"`
}

type CategorySpending struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	TotalAmount  int64   `json:"totalAmount"`
	Percentage   float64 `json:"percentage"`
}

type MemberSpending struct {
	UserID       string `json:"userId"`
	UserName     string `json:"userName"`
	TotalExpense int64  `json:"totalExpense"`
	TotalIncome  int64  `json:"totalIncome"`
	Net          int64  `json:"net"`
}

type ReportFilters struct {
	Month     int
	Year      int
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
	UserID    string
}

type SearchFilters struct {
	Description string
	MinAmount   *int64
	MaxAmount   *int64
	StartDate   string
	EndDate     string
	CategoryID  string
	AccountID   string
	Tags        []string
	UserID      string
}

type SearchResult struct {
	Transactions []*Transaction `json:"transactions"`
	TotalCount   int            `json:"totalCount"`
}

// Helpers to compute date ranges from month/year
func (f *ReportFilters) DateRange() (time.Time, time.Time) {
	start := time.Date(f.Year, time.Month(f.Month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1) // last day of month
	return start, end
}
