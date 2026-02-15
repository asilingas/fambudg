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
import { useLanguage } from "@/context/language-context"
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
  const { t } = useLanguage()

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
        toast.success(t("bills.updated"))
      } else {
        await api.post("/bill-reminders", payload)
        toast.success(t("bills.created"))
      }
      setDialogOpen(false)
      fetchData()
    } catch {
      toast.error(t("bills.saveFailed"))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    setSubmitting(true)
    try {
      await api.delete(`/bill-reminders/${deleting.id}`)
      toast.success(t("bills.deleted"))
      setDeleteDialogOpen(false)
      setDeleting(null)
      fetchData()
    } catch {
      toast.error(t("bills.deleteFailed"))
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
      toast.success(t("bills.paid"))
      setPayDialogOpen(false)
      setPaying(null)
      fetchData()
    } catch {
      toast.error(t("bills.payFailed"))
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
        <h1 className="text-2xl font-bold">{t("bills.title")}</h1>
        {isAdmin && (
          <Button onClick={openCreate} size="sm">
            <Plus className="mr-1 h-4 w-4" />
            {t("bills.add")}
          </Button>
        )}
      </div>

      {/* Upcoming Bills */}
      {upcoming.length > 0 && (
        <div className="space-y-3">
          <h2 className="text-lg font-semibold">{t("bills.upcoming")}</h2>
          <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
            {upcoming.map((bill) => {
              const overdue = isOverdue(bill.nextDueDate)
              return (
                <Card key={bill.id} className={overdue ? "border-destructive" : ""}>
                  <CardHeader className="pb-2">
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-sm font-medium">{bill.name}</CardTitle>
                      {overdue && <Badge variant="destructive">{t("bills.overdue")}</Badge>}
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p className="text-lg font-bold">{formatCents(bill.amount)}</p>
                    <p className="text-xs text-muted-foreground">
                      {t("bills.due")} {bill.nextDueDate.slice(0, 10)} · {bill.frequency}
                    </p>
                    <Button
                      variant="outline"
                      size="sm"
                      className="mt-3 w-full"
                      onClick={() => openPay(bill)}
                    >
                      <CreditCard className="mr-1 h-3 w-3" />
                      {t("bills.markAsPaid")}
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
        <h2 className="text-lg font-semibold">{t("bills.allBills")}</h2>
        {bills.length === 0 ? (
          <p className="text-sm text-muted-foreground">{t("bills.noData")}</p>
        ) : (
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t("bills.name")}</TableHead>
                  <TableHead>{t("bills.amount")}</TableHead>
                  <TableHead>{t("bills.frequency")}</TableHead>
                  <TableHead>{t("bills.nextDue")}</TableHead>
                  <TableHead>{t("bills.category")}</TableHead>
                  <TableHead>{t("bills.account")}</TableHead>
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
            <DialogTitle>{editing ? t("bills.editTitle") : t("bills.newTitle")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="bill-name">{t("bills.name")}</Label>
              <Input id="bill-name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="e.g. Rent" />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="bill-amount">{t("bills.amount")}</Label>
                <Input id="bill-amount" type="number" step="0.01" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="bill-dueday">{t("bills.dueDay")}</Label>
                <Input id="bill-dueday" type="number" min="1" max="31" value={form.dueDay} onChange={(e) => setForm({ ...form, dueDay: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>{t("bills.frequency")}</Label>
                <Select value={form.frequency} onValueChange={(v) => setForm({ ...form, frequency: v })}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="monthly">{t("bills.monthly")}</SelectItem>
                    <SelectItem value="quarterly">{t("bills.quarterly")}</SelectItem>
                    <SelectItem value="yearly">{t("bills.yearly")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="bill-nextdue">{t("bills.nextDueDate")}</Label>
                <Input id="bill-nextdue" type="date" value={form.nextDueDate} onChange={(e) => setForm({ ...form, nextDueDate: e.target.value })} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>{t("bills.categoryOptional")}</Label>
                <Select value={form.categoryId} onValueChange={(v) => setForm({ ...form, categoryId: v === "none" ? "" : v })}>
                  <SelectTrigger><SelectValue placeholder={t("bills.none")} /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">{t("bills.none")}</SelectItem>
                    {categories.filter((c) => c.type === "expense").map((c) => (
                      <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>{t("bills.accountOptional")}</Label>
                <Select value={form.accountId} onValueChange={(v) => setForm({ ...form, accountId: v === "none" ? "" : v })}>
                  <SelectTrigger><SelectValue placeholder={t("bills.none")} /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="none">{t("bills.none")}</SelectItem>
                    {accounts.map((a) => (
                      <SelectItem key={a.id} value={a.id}>{a.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button onClick={handleSubmit} disabled={submitting || !form.name || !form.amount || !form.nextDueDate}>
              {submitting ? t("common.saving") : t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("bills.deleteTitle")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            {t("bills.deleteConfirm").replace("{name}", deleting?.name ?? "")}
          </p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={submitting}>
              {submitting ? t("common.deleting") : t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Pay Dialog */}
      <Dialog open={payDialogOpen} onOpenChange={setPayDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("bills.payTitle").replace("{name}", paying?.name ?? "")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">
              {t("bills.payAmount").replace("{amount}", paying ? formatCents(paying.amount) : "")}
            </p>
            <div className="space-y-2">
              <Label>{t("bills.payFromAccount")}</Label>
              <Select value={payAccountId} onValueChange={setPayAccountId}>
                <SelectTrigger><SelectValue placeholder={t("bills.selectAccount")} /></SelectTrigger>
                <SelectContent>
                  {accounts.map((a) => (
                    <SelectItem key={a.id} value={a.id}>{a.name}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="pay-date">{t("bills.date")}</Label>
              <Input id="pay-date" type="date" value={payDate} onChange={(e) => setPayDate(e.target.value)} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPayDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button onClick={handlePay} disabled={submitting || !payAccountId || !payDate}>
              {submitting ? t("bills.paying") : t("bills.pay")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
