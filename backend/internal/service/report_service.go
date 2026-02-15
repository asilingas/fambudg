package service

import (
	"context"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
)

type ReportService struct {
	reportRepo  *repository.ReportRepository
	accountRepo *repository.AccountRepository
}

func NewReportService(reportRepo *repository.ReportRepository, accountRepo *repository.AccountRepository) *ReportService {
	return &ReportService{
		reportRepo:  reportRepo,
		accountRepo: accountRepo,
	}
}

func (s *ReportService) GetDashboard(ctx context.Context, userID string, month, year int) (*model.DashboardResponse, error) {
	accounts, err := s.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	monthSummary, err := s.reportRepo.GetMonthSummary(ctx, userID, month, year)
	if err != nil {
		return nil, err
	}

	recentTransactions, err := s.reportRepo.GetRecentTransactions(ctx, userID, 10)
	if err != nil {
		return nil, err
	}

	return &model.DashboardResponse{
		Accounts:           accounts,
		MonthSummary:       monthSummary,
		RecentTransactions: recentTransactions,
	}, nil
}

// GetDashboardAll returns dashboard for all users (admin)
func (s *ReportService) GetDashboardAll(ctx context.Context, month, year int) (*model.DashboardResponse, error) {
	accounts, err := s.accountRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	monthSummary, err := s.reportRepo.GetMonthSummaryAll(ctx, month, year)
	if err != nil {
		return nil, err
	}

	recentTransactions, err := s.reportRepo.GetRecentTransactionsAll(ctx, 10)
	if err != nil {
		return nil, err
	}

	return &model.DashboardResponse{
		Accounts:           accounts,
		MonthSummary:       monthSummary,
		RecentTransactions: recentTransactions,
	}, nil
}

func (s *ReportService) GetMonthlySummary(ctx context.Context, userID string, month, year int) (*model.MonthSummary, error) {
	return s.reportRepo.GetMonthSummary(ctx, userID, month, year)
}

// GetMonthlySummaryAll returns monthly summary for all users (admin)
func (s *ReportService) GetMonthlySummaryAll(ctx context.Context, month, year int) (*model.MonthSummary, error) {
	return s.reportRepo.GetMonthSummaryAll(ctx, month, year)
}

func (s *ReportService) GetSpendingByCategory(ctx context.Context, userID string, month, year int) ([]*model.CategorySpending, error) {
	return s.reportRepo.GetSpendingByCategory(ctx, userID, month, year)
}

// GetSpendingByCategoryAll returns spending by category for all users (admin)
func (s *ReportService) GetSpendingByCategoryAll(ctx context.Context, month, year int) ([]*model.CategorySpending, error) {
	return s.reportRepo.GetSpendingByCategoryAll(ctx, month, year)
}

func (s *ReportService) GetSpendingByMember(ctx context.Context, month, year int) ([]*model.MemberSpending, error) {
	return s.reportRepo.GetSpendingByMember(ctx, month, year)
}

func (s *ReportService) SearchTransactions(ctx context.Context, filters *model.SearchFilters) (*model.SearchResult, error) {
	return s.reportRepo.SearchTransactions(ctx, filters)
}

// SearchTransactionsAll searches all transactions without user filter (admin)
func (s *ReportService) SearchTransactionsAll(ctx context.Context, filters *model.SearchFilters) (*model.SearchResult, error) {
	return s.reportRepo.SearchTransactionsAll(ctx, filters)
}

func (s *ReportService) GetTrends(ctx context.Context, userID string, months int) ([]*model.TrendPoint, error) {
	return s.reportRepo.GetTrends(ctx, userID, months)
}

// GetTrendsAll returns trends for all users (admin)
func (s *ReportService) GetTrendsAll(ctx context.Context, months int) ([]*model.TrendPoint, error) {
	return s.reportRepo.GetTrendsAll(ctx, months)
}
