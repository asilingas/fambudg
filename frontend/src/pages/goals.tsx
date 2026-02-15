import { useEffect, useState, useCallback } from "react"
import { PageSkeleton } from "@/components/loading-skeleton"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Progress } from "@/components/ui/progress"
import { Badge } from "@/components/ui/badge"
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
import { Plus, Pencil, CheckCircle2 } from "lucide-react"
import { toast } from "sonner"
import { useAuth } from "@/context/auth-context"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { SavingGoal } from "@/lib/types"

interface GoalFormData {
  name: string
  targetAmount: string
  targetDate: string
  priority: string
}

const emptyForm: GoalFormData = { name: "", targetAmount: "", targetDate: "", priority: "1" }

export default function GoalsPage() {
  const { user } = useAuth()
  const isAdmin = user?.role === "admin"

  const [goals, setGoals] = useState<SavingGoal[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [contributeDialogOpen, setContributeDialogOpen] = useState(false)
  const [editing, setEditing] = useState<SavingGoal | null>(null)
  const [contributingGoal, setContributingGoal] = useState<SavingGoal | null>(null)
  const [form, setForm] = useState<GoalFormData>(emptyForm)
  const [contributeAmount, setContributeAmount] = useState("")
  const [submitting, setSubmitting] = useState(false)

  const fetchGoals = useCallback(() => {
    api.get("/saving-goals").then((res) => {
      setGoals(res.data ?? [])
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    fetchGoals()
  }, [fetchGoals])

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDialogOpen(true)
  }

  function openEdit(goal: SavingGoal) {
    setEditing(goal)
    setForm({
      name: goal.name,
      targetAmount: String(goal.targetAmount / 100),
      targetDate: goal.targetDate?.slice(0, 10) ?? "",
      priority: String(goal.priority),
    })
    setDialogOpen(true)
  }

  function openContribute(goal: SavingGoal) {
    setContributingGoal(goal)
    setContributeAmount("")
    setContributeDialogOpen(true)
  }

  async function handleSubmit() {
    setSubmitting(true)
    try {
      const payload = {
        name: form.name,
        targetAmount: Math.round(parseFloat(form.targetAmount) * 100),
        targetDate: form.targetDate || undefined,
        priority: parseInt(form.priority) || 1,
      }
      if (editing) {
        await api.put(`/saving-goals/${editing.id}`, payload)
        toast.success("Goal updated")
      } else {
        await api.post("/saving-goals", payload)
        toast.success("Goal created")
      }
      setDialogOpen(false)
      fetchGoals()
    } catch {
      toast.error("Failed to save goal")
    } finally {
      setSubmitting(false)
    }
  }

  async function handleContribute() {
    if (!contributingGoal) return
    setSubmitting(true)
    try {
      await api.post(`/saving-goals/${contributingGoal.id}/contribute`, {
        amount: Math.round(parseFloat(contributeAmount) * 100),
      })
      toast.success("Contribution added")
      setContributeDialogOpen(false)
      fetchGoals()
    } catch {
      toast.error("Failed to contribute")
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
        <h1 className="text-2xl font-bold">Saving Goals</h1>
        {isAdmin && (
          <Button onClick={openCreate} size="sm">
            <Plus className="mr-1 h-4 w-4" />
            Add Goal
          </Button>
        )}
      </div>

      {goals.length === 0 ? (
        <p className="text-sm text-muted-foreground">No saving goals yet.</p>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {goals.map((goal) => {
            const pct = goal.targetAmount > 0
              ? Math.min((goal.currentAmount / goal.targetAmount) * 100, 100)
              : 0
            const isCompleted = goal.status === "completed"

            return (
              <Card key={goal.id}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <CardTitle className="text-sm font-medium">{goal.name}</CardTitle>
                      {isCompleted && (
                        <Badge variant="default" className="gap-1">
                          <CheckCircle2 className="h-3 w-3" />
                          Completed
                        </Badge>
                      )}
                      {goal.status === "cancelled" && (
                        <Badge variant="secondary">Cancelled</Badge>
                      )}
                    </div>
                    {isAdmin && !isCompleted && (
                      <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openEdit(goal)}>
                        <Pencil className="h-3 w-3" />
                      </Button>
                    )}
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between text-sm mb-2">
                    <span>{formatCents(goal.currentAmount)}</span>
                    <span>of {formatCents(goal.targetAmount)}</span>
                  </div>
                  <Progress value={pct} />
                  <div className="flex items-center justify-between mt-2">
                    <p className="text-xs text-muted-foreground">
                      {pct.toFixed(0)}% reached
                      {goal.targetDate && ` · Due ${goal.targetDate.slice(0, 10)}`}
                    </p>
                    {isAdmin && !isCompleted && goal.status !== "cancelled" && (
                      <Button variant="outline" size="sm" onClick={() => openContribute(goal)}>
                        Contribute
                      </Button>
                    )}
                  </div>
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
            <DialogTitle>{editing ? "Edit Goal" : "New Saving Goal"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="goal-name">Name</Label>
              <Input
                id="goal-name"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="e.g. Vacation Fund"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="goal-target">Target Amount</Label>
              <Input
                id="goal-target"
                type="number"
                step="0.01"
                value={form.targetAmount}
                onChange={(e) => setForm({ ...form, targetAmount: e.target.value })}
                placeholder="e.g. 5000.00"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="goal-date">Target Date (optional)</Label>
              <Input
                id="goal-date"
                type="date"
                value={form.targetDate}
                onChange={(e) => setForm({ ...form, targetDate: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="goal-priority">Priority</Label>
              <Input
                id="goal-priority"
                type="number"
                value={form.priority}
                onChange={(e) => setForm({ ...form, priority: e.target.value })}
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleSubmit} disabled={submitting || !form.name || !form.targetAmount}>
              {submitting ? "Saving..." : "Save"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Contribute Dialog */}
      <Dialog open={contributeDialogOpen} onOpenChange={setContributeDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Contribute to {contributingGoal?.name}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Current: {contributingGoal ? formatCents(contributingGoal.currentAmount) : ""} / {contributingGoal ? formatCents(contributingGoal.targetAmount) : ""}
            </p>
            <div className="space-y-2">
              <Label htmlFor="contribute-amount">Amount</Label>
              <Input
                id="contribute-amount"
                type="number"
                step="0.01"
                value={contributeAmount}
                onChange={(e) => setContributeAmount(e.target.value)}
                placeholder="e.g. 100.00"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setContributeDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleContribute} disabled={submitting || !contributeAmount}>
              {submitting ? "Contributing..." : "Contribute"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
