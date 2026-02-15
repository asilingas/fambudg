import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { DashboardResponse } from "@/lib/types"

export default function DashboardPage() {
  const [data, setData] = useState<DashboardResponse | null>(null)
  const [error, setError] = useState("")

  useEffect(() => {
    const now = new Date()
    api
      .get("/reports/dashboard", {
        params: { month: now.getMonth() + 1, year: now.getFullYear() },
      })
      .then((res) => setData(res.data))
      .catch(() => setError("Failed to load dashboard"))
  }, [])

  if (error) {
    return <p className="text-destructive">{error}</p>
  }

  if (!data) {
    return <p className="text-muted-foreground">Loading...</p>
  }

  const { accounts, monthSummary, recentTransactions } = data

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      {/* Month Summary */}
      <div className="grid gap-4 sm:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Income
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-income">
              {formatCents(monthSummary.totalIncome)}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Expenses
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold text-expense">
              {formatCents(Math.abs(monthSummary.totalExpense))}
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              Net
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p
              className={`text-2xl font-bold ${monthSummary.net >= 0 ? "text-income" : "text-expense"}`}
            >
              {formatCents(monthSummary.net)}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Accounts */}
      <div>
        <h2 className="mb-3 text-lg font-semibold">Accounts</h2>
        {accounts.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            No accounts yet. Create one to get started.
          </p>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {accounts.map((account) => (
              <Card key={account.id}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium">
                      {account.name}
                    </CardTitle>
                    <Badge variant="secondary">{account.type}</Badge>
                  </div>
                </CardHeader>
                <CardContent>
                  <p className="text-xl font-bold">
                    {formatCents(account.balance)}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {account.currency}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      {/* Recent Transactions */}
      <div>
        <h2 className="mb-3 text-lg font-semibold">Recent Transactions</h2>
        {recentTransactions.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            No transactions yet.
          </p>
        ) : (
          <div className="space-y-2">
            {recentTransactions.map((tx) => (
              <div
                key={tx.id}
                className="flex items-center justify-between rounded-lg border p-3"
              >
                <div>
                  <p className="text-sm font-medium">
                    {tx.description || "No description"}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {new Date(tx.date).toLocaleDateString()}
                  </p>
                </div>
                <p
                  className={`text-sm font-semibold ${tx.amount >= 0 ? "text-income" : "text-expense"}`}
                >
                  {formatCents(tx.amount)}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
