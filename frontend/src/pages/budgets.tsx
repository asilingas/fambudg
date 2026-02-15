import { useEffect, useState, useCallback } from "react"
import { PageSkeleton } from "@/components/loading-skeleton"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Progress } from "@/components/ui/progress"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Plus, Pencil, Trash2, AlertTriangle } from "lucide-react"
import { toast } from "sonner"
import { useAuth } from "@/context/auth-context"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { Budget, BudgetSummary, Category } from "@/lib/types"

const now = new Date()

interface FormData {
  categoryId: string
  amount: string
  month: string
  year: string
}

const emptyForm: FormData = {
  categoryId: "",
  amount: "",
  month: String(now.getMonth() + 1),
  year: String(now.getFullYear()),
}

export default function BudgetsPage() {
  const { user } = useAuth()
  const isAdmin = user?.role === "admin"

  const [budgets, setBudgets] = useState<Budget[]>([])
  const [summary, setSummary] = useState<BudgetSummary[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Budget | null>(null)
  const [deleting, setDeleting] = useState<Budget | null>(null)
  const [form, setForm] = useState<FormData>(emptyForm)
  const [submitting, setSubmitting] = useState(false)

  const [month, setMonth] = useState(now.getMonth() + 1)
  const [year, setYear] = useState(now.getFullYear())

  const fetchData = useCallback(() => {
    setLoading(true)
    Promise.all([
      api.get(`/budgets?month=${month}&year=${year}`),
      api.get(`/budgets/summary?month=${month}&year=${year}`),
      api.get("/categories"),
    ]).then(([budgetsRes, summaryRes, catsRes]) => {
      setBudgets(budgetsRes.data ?? [])
      setSummary(summaryRes.data ?? [])
      setCategories(catsRes.data ?? [])
      setLoading(false)
    })
  }, [month, year])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  function categoryName(id: string) {
    return categories.find((c) => c.id === id)?.name ?? "Unknown"
  }

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDialogOpen(true)
  }

  function openEdit(b: Budget) {
    setEditing(b)
    setForm({
      categoryId: b.categoryId,
      amount: String(b.amount / 100),
      month: String(b.month),
      year: String(b.year),
    })
    setDialogOpen(true)
  }

  function openDelete(b: Budget) {
    setDeleting(b)
    setDeleteDialogOpen(true)
  }

  async function handleSubmit() {
    setSubmitting(true)
    try {
      const payload = {
        categoryId: form.categoryId,
        amount: Math.round(parseFloat(form.amount) * 100),
        month: parseInt(form.month),
        year: parseInt(form.year),
      }
      if (editing) {
        await api.put(`/budgets/${editing.id}`, payload)
        toast.success("Budget updated")
      } else {
        await api.post("/budgets", payload)
        toast.success("Budget created")
      }
      setDialogOpen(false)
      fetchData()
    } catch {
      toast.error("Failed to save budget")
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    setSubmitting(true)
    try {
      await api.delete(`/budgets/${deleting.id}`)
      toast.success("Budget deleted")
      setDeleteDialogOpen(false)
      setDeleting(null)
      fetchData()
    } catch {
      toast.error("Failed to delete budget")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return <PageSkeleton />
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Budgets</h1>
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
          {isAdmin && (
            <Button onClick={openCreate} size="sm">
              <Plus className="mr-1 h-4 w-4" />
              Add Budget
            </Button>
          )}
        </div>
      </div>

      {summary.length === 0 ? (
        <p className="text-sm text-muted-foreground">No budgets set for this month.</p>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {summary.map((s) => {
            const spent = Math.abs(s.actualAmount)
            const pct = s.budgetAmount > 0 ? Math.min((spent / s.budgetAmount) * 100, 100) : 0
            const overspent = spent > s.budgetAmount
            const budget = budgets.find((b) => b.categoryId === s.categoryId)

            return (
              <Card key={s.categoryId}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium">
                      {s.categoryName}
                    </CardTitle>
                    {isAdmin && budget && (
                      <div className="flex gap-1">
                        <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openEdit(budget)}>
                          <Pencil className="h-3 w-3" />
                        </Button>
                        <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openDelete(budget)}>
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
                    )}
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between text-sm mb-2">
                    <span>{formatCents(spent)} spent</span>
                    <span>of {formatCents(s.budgetAmount)}</span>
                  </div>
                  <Progress value={pct} className={overspent ? "[&>div]:bg-destructive" : ""} />
                  {overspent && (
                    <div className="flex items-center gap-1 mt-2 text-xs text-destructive">
                      <AlertTriangle className="h-3 w-3" />
                      Overspent by {formatCents(spent - s.budgetAmount)}
                    </div>
                  )}
                  {!overspent && (
                    <p className="text-xs text-muted-foreground mt-2">
                      {formatCents(s.remaining)} remaining
                    </p>
                  )}
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}

      {/* Create/Edit Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? "Edit Budget" : "New Budget"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Category</Label>
              <Select value={form.categoryId} onValueChange={(v) => setForm({ ...form, categoryId: v })}>
                <SelectTrigger>
                  <SelectValue placeholder="Select category" />
                </SelectTrigger>
                <SelectContent>
                  {categories
                    .filter((c) => c.type === "expense")
                    .map((c) => (
                      <SelectItem key={c.id} value={c.id}>
                        {c.name}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="budget-amount">Amount</Label>
              <Input
                id="budget-amount"
                type="number"
                step="0.01"
                value={form.amount}
                onChange={(e) => setForm({ ...form, amount: e.target.value })}
                placeholder="e.g. 500.00"
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Month</Label>
                <Select value={form.month} onValueChange={(v) => setForm({ ...form, month: v })}>
                  <SelectTrigger>
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
              </div>
              <div className="space-y-2">
                <Label htmlFor="budget-year">Year</Label>
                <Input
                  id="budget-year"
                  type="number"
                  value={form.year}
                  onChange={(e) => setForm({ ...form, year: e.target.value })}
                />
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleSubmit} disabled={submitting || !form.categoryId || !form.amount}>
              {submitting ? "Saving..." : "Save"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Budget</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            Are you sure you want to delete the budget for &quot;{deleting ? categoryName(deleting.categoryId) : ""}&quot;?
          </p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={submitting}>
              {submitting ? "Deleting..." : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
