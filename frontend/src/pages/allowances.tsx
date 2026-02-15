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
import { Plus, Pencil } from "lucide-react"
import { toast } from "sonner"
import { useAuth } from "@/context/auth-context"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { Allowance, User } from "@/lib/types"

interface FormData {
  userId: string
  amount: string
  periodStart: string
}

const emptyForm: FormData = { userId: "", amount: "", periodStart: "" }

export default function AllowancesPage() {
  const { user } = useAuth()
  const isAdmin = user?.role === "admin"

  const [allowances, setAllowances] = useState<Allowance[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Allowance | null>(null)
  const [form, setForm] = useState<FormData>(emptyForm)
  const [submitting, setSubmitting] = useState(false)

  const fetchData = useCallback(() => {
    const requests: Promise<any>[] = [api.get("/allowances")]
    if (isAdmin) requests.push(api.get("/users"))

    Promise.all(requests).then(([allowancesRes, usersRes]) => {
      setAllowances(allowancesRes.data ?? [])
      if (usersRes) setUsers(usersRes.data ?? [])
      setLoading(false)
    })
  }, [isAdmin])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  function userName(userId: string) {
    return users.find((u) => u.id === userId)?.name ?? userId
  }

  const childUsers = users.filter((u) => u.role === "child")

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDialogOpen(true)
  }

  function openEdit(allowance: Allowance) {
    setEditing(allowance)
    setForm({
      userId: allowance.userId,
      amount: String(allowance.amount / 100),
      periodStart: allowance.periodStart.slice(0, 10),
    })
    setDialogOpen(true)
  }

  async function handleSubmit() {
    setSubmitting(true)
    try {
      const payload = {
        userId: form.userId,
        amount: Math.round(parseFloat(form.amount) * 100),
        periodStart: form.periodStart,
      }
      if (editing) {
        await api.put(`/allowances/${editing.id}`, payload)
        toast.success("Allowance updated")
      } else {
        await api.post("/allowances", payload)
        toast.success("Allowance created")
      }
      setDialogOpen(false)
      fetchData()
    } catch {
      toast.error("Failed to save allowance")
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
        <h1 className="text-2xl font-bold">Allowances</h1>
        {isAdmin && (
          <Button onClick={openCreate} size="sm">
            <Plus className="mr-1 h-4 w-4" />
            Set Allowance
          </Button>
        )}
      </div>

      {allowances.length === 0 ? (
        <p className="text-sm text-muted-foreground">No allowances set.</p>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {allowances.map((a) => {
            const pct = a.amount > 0 ? Math.min((a.spent / a.amount) * 100, 100) : 0
            const overspent = a.spent > a.amount

            return (
              <Card key={a.id}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium">
                      {isAdmin ? userName(a.userId) : "My Allowance"}
                    </CardTitle>
                    {isAdmin && (
                      <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openEdit(a)}>
                        <Pencil className="h-3 w-3" />
                      </Button>
                    )}
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between text-sm mb-1">
                    <span>Monthly: {formatCents(a.amount)}</span>
                    <span>Period: {a.periodStart.slice(0, 10)}</span>
                  </div>
                  <div className="flex items-center justify-between text-sm mb-2">
                    <span>Spent: {formatCents(a.spent)}</span>
                    <span className={overspent ? "text-destructive font-medium" : "text-income"}>
                      Remaining: {formatCents(a.remaining)}
                    </span>
                  </div>
                  <Progress value={pct} className={overspent ? "[&>div]:bg-destructive" : ""} />
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
            <DialogTitle>{editing ? "Edit Allowance" : "Set Allowance"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Child</Label>
              <Select value={form.userId} onValueChange={(v) => setForm({ ...form, userId: v })}>
                <SelectTrigger><SelectValue placeholder="Select child" /></SelectTrigger>
                <SelectContent>
                  {childUsers.map((u) => (
                    <SelectItem key={u.id} value={u.id}>{u.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="allowance-amount">Monthly Amount</Label>
              <Input id="allowance-amount" type="number" step="0.01" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} placeholder="e.g. 50.00" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="allowance-start">Period Start</Label>
              <Input id="allowance-start" type="date" value={form.periodStart} onChange={(e) => setForm({ ...form, periodStart: e.target.value })} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleSubmit} disabled={submitting || !form.userId || !form.amount || !form.periodStart}>
              {submitting ? "Saving..." : "Save"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
