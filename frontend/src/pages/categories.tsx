import { useEffect, useState, useCallback } from "react"
import { PageSkeleton } from "@/components/loading-skeleton"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
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
import { useAuth } from "@/context/auth-context"
import { useLanguage } from "@/context/language-context"
import api from "@/lib/api"
import type { Category } from "@/lib/types"

interface FormData {
  name: string
  type: "expense" | "income"
  icon: string
  sortOrder: string
}

const emptyForm: FormData = { name: "", type: "expense", icon: "", sortOrder: "0" }

export default function CategoriesPage() {
  const { user } = useAuth()
  const { t } = useLanguage()
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Category | null>(null)
  const [deleting, setDeleting] = useState<Category | null>(null)
  const [form, setForm] = useState<FormData>(emptyForm)
  const [submitting, setSubmitting] = useState(false)

  const canCreate = user?.role === "admin" || user?.role === "member"
  const canEditDelete = user?.role === "admin"

  const fetchCategories = useCallback(() => {
    api.get("/categories").then((res) => {
      setCategories(res.data)
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    fetchCategories()
  }, [fetchCategories])

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDialogOpen(true)
  }

  function openEdit(cat: Category) {
    setEditing(cat)
    setForm({
      name: cat.name,
      type: cat.type,
      icon: cat.icon ?? "",
      sortOrder: String(cat.sortOrder),
    })
    setDialogOpen(true)
  }

  function openDelete(cat: Category) {
    setDeleting(cat)
    setDeleteDialogOpen(true)
  }

  async function handleSubmit() {
    setSubmitting(true)
    try {
      const payload = {
        name: form.name,
        type: form.type,
        icon: form.icon || undefined,
        sortOrder: parseInt(form.sortOrder) || 0,
      }
      if (editing) {
        await api.put(`/categories/${editing.id}`, payload)
        toast.success(t("categories.updated"))
      } else {
        await api.post("/categories", payload)
        toast.success(t("categories.created"))
      }
      setDialogOpen(false)
      fetchCategories()
    } catch {
      toast.error(t("categories.saveFailed"))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    setSubmitting(true)
    try {
      await api.delete(`/categories/${deleting.id}`)
      toast.success(t("categories.deleted"))
      setDeleteDialogOpen(false)
      setDeleting(null)
      fetchCategories()
    } catch {
      toast.error(t("categories.deleteFailed"))
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
        <h1 className="text-2xl font-bold">{t("categories.title")}</h1>
        {canCreate && (
          <Button onClick={openCreate} size="sm">
            <Plus className="mr-1 h-4 w-4" />
            {t("categories.add")}
          </Button>
        )}
      </div>

      {categories.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t("categories.noData")}</p>
      ) : (
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("categories.name")}</TableHead>
                <TableHead>{t("categories.type")}</TableHead>
                <TableHead>{t("categories.icon")}</TableHead>
                <TableHead>{t("categories.order")}</TableHead>
                {canEditDelete && <TableHead />}
              </TableRow>
            </TableHeader>
            <TableBody>
              {categories.map((cat) => (
                <TableRow key={cat.id}>
                  <TableCell className="font-medium">{cat.name}</TableCell>
                  <TableCell>
                    <Badge
                      variant={cat.type === "income" ? "default" : "secondary"}
                    >
                      {cat.type}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-sm">
                    {cat.icon || "—"}
                  </TableCell>
                  <TableCell className="text-sm">{cat.sortOrder}</TableCell>
                  {canEditDelete && (
                    <TableCell>
                      <div className="flex gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => openEdit(cat)}
                        >
                          <Pencil className="h-3 w-3" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => openDelete(cat)}
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
                    </TableCell>
                  )}
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
              {editing ? t("categories.editTitle") : t("categories.newTitle")}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="cat-name">{t("categories.name")}</Label>
              <Input
                id="cat-name"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="e.g. Groceries"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cat-type">{t("categories.type")}</Label>
              <Select
                value={form.type}
                onValueChange={(v) =>
                  setForm({ ...form, type: v as "expense" | "income" })
                }
              >
                <SelectTrigger id="cat-type">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="expense">{t("transactions.expense")}</SelectItem>
                  <SelectItem value="income">{t("transactions.income")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="cat-icon">{t("categories.iconOptional")}</Label>
              <Input
                id="cat-icon"
                value={form.icon}
                onChange={(e) => setForm({ ...form, icon: e.target.value })}
                placeholder="e.g. shopping-cart"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cat-order">{t("categories.sortOrder")}</Label>
              <Input
                id="cat-order"
                type="number"
                value={form.sortOrder}
                onChange={(e) =>
                  setForm({ ...form, sortOrder: e.target.value })
                }
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDialogOpen(false)}>
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
            <DialogTitle>{t("categories.deleteTitle")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            {t("categories.deleteConfirm").replace("{name}", deleting?.name ?? "")}
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
