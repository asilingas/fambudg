import { BrowserRouter, Routes, Route } from "react-router-dom"
import { AuthProvider } from "@/context/auth-context"
import { ProtectedRoute } from "@/components/protected-route"
import { AppLayout } from "@/components/layout/app-layout"
import { Toaster } from "@/components/ui/sonner"
import LoginPage from "@/pages/login"
import DashboardPage from "@/pages/dashboard"
import TransactionsPage from "@/pages/transactions"
import AccountsPage from "@/pages/accounts"
import CategoriesPage from "@/pages/categories"
import BudgetsPage from "@/pages/budgets"
import ReportsPage from "@/pages/reports"
import SearchPage from "@/pages/search"
import GoalsPage from "@/pages/goals"
import BillsPage from "@/pages/bills"
import TransfersPage from "@/pages/transfers"

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<ProtectedRoute />}>
            <Route element={<AppLayout />}>
              <Route path="/" element={<DashboardPage />} />
              <Route path="/transactions" element={<TransactionsPage />} />
              <Route path="/accounts" element={<AccountsPage />} />
              <Route path="/categories" element={<CategoriesPage />} />
              <Route path="/budgets" element={<ProtectedRoute allowedRoles={["admin", "member"]} />}>
                <Route index element={<BudgetsPage />} />
              </Route>
              <Route path="/goals" element={<ProtectedRoute allowedRoles={["admin", "member"]} />}>
                <Route index element={<GoalsPage />} />
              </Route>
              <Route path="/bills" element={<ProtectedRoute allowedRoles={["admin", "member"]} />}>
                <Route index element={<BillsPage />} />
              </Route>
              <Route path="/transfers" element={<ProtectedRoute allowedRoles={["admin", "member"]} />}>
                <Route index element={<TransfersPage />} />
              </Route>
              <Route path="/reports" element={<ReportsPage />} />
              <Route path="/search" element={<SearchPage />} />
            </Route>
          </Route>
        </Routes>
        <Toaster />
      </AuthProvider>
    </BrowserRouter>
  )
}
