export interface Account {
  id: string
  userId: string
  name: string
  type: "checking" | "savings" | "credit" | "cash"
  currency: string
  balance: number
  createdAt: string
}

export interface RecurringRule {
  frequency: "daily" | "weekly" | "monthly" | "yearly"
  day?: number
  dayOfWeek?: number
}

export interface Transaction {
  id: string
  userId: string
  accountId: string
  categoryId: string
  amount: number
  type: "expense" | "income" | "transfer"
  description?: string
  date: string
  isShared: boolean
  isRecurring: boolean
  recurringRule?: RecurringRule
  tags?: string[]
  transferToAccountId?: string
  createdAt: string
  updatedAt: string
}

export interface Category {
  id: string
  parentId?: string
  name: string
  type: "expense" | "income"
  icon?: string
  sortOrder: number
}

export interface Budget {
  id: string
  categoryId: string
  amount: number
  month: number
  year: number
  createdAt: string
}

export interface BudgetSummary {
  categoryId: string
  categoryName: string
  budgetAmount: number
  actualAmount: number
  remaining: number
}

export interface SavingGoal {
  id: string
  name: string
  targetAmount: number
  currentAmount: number
  targetDate?: string
  priority: number
  status: "active" | "completed" | "cancelled"
  createdAt: string
  updatedAt: string
}

export interface BillReminder {
  id: string
  name: string
  amount: number
  dueDay: number
  frequency: "monthly" | "quarterly" | "yearly"
  categoryId?: string
  accountId?: string
  isActive: boolean
  nextDueDate: string
  createdAt: string
  updatedAt: string
}

export interface Allowance {
  id: string
  userId: string
  amount: number
  spent: number
  remaining: number
  periodStart: string
  createdAt: string
  updatedAt: string
}

export interface MonthSummary {
  month: number
  year: number
  totalIncome: number
  totalExpense: number
  net: number
}

export interface DashboardResponse {
  accounts: Account[]
  monthSummary: MonthSummary
  recentTransactions: Transaction[]
}

export interface CategorySpending {
  categoryId: string
  categoryName: string
  totalAmount: number
  percentage: number
}

export interface MemberSpending {
  userId: string
  userName: string
  totalExpense: number
  totalIncome: number
  net: number
}

export interface TrendPoint {
  month: number
  year: number
  totalIncome: number
  totalExpense: number
  net: number
}

export interface SearchResult {
  transactions: Transaction[]
  totalCount: number
}

export interface User {
  id: string
  email: string
  name: string
  role: "admin" | "member" | "child"
  createdAt: string
  updatedAt: string
}
