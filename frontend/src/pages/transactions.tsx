import { useEffect, useState, useCallback } from "react"
import { PageSkeleton } from "@/components/loading-skeleton"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
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
import { Plus, Pencil, Trash2 } from "lucide-react"
import { toast } from "sonner"
import api from "@/lib/api"
import { formatCents, inputToCents, centsToInput } from "@/lib/format"
import type { Transaction, Account, Category } from "@/lib/types"
import { useLanguage } from "@/context/language-context"

interface FormData {
  amount: string
  type: "income" | "expense"
  accountId: string
  categoryId: string
  description: string
  date: string
  isShared: boolean
  tags: string
}

const emptyForm: FormData = {
  amount: "",
  type: "expense",
  accountId: "",
  categoryId: "",
  description: "",
  date: new Date().toISOString().split("T")[0],
  isShared: true,
  tags: "",
}

export default function TransactionsPage() {
  const { t } = useLanguage()
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Transaction | null>(null)
  const [deleting, setDeleting] = useState<Transaction | null>(null)
  const [form, setForm] = useState<FormData>(emptyForm)
  const [submitting, setSubmitting] = useState(false)

  // Filters
  const [filterAccount, setFilterAccount] = useState("")
  const [filterCategory, setFilterCategory] = useState("")
  const [filterType, setFilterType] = useState("")
  const [filterStartDate, setFilterStartDate] = useState("")
  const [filterEndDate, setFilterEndDate] = useState("")

  const fetchTransactions = useCallback(() => {
    const params: Record<string, string> = {}
    if (filterAccount) params.accountId = filterAccount
    if (filterCategory) params.categoryId = filterCategory
    if (filterType) params.type = filterType
    if (filterStartDate) params.startDate = filterStartDate
    if (filterEndDate) params.endDate = filterEndDate

    api.get("/transactions", { params }).then((res) => {
      setTransactions(res.data)
      setLoading(false)
    })
  }, [filterAccount, filterCategory, filterType, filterStartDate, filterEndDate])

  useEffect(() => {
    fetchTransactions()
  }, [fetchTransactions])

  useEffect(() => {
    Promise.all([api.get("/accounts"), api.get("/categories")]).then(
      ([accRes, catRes]) => {
        setAccounts(accRes.data)
        setCategories(catRes.data)
      },
    )
  }, [])

  function getCategoryName(id: string) {
    return categories.find((c) => c.id === id)?.name ?? "—"
  }

  function getAccountName(id: string) {
    return accounts.find((a) => a.id === id)?.name ?? "—"
  }

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDialogOpen(true)
  }

  function openEdit(tx: Transaction) {
    setEditing(tx)
    setForm({
      amount: centsToInput(Math.abs(tx.amount)),
      type: tx.type === "income" ? "income" : "expense",
      accountId: tx.accountId,
      categoryId: tx.categoryId,
      description: tx.description ?? "",
      date: tx.date.split("T")[0],
      isShared: tx.isShared,
      tags: tx.tags?.join(", ") ?? "",
    })
    setDialogOpen(true)
  }

  function openDelete(tx: Transaction) {
    setDeleting(tx)
    setDeleteDialogOpen(true)
  }

  async function handleSubmit() {
    setSubmitting(true)
    try {
      const cents = inputToCents(form.amount)
      const amount = form.type === "expense" ? -Math.abs(cents) : Math.abs(cents)
      const tags = form.tags
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean)

      const payload = {
        amount,
        type: form.type,
        accountId: form.accountId,
        categoryId: form.categoryId,
        description: form.description,
        date: form.date,
        isShared: form.isShared,
        tags: tags.length > 0 ? tags : undefined,
      }

      if (editing) {
        await api.put(`/transactions/${editing.id}`, payload)
        toast.success(t("transactions.updated"))
      } else {
        await api.post("/transactions", payload)
        toast.success(t("transactions.created"))
      }
      setDialogOpen(false)
      fetchTransactions()
    } catch {
      toast.error(t("transactions.saveFailed"))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    setSubmitting(true)
    try {
      await api.delete(`/transactions/${deleting.id}`)
      toast.success(t("transactions.deleted"))
      setDeleteDialogOpen(false)
      setDeleting(null)
      fetchTransactions()
    } catch {
      toast.error(t("transactions.deleteFailed"))
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
        <h1 className="text-2xl font-bold">{t("transactions.title")}</h1>
        <Button onClick={openCreate} size="sm">
          <Plus className="mr-1 h-4 w-4" />
          {t("transactions.add")}
        </Button>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-3">
        <Select value={filterAccount} onValueChange={setFilterAccount}>
          <SelectTrigger className="w-40">
            <SelectValue placeholder={t("transactions.allAccounts")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("transactions.allAccounts")}</SelectItem>
            {accounts.map((a) => (
              <SelectItem key={a.id} value={a.id}>
                {a.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select value={filterCategory} onValueChange={setFilterCategory}>
          <SelectTrigger className="w-40">
            <SelectValue placeholder={t("transactions.allCategories")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("transactions.allCategories")}</SelectItem>
            {categories.map((c) => (
              <SelectItem key={c.id} value={c.id}>
                {c.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select value={filterType} onValueChange={setFilterType}>
          <SelectTrigger className="w-36">
            <SelectValue placeholder={t("transactions.allTypes")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("transactions.allTypes")}</SelectItem>
            <SelectItem value="income">{t("transactions.income")}</SelectItem>
            <SelectItem value="expense">{t("transactions.expense")}</SelectItem>
            <SelectItem value="transfer">{t("transactions.transfer")}</SelectItem>
          </SelectContent>
        </Select>

        <Input
          type="date"
          className="w-40"
          value={filterStartDate}
          onChange={(e) => setFilterStartDate(e.target.value)}
          placeholder="Start date"
        />
        <Input
          type="date"
          className="w-40"
          value={filterEndDate}
          onChange={(e) => setFilterEndDate(e.target.value)}
          placeholder="End date"
        />
      </div>

      {/* Table */}
      {transactions.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t("transactions.noData")}</p>
      ) : (
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("transactions.date")}</TableHead>
                <TableHead>{t("transactions.description")}</TableHead>
                <TableHead>{t("transactions.category")}</TableHead>
                <TableHead>{t("transactions.account")}</TableHead>
                <TableHead className="text-right">{t("transactions.amount")}</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.map((tx) => (
                <TableRow key={tx.id}>
                  <TableCell className="whitespace-nowrap text-sm">
                    {new Date(tx.date).toLocaleDateString()}
                  </TableCell>
                  <TableCell className="text-sm">
                    {tx.description || "—"}
                    {tx.tags && tx.tags.length > 0 && (
                      <span className="ml-2">
                        {tx.tags.map((tag) => (
                          <Badge key={tag} variant="secondary" className="mr-1 text-xs">
                            {tag}
                          </Badge>
                        ))}
                      </span>
                    )}
                  </TableCell>
                  <TableCell className="text-sm">
                    {getCategoryName(tx.categoryId)}
                  </TableCell>
                  <TableCell className="text-sm">
                    {getAccountName(tx.accountId)}
                  </TableCell>
                  <TableCell
                    className={`text-right text-sm font-medium ${tx.amount >= 0 ? "text-income" : "text-expense"}`}
                  >
                    {formatCents(tx.amount)}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => openEdit(tx)}
                      >
                        <Pencil className="h-3 w-3" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => openDelete(tx)}
                      >
                        <Trash2 className="h-3 w-3" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Create/Edit Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editing ? t("transactions.editTitle") : t("transactions.newTitle")}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="tx-type">{t("transactions.type")}</Label>
                <Select
                  value={form.type}
                  onValueChange={(v) =>
                    setForm({ ...form, type: v as "income" | "expense", categoryId: "" })
                  }
                >
                  <SelectTrigger id="tx-type">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="income">{t("transactions.income")}</SelectItem>
                    <SelectItem value="expense">{t("transactions.expense")}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="tx-amount">{t("transactions.amount")}</Label>
                <Input
                  id="tx-amount"
                  type="number"
                  step="0.01"
                  min="0"
                  value={form.amount}
                  onChange={(e) => setForm({ ...form, amount: e.target.value })}
                  placeholder="0.00"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="tx-account">{t("transactions.account")}</Label>
              <Select
                value={form.accountId}
                onValueChange={(v) => setForm({ ...form, accountId: v })}
              >
                <SelectTrigger id="tx-account">
                  <SelectValue placeholder={t("transactions.selectAccount")} />
                </SelectTrigger>
                <SelectContent>
                  {accounts.map((a) => (
                    <SelectItem key={a.id} value={a.id}>
                      {a.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="tx-category">{t("transactions.category")}</Label>
              <Select
                value={form.categoryId}
                onValueChange={(v) => setForm({ ...form, categoryId: v })}
              >
                <SelectTrigger id="tx-category">
                  <SelectValue placeholder={t("transactions.selectCategory")} />
                </SelectTrigger>
                <SelectContent>
                  {categories
                    .filter((c) => c.type === form.type)
                    .map((c) => (
                      <SelectItem key={c.id} value={c.id}>
                        {c.name}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="tx-desc">{t("transactions.description")}</Label>
              <Input
                id="tx-desc"
                value={form.description}
                onChange={(e) =>
                  setForm({ ...form, description: e.target.value })
                }
                placeholder="e.g. Weekly groceries"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="tx-date">{t("transactions.date")}</Label>
              <Input
                id="tx-date"
                type="date"
                value={form.date}
                onChange={(e) => setForm({ ...form, date: e.target.value })}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="tx-tags">{t("transactions.tags")}</Label>
              <Input
                id="tx-tags"
                value={form.tags}
                onChange={(e) => setForm({ ...form, tags: e.target.value })}
                placeholder="e.g. vacation, birthday"
              />
            </div>

            {form.type === "expense" && (
              <div className="flex items-center gap-2">
                <input
                  id="tx-shared"
                  type="checkbox"
                  checked={form.isShared}
                  onChange={(e) =>
                    setForm({ ...form, isShared: e.target.checked })
                  }
                  className="h-4 w-4 rounded border-input"
                />
                <Label htmlFor="tx-shared">{t("transactions.sharedExpense")}</Label>
              </div>
            )}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>
              {t("common.cancel")}
            </Button>
            <Button
              onClick={handleSubmit}
              disabled={
                submitting || !form.amount || !form.accountId || !form.categoryId
              }
            >
              {submitting ? t("common.saving") : t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("transactions.deleteTitle")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            {t("transactions.deleteConfirm")}
          </p>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
            >
              {t("common.cancel")}
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={submitting}
            >
              {submitting ? t("common.deleting") : t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
