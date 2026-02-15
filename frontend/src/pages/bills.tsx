import { useEffect, useState, useCallback } from "react"
import { PageSkeleton } from "@/components/loading-skeleton"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Plus, Pencil, Trash2, CreditCard } from "lucide-react"
import { toast } from "sonner"
import { useAuth } from "@/context/auth-context"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { BillReminder, Category, Account } from "@/lib/types"

interface FormData {
  name: string
  amount: string
  dueDay: string
  frequency: string
  categoryId: string
  accountId: string
  nextDueDate: string
}

const emptyForm: FormData = {
  name: "",
  amount: "",
  dueDay: "1",
  frequency: "monthly",
  categoryId: "",
  accountId: "",
  nextDueDate: "",
}

export default function BillsPage() {
  const { user } = useAuth()
  const isAdmin = user?.role === "admin"

  const [bills, setBills] = useState<BillReminder[]>([])
  const [upcoming, setUpcoming] = useState<BillReminder[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)

  const [dialogOpen, setDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [payDialogOpen, setPayDialogOpen] = useState(false)
  const [editing, setEditing] = useState<BillReminder | null>(null)
  const [deleting, setDeleting] = useState<BillReminder | null>(null)
  const [paying, setPaying] = useState<BillReminder | null>(null)
  const [form, setForm] = useState<FormData>(emptyForm)
  const [payAccountId, setPayAccountId] = useState("")
  const [payDate, setPayDate] = useState("")
  const [submitting, setSubmitting] = useState(false)

  const fetchData = useCallback(() => {
    Promise.all([
      api.get("/bill-reminders"),
      api.get("/bill-reminders/upcoming"),
      api.get("/categories"),
      api.get("/accounts"),
    ]).then(([billsRes, upcomingRes, catsRes, acctsRes]) => {
      setBills(billsRes.data ?? [])
      setUpcoming(upcomingRes.data ?? [])
      setCategories(catsRes.data ?? [])
      setAccounts(acctsRes.data ?? [])
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  function categoryName(id?: string) {
    if (!id) return "—"
    return categories.find((c) => c.id === id)?.name ?? "—"
  }

  function accountName(id?: string) {
    if (!id) return "—"
    return accounts.find((a) => a.id === id)?.name ?? "—"
  }

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDialogOpen(true)
  }

  function openEdit(bill: BillReminder) {
    setEditing(bill)
    setForm({
      name: bill.name,
      amount: String(bill.amount / 100),
      dueDay: String(bill.dueDay),
      frequency: bill.frequency,
      categoryId: bill.categoryId ?? "",
      accountId: bill.accountId ?? "",
      nextDueDate: bill.nextDueDate.slice(0, 10),
    })
    setDialogOpen(true)
  }

  function openDelete(bill: BillReminder) {
    setDeleting(bill)
    setDeleteDialogOpen(true)
  }

  function openPay(bill: BillReminder) {
    setPaying(bill)
    setPayAccountId(bill.accountId ?? "")
    setPayDate(new Date().toISOString().slice(0, 10))
    setPayDialogOpen(true)
  }

  function isOverdue(dateStr: string) {
    return new Date(dateStr) < new Date()
  }

  async function handleSubmit() {
    setSubmitting(true)
    try {
      const payload = {
        name: form.name,
        amount: Math.round(parseFloat(form.amount) * 100),
        dueDay: parseInt(form.dueDay),
        frequency: form.frequency,
        categoryId: form.categoryId || undefined,
        accountId: form.accountId || undefined,
        nextDueDate: form.nextDueDate,
      }
      if (editing) {
        await api.put(`/bill-reminders/${editing.id}`, payload)
        toast.success("Bill reminder updated")
      } else {
        await api.post("/bill-reminders", payload)
        toast.success("Bill reminder created")
      }
      setDialogOpen(false)
      fetchData()
    } catch {
      toast.error("Failed to save bill reminder")
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    setSubmitting(true)
    try {
      await api.delete(`/bill-reminders/${deleting.id}`)
      toast.success("Bill reminder deleted")
      setDeleteDialogOpen(false)
      setDeleting(null)
      fetchData()
    } catch {
      toast.error("Failed to delete bill reminder")
    } finally {
      setSubmitting(false)
    }
  }

  async function handlePay() {
    if (!paying) return
    setSubmitting(true)
    try {
      await api.post(`/bill-reminders/${paying.id}/pay`, {
        accountId: payAccountId,
        date: payDate,
      })
      toast.success("Bill marked as paid")
      setPayDialogOpen(false)
      setPaying(null)
      fetchData()
    } catch {
      toast.error("Failed to pay bill")
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
        <h1 className="text-2xl font-bold">Bill Reminders</h1>
        {isAdmin && (
          <Button onClick={openCreate} size="sm">
            <Plus className="mr-1 h-4 w-4" />
            Add Bill
          </Button>
        )}
      </div>

      {/* Upcoming Bills */}
      {upcoming.length > 0 && (
        <div className="space-y-3">
          <h2 className="text-lg font-semibold">Upcoming</h2>
          <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
            {upcoming.map((bill) => {
              const overdue = isOverdue(bill.nextDueDate)
              return (
                <Card key={bill.id} className={overdue ? "border-destructive" : ""}>
                  <CardHeader className="pb-2">
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-sm font-medium">{bill.name}</CardTitle>
                      {overdue && <Badge variant="destructive">Overdue</Badge>}
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p className="text-lg font-bold">{formatCents(bill.amount)}</p>
                    <p className="text-xs text-muted-foreground">
                      Due {bill.nextDueDate.slice(0, 10)} · {bill.frequency}
                    </p>
                    <Button
                      variant="outline"
                      size="sm"
                      className="mt-3 w-full"
                      onClick={() => openPay(bill)}
                    >
                      <CreditCard className="mr-1 h-3 w-3" />
                      Mark as Paid
                    </Button>
                  </CardContent>
                </Card>
              )
            })}
          </div>
        </div>
      )}

      {/* All Bills Table */}
      <div className="space-y-3">
        <h2 className="text-lg font-semibold">All Bills</h2>
        {bills.length === 0 ? (
          <p className="text-sm text-muted-foreground">No bill reminders yet.</p>
        ) : (
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Amount</TableHead>
                  <TableHead>Frequency</TableHead>
                  <TableHead>Next Due</TableHead>
                  <TableHead>Category</TableHead>
                  <TableHead>Account</TableHead>
                  <TableHead />
                </TableRow>
              </TableHeader>
              <TableBody>
                {bills.map((bill) => (
                  <TableRow key={bill.id}>
                    <TableCell className="font-medium">{bill.name}</TableCell>
                    <TableCell>{formatCents(bill.amount)}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{bill.frequency}</Badge>
                    </TableCell>
                    <TableCell className={isOverdue(bill.nextDueDate) ? "text-destructive font-medium" : ""}>
                      {bill.nextDueDate.slice(0, 10)}
                    </TableCell>
                    <TableCell className="text-sm">{categoryName(bill.categoryId)}</TableCell>
                    <TableCell className="text-sm">{accountName(bill.accountId)}</TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openPay(bill)}>
                          <CreditCard className="h-3 w-3" />
                        </Button>
                        {isAdmin && (
                          <>
                            <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openEdit(bill)}>
                              <Pencil className="h-3 w-3" />
                            </Button>
                            <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openDelete(bill)}>
                              <Trash2 className="h-3 w-3" />
                            </Button>
                          </>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </div>

      {/* Create/Edit Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editing ? "Edit Bill Reminder" : "New Bill Reminder"}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="bill-name">Name</Label>
              <Input id="bill-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="e.g. Rent" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="bill-amount">Amount</Label>
                <Input id="bill-amount" type="number" step="0.01" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="bill-dueday">Due Day</Label>
                <Input id="bill-dueday" type="number" min="1" max="31" value={form.dueDay} onChange={(e) => setForm({ ...form, dueDay: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Frequency</Label>
                <Select value={form.frequency} onValueChange={(v) => setForm({ ...form, frequency: v })}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="monthly">Monthly</SelectItem>
                    <SelectItem value="quarterly">Quarterly</SelectItem>
                    <SelectItem value="yearly">Yearly</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="bill-nextdue">Next Due Date</Label>
                <Input id="bill-nextdue" type="date" value={form.nextDueDate} onChange={(e) => setForm({ ...form, nextDueDate: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Category (optional)</Label>
                <Select value={form.categoryId} onValueChange={(v) => setForm({ ...form, categoryId: v === "none" ? "" : v })}>
                  <SelectTrigger><SelectValue placeholder="None" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">None</SelectItem>
                    {categories.filter((c) => c.type === "expense").map((c) => (
                      <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Account (optional)</Label>
                <Select value={form.accountId} onValueChange={(v) => setForm({ ...form, accountId: v === "none" ? "" : v })}>
                  <SelectTrigger><SelectValue placeholder="None" /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">None</SelectItem>
                    {accounts.map((a) => (
                      <SelectItem key={a.id} value={a.id}>{a.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>Cancel</Button>
            <Button onClick={handleSubmit} disabled={submitting || !form.name || !form.amount || !form.nextDueDate}>
              {submitting ? "Saving..." : "Save"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Bill Reminder</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            Are you sure you want to delete &quot;{deleting?.name}&quot;?
          </p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={submitting}>
              {submitting ? "Deleting..." : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Pay Dialog */}
      <Dialog open={payDialogOpen} onOpenChange={setPayDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Pay Bill: {paying?.name}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Amount: {paying ? formatCents(paying.amount) : ""}
            </p>
            <div className="space-y-2">
              <Label>Pay from Account</Label>
              <Select value={payAccountId} onValueChange={setPayAccountId}>
                <SelectTrigger><SelectValue placeholder="Select account" /></SelectTrigger>
                <SelectContent>
                  {accounts.map((a) => (
                    <SelectItem key={a.id} value={a.id}>{a.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="pay-date">Date</Label>
              <Input id="pay-date" type="date" value={payDate} onChange={(e) => setPayDate(e.target.value)} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPayDialogOpen(false)}>Cancel</Button>
            <Button onClick={handlePay} disabled={submitting || !payAccountId || !payDate}>
              {submitting ? "Paying..." : "Pay"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
