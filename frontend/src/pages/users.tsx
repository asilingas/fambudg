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
import api from "@/lib/api"
import type { User } from "@/lib/types"
import { useLanguage } from "@/context/language-context"

interface CreateFormData {
  email: string
  password: string
  name: string
  role: string
}

interface EditFormData {
  name: string
  role: string
}

const emptyCreateForm: CreateFormData = { email: "", password: "", name: "", role: "member" }
const emptyEditForm: EditFormData = { name: "", role: "" }

const roleBadgeVariant: Record<string, "default" | "secondary" | "outline"> = {
  admin: "default",
  member: "secondary",
  child: "outline",
}

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editDialogOpen, setEditDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)
  const [deleting, setDeleting] = useState<User | null>(null)
  const [createForm, setCreateForm] = useState<CreateFormData>(emptyCreateForm)
  const [editForm, setEditForm] = useState<EditFormData>(emptyEditForm)
  const [submitting, setSubmitting] = useState(false)
  const { t } = useLanguage()

  const fetchUsers = useCallback(() => {
    api.get("/users").then((res) => {
      setUsers(res.data ?? [])
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    fetchUsers()
  }, [fetchUsers])

  function openCreate() {
    setCreateForm(emptyCreateForm)
    setCreateDialogOpen(true)
  }

  function openEdit(user: User) {
    setEditing(user)
    setEditForm({ name: user.name, role: user.role })
    setEditDialogOpen(true)
  }

  function openDelete(user: User) {
    setDeleting(user)
    setDeleteDialogOpen(true)
  }

  async function handleCreate() {
    setSubmitting(true)
    try {
      await api.post("/users", createForm)
      toast.success(t("users.userCreated"))
      setCreateDialogOpen(false)
      fetchUsers()
    } catch {
      toast.error(t("users.createFailed"))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleEdit() {
    if (!editing) return
    setSubmitting(true)
    try {
      await api.put(`/users/${editing.id}`, editForm)
      toast.success(t("users.userUpdated"))
      setEditDialogOpen(false)
      fetchUsers()
    } catch {
      toast.error(t("users.updateFailed"))
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    setSubmitting(true)
    try {
      await api.delete(`/users/${deleting.id}`)
      toast.success(t("users.userDeleted"))
      setDeleteDialogOpen(false)
      setDeleting(null)
      fetchUsers()
    } catch {
      toast.error(t("users.deleteFailed"))
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
        <h1 className="text-2xl font-bold">{t("users.title")}</h1>
        <Button onClick={openCreate} size="sm">
          <Plus className="mr-1 h-4 w-4" />
          {t("users.add")}
        </Button>
      </div>

      {users.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t("users.noData")}</p>
      ) : (
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("users.name")}</TableHead>
                <TableHead>{t("users.email")}</TableHead>
                <TableHead>{t("users.role")}</TableHead>
                <TableHead>{t("users.created")}</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((u) => (
                <TableRow key={u.id}>
                  <TableCell className="font-medium">{u.name}</TableCell>
                  <TableCell className="text-sm">{u.email}</TableCell>
                  <TableCell>
                    <Badge variant={roleBadgeVariant[u.role] ?? "secondary"}>
                      {u.role}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-sm">{u.createdAt.slice(0, 10)}</TableCell>
                  <TableCell>
                    <div className="flex gap-1">
                      <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openEdit(u)}>
                        <Pencil className="h-3 w-3" />
                      </Button>
                      <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => openDelete(u)}>
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

      {/* Create User Dialog */}
      <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("users.createTitle")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="user-email">{t("users.email")}</Label>
              <Input id="user-email" type="email" value={createForm.email} onChange={(e) => setCreateForm({ ...createForm, email: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="user-password">{t("users.password")}</Label>
              <Input id="user-password" type="password" value={createForm.password} onChange={(e) => setCreateForm({ ...createForm, password: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="user-name">{t("users.name")}</Label>
              <Input id="user-name" value={createForm.name} onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>{t("users.role")}</Label>
              <Select value={createForm.role} onValueChange={(v) => setCreateForm({ ...createForm, role: v })}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="admin">{t("users.admin")}</SelectItem>
                  <SelectItem value="member">{t("users.member")}</SelectItem>
                  <SelectItem value="child">{t("users.child")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button onClick={handleCreate} disabled={submitting || !createForm.email || !createForm.password || !createForm.name}>
              {submitting ? t("users.creating") : t("users.create")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit User Dialog */}
      <Dialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("users.editTitle")}</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="edit-name">{t("users.name")}</Label>
              <Input id="edit-name" value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} />
            </div>
            <div className="space-y-2">
              <Label>{t("users.role")}</Label>
              <Select value={editForm.role} onValueChange={(v) => setEditForm({ ...editForm, role: v })}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="admin">{t("users.admin")}</SelectItem>
                  <SelectItem value="member">{t("users.member")}</SelectItem>
                  <SelectItem value="child">{t("users.child")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button onClick={handleEdit} disabled={submitting || !editForm.name}>
              {submitting ? t("common.saving") : t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("users.deleteTitle")}</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            {t("users.deleteConfirm").replace("{name}", deleting?.name ?? "").replace("{email}", deleting?.email ?? "")}
          </p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteDialogOpen(false)}>{t("common.cancel")}</Button>
            <Button variant="destructive" onClick={handleDelete} disabled={submitting}>
              {submitting ? t("common.deleting") : t("common.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
