package model

type TrendPoint struct {
	Month        int   `json:"month"`
	Year         int   `json:"year"`
	TotalIncome  int64 `json:"totalIncome"`
	TotalExpense int64 `json:"totalExpense"`
	Net          int64 `json:"net"`
}
