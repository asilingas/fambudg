import { useEffect, useState, useCallback } from "react"
import { PageSkeleton } from "@/components/loading-skeleton"
import { Input } from "@/components/ui/input"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
  Legend,
} from "recharts"
import { useAuth } from "@/context/auth-context"
import { useLanguage } from "@/context/language-context"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { MonthSummary, CategorySpending, TrendPoint, MemberSpending } from "@/lib/types"

const now = new Date()

function centsToEuros(cents: number) {
  return Math.round(Math.abs(cents)) / 100
}

function monthLabel(m: number, y: number) {
  return new Date(y, m - 1).toLocaleString("default", { month: "short", year: "2-digit" })
}

export default function ReportsPage() {
  const { user } = useAuth()
  const isAdmin = user?.role === "admin"
  const { t } = useLanguage()

  const [month, setMonth] = useState(now.getMonth() + 1)
  const [year, setYear] = useState(now.getFullYear())

  const [monthlySummary, setMonthlySummary] = useState<MonthSummary | null>(null)
  const [categorySpending, setCategorySpending] = useState<CategorySpending[]>([])
  const [trends, setTrends] = useState<TrendPoint[]>([])
  const [memberSpending, setMemberSpending] = useState<MemberSpending[]>([])
  const [loading, setLoading] = useState(true)

  const fetchData = useCallback(() => {
    setLoading(true)
    const requests = [
      api.get(`/reports/monthly?month=${month}&year=${year}`),
      api.get(`/reports/by-category?month=${month}&year=${year}`),
      api.get("/reports/trends?months=6"),
    ]
    if (isAdmin) {
      requests.push(api.get(`/reports/by-member?month=${month}&year=${year}`))
    }
    Promise.all(requests).then(([monthlyRes, catRes, trendsRes, memberRes]) => {
      setMonthlySummary(monthlyRes.data)
      setCategorySpending(catRes.data ?? [])
      setTrends(trendsRes.data ?? [])
      if (memberRes) setMemberSpending(memberRes.data ?? [])
      setLoading(false)
    })
  }, [month, year, isAdmin])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const chartData = categorySpending.map((cs) => ({
    name: cs.categoryName,
    amount: centsToEuros(cs.totalAmount),
  }))

  const trendData = trends.map((t) => ({
    name: monthLabel(t.month, t.year),
    income: centsToEuros(t.totalIncome),
    expenses: centsToEuros(t.totalExpense),
    net: t.net / 100,
  }))

  const memberData = memberSpending.map((ms) => ({
    name: ms.userName,
    income: centsToEuros(ms.totalIncome),
    expenses: centsToEuros(ms.totalExpense),
    net: ms.net / 100,
  }))

  if (loading) {
    return <PageSkeleton />
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{t("reports.title")}</h1>
        <div className="flex items-center gap-2">
          <Select value={String(month)} onValueChange={(v) => setMonth(parseInt(v))}>
            <SelectTrigger className="w-[120px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {Array.from({ length: 12 }, (_, i) => (
                <SelectItem key={i + 1} value={String(i + 1)}>
                  {new Date(2026, i).toLocaleString("default", { month: "long" })}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Input
            type="number"
            value={year}
            onChange={(e) => setYear(parseInt(e.target.value) || year)}
            className="w-[80px]"
          />
        </div>
      </div>

      <Tabs defaultValue="monthly">
        <TabsList>
          <TabsTrigger value="monthly">{t("reports.monthly")}</TabsTrigger>
          <TabsTrigger value="categories">{t("reports.byCategory")}</TabsTrigger>
          <TabsTrigger value="trends">{t("reports.trends")}</TabsTrigger>
          {isAdmin && <TabsTrigger value="family">{t("reports.family")}</TabsTrigger>}
        </TabsList>

        {/* Monthly Summary Tab */}
        <TabsContent value="monthly" className="space-y-4 pt-4">
          {monthlySummary ? (
            <div className="grid gap-4 sm:grid-cols-3">
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">{t("reports.income")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-2xl font-bold text-income">{formatCents(monthlySummary.totalIncome)}</p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">{t("reports.expenses")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-2xl font-bold text-expense">{formatCents(Math.abs(monthlySummary.totalExpense))}</p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">{t("reports.net")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className={`text-2xl font-bold ${monthlySummary.net >= 0 ? "text-income" : "text-expense"}`}>
                    {formatCents(monthlySummary.net)}
                  </p>
                </CardContent>
              </Card>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">{t("reports.noDataPeriod")}</p>
          )}
        </TabsContent>

        {/* Category Breakdown Tab */}
        <TabsContent value="categories" className="space-y-4 pt-4">
          {chartData.length > 0 ? (
            <Card>
              <CardHeader>
                <CardTitle className="text-sm font-medium">{t("reports.spendingByCategory")}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="h-[300px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={chartData} layout="vertical" margin={{ left: 80 }}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis type="number" tickFormatter={(v) => `€${v}`} />
                      <YAxis type="category" dataKey="name" width={80} />
                      <Tooltip formatter={(value: number) => [`€${value.toFixed(2)}`, t("reports.spent")]} />
                      <Bar dataKey="amount" fill="hsl(var(--primary))" radius={[0, 4, 4, 0]} />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>
          ) : (
            <p className="text-sm text-muted-foreground">{t("reports.noSpendingData")}</p>
          )}

          {categorySpending.length > 0 && (
            <div className="space-y-2">
              {categorySpending.map((cs) => (
                <div key={cs.categoryId} className="flex items-center justify-between rounded-md border p-3">
                  <div>
                    <p className="text-sm font-medium">{cs.categoryName}</p>
                    <p className="text-xs text-muted-foreground">{cs.percentage.toFixed(1)}%</p>
                  </div>
                  <p className="text-sm font-medium">{formatCents(Math.abs(cs.totalAmount))}</p>
                </div>
              ))}
            </div>
          )}
        </TabsContent>

        {/* Trends Tab */}
        <TabsContent value="trends" className="space-y-4 pt-4">
          {trendData.length > 0 ? (
            <Card>
              <CardHeader>
                <CardTitle className="text-sm font-medium">{t("reports.trendTitle")}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="h-[300px]">
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={trendData}>
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="name" />
                      <YAxis tickFormatter={(v) => `€${v}`} />
                      <Tooltip formatter={(value: number) => `€${value.toFixed(2)}`} />
                      <Legend />
                      <Line type="monotone" dataKey="income" stroke="hsl(var(--income))" strokeWidth={2} name={t("reports.income")} />
                      <Line type="monotone" dataKey="expenses" stroke="hsl(var(--expense))" strokeWidth={2} name={t("reports.expenses")} />
                      <Line type="monotone" dataKey="net" stroke="hsl(var(--primary))" strokeWidth={2} strokeDasharray="5 5" name={t("reports.net")} />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </CardContent>
            </Card>
          ) : (
            <p className="text-sm text-muted-foreground">{t("reports.noTrendData")}</p>
          )}
        </TabsContent>

        {/* Family Spending Comparison Tab (Admin Only) */}
        {isAdmin && (
          <TabsContent value="family" className="space-y-4 pt-4">
            {memberData.length > 0 ? (
              <Card>
                <CardHeader>
                  <CardTitle className="text-sm font-medium">{t("reports.familyComparison")}</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="h-[300px]">
                    <ResponsiveContainer width="100%" height="100%">
                      <BarChart data={memberData}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="name" />
                        <YAxis tickFormatter={(v) => `€${v}`} />
                        <Tooltip formatter={(value: number) => `€${value.toFixed(2)}`} />
                        <Legend />
                        <Bar dataKey="income" fill="hsl(var(--income))" name={t("reports.income")} />
                        <Bar dataKey="expenses" fill="hsl(var(--expense))" name={t("reports.expenses")} />
                      </BarChart>
                    </ResponsiveContainer>
                  </div>
                </CardContent>
              </Card>
            ) : (
              <p className="text-sm text-muted-foreground">{t("reports.noFamilyData")}</p>
            )}
          </TabsContent>
        )}
      </Tabs>
    </div>
  )
}
