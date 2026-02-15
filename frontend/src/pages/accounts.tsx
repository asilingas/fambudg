import { useEffect, useState, useCallback } from "react"
import { PageSkeleton } from "@/components/loading-skeleton"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
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
import type { Account } from "@/lib/types"
import { useLanguage } from "@/context/language-context"

const ACCOUNT_TYPES = ["checking", "savings", "credit", "cash"] as const

interface FormData {
  name: string
  type: string
  currency: string
  balance: string
}

const emptyForm: FormData = { name: "", type: "checking", currency: "EUR", balance: "0.00" }

export default function AccountsPage() {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Account | null>(null)
  const [deleting, setDeleting] = useState<Account | null>(null)
  const [form, setForm] = useState<FormData>(emptyForm)
  const [submitting, setSubmitting] = useState(false)
  const { t } = useLanguage()

  const fetchAccounts = useCallback(() => {
    api.get("/accounts").then((res) => {
      setAccounts(res.data)
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    fetchAccounts()
  }, [fetchAccounts])

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDialogOpen(true)
  }

  function openEdit(account: Account) {
    setEditing(account)
    setForm({
      name: account.name,
      type: account.type,
      currency: account.currency,
      balance: centsToInput(account.balance),
    })
    setDialogOpen(true)
  }

  function openDelete(account: Account) {
    setDeleting(account)
    setDeleteDialogOpen(true)
  }

  async function handleSubmit() {
    setSubmitting(true)
    try {
      const payload = {
        name: form.name,
        type: form.type,
        currency: form.currency,
        balance: inputToCents(form.balance),
      }
      if (editing) {
        await api.put(`/accounts/${editing.id}`, payload)
        toast.success(t("accounts.updated"))
      } else {
        await api.post("/accounts", payload)
        toast.success(t("accounts.created"))
      }
      setDialogOpen(false)
      fetchAccounts()
    } catch {
      toast.error(t("accounts.saveFailed"))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    setSubmitting(true)
    try {
      await api.delete(`/accounts/${deleting.id}`)
      toast.success(t("accounts.deleted"))
      setDeleteDialogOpen(false)
      setDeleting(null)
      fetchAccounts()
    } catch {
      toast.error(t("accounts.deleteFailed"))
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
        <h1 className="text-2xl font-bold">{t("accounts.title")}</h1>
        <Button onClick={openCreate} size="sm">
          <Plus className="mr-1 h-4 w-4" />
          {t("accounts.add")}
        </Button>
      </div>

      {accounts.length === 0 ? (
        <p className="text-sm text-muted-foreground">
          {t("accounts.noData")}
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
                <div className="mt-3 flex gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => openEdit(account)}
                  >
                    <Pencil className="mr-1 h-3 w-3" />
                    {t("common.edit")}
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => openDelete(account)}
                  >
                    <Trash2 className="mr-1 h-3 w-3" />
                    {t("common.delete")}
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Create/Edit Dialog */}
      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {editing ? t("accounts.editTitle") : t("accounts.newTitle")}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">{t("accounts.name")}</Label>
              <Input
                id="name"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="e.g. Chase Checking"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="type">{t("accounts.type")}</Label>
              <Select
                value={form.type}
                onValueChange={(v) => setForm({ ...form, type: v })}
              >
                <SelectTrigger id="type">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {ACCOUNT_TYPES.map((t) => (
                    <SelectItem key={t} value={t}>
                      {t.charAt(0).toUpperCase() + t.slice(1)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="currency">{t("accounts.currency")}</Label>
              <Input
                id="currency"
                value={form.currency}
                onChange={(e) =>
                  setForm({ ...form, currency: e.target.value.toUpperCase() })
                }
                maxLength={3}
                placeholder="EUR"
              />
            </div>
            {!editing && (
              <div className="space-y-2">
                <Label htmlFor="balance">{t("accounts.startingBalance")}</Label>
                <Input
                  id="balance"
                  type="number"
                  step="0.01"
                  value={form.balance}
                  onChange={(e) => setForm({ ...form, balance: e.target.value })}
                />
              </div>
            )}
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDialogOpen(false)}
            >
              {t("common.cancel")}
            </Button>
            <Button onClick={handleSubmit} disabled={submitting || !form.name}>
              {submitting ? t("common.saving") : t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("accounts.deleteTitle")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            {t("accounts.deleteConfirm").replace("{name}", deleting?.name ?? "")}
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
