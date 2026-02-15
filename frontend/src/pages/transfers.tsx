import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
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
import { ArrowRight } from "lucide-react"
import { toast } from "sonner"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { Account } from "@/lib/types"

export default function TransfersPage() {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)

  const [fromAccountId, setFromAccountId] = useState("")
  const [toAccountId, setToAccountId] = useState("")
  const [amount, setAmount] = useState("")
  const [description, setDescription] = useState("")
  const [date, setDate] = useState(new Date().toISOString().slice(0, 10))

  useEffect(() => {
    api.get("/accounts").then((res) => {
      setAccounts(res.data ?? [])
      setLoading(false)
    })
  }, [])

  async function handleTransfer() {
    if (fromAccountId === toAccountId) {
      toast.error("Cannot transfer to the same account")
      return
    }

    setSubmitting(true)
    try {
      await api.post("/transfers", {
        fromAccountId,
        toAccountId,
        amount: Math.round(parseFloat(amount) * 100),
        description: description || undefined,
        date,
      })
      toast.success("Transfer completed")
      setAmount("")
      setDescription("")
      // Refresh account balances
      const res = await api.get("/accounts")
      setAccounts(res.data ?? [])
    } catch {
      toast.error("Transfer failed")
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return <p className="text-muted-foreground">Loading...</p>
  }

  const fromAccount = accounts.find((a) => a.id === fromAccountId)
  const toAccount = accounts.find((a) => a.id === toAccountId)
  const canSubmit = fromAccountId && toAccountId && amount && date && fromAccountId !== toAccountId

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Transfer</h1>

      <Card className="max-w-lg">
        <CardHeader>
          <CardTitle className="text-sm font-medium">Move money between accounts</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-end gap-3">
            <div className="flex-1 space-y-2">
              <Label>From</Label>
              <Select value={fromAccountId} onValueChange={setFromAccountId}>
                <SelectTrigger>
                  <SelectValue placeholder="Select account" />
                </SelectTrigger>
                <SelectContent>
                  {accounts.map((a) => (
                    <SelectItem key={a.id} value={a.id}>
                      {a.name} ({formatCents(a.balance)})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <ArrowRight className="mb-2 h-5 w-5 text-muted-foreground shrink-0" />
            <div className="flex-1 space-y-2">
              <Label>To</Label>
              <Select value={toAccountId} onValueChange={setToAccountId}>
                <SelectTrigger>
                  <SelectValue placeholder="Select account" />
                </SelectTrigger>
                <SelectContent>
                  {accounts.map((a) => (
                    <SelectItem key={a.id} value={a.id}>
                      {a.name} ({formatCents(a.balance)})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="transfer-amount">Amount</Label>
            <Input
              id="transfer-amount"
              type="number"
              step="0.01"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="e.g. 500.00"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="transfer-desc">Description (optional)</Label>
            <Input
              id="transfer-desc"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="e.g. Monthly savings"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="transfer-date">Date</Label>
            <Input
              id="transfer-date"
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
            />
          </div>

          {fromAccount && toAccount && fromAccountId !== toAccountId && amount && (
            <div className="rounded-md bg-muted p-3 text-sm">
              Transfer {formatCents(Math.round(parseFloat(amount) * 100))} from{" "}
              <span className="font-medium">{fromAccount.name}</span> to{" "}
              <span className="font-medium">{toAccount.name}</span>
            </div>
          )}

          {fromAccountId && toAccountId && fromAccountId === toAccountId && (
            <p className="text-sm text-destructive">Cannot transfer to the same account.</p>
          )}

          <Button onClick={handleTransfer} disabled={submitting || !canSubmit} className="w-full">
            {submitting ? "Transferring..." : "Transfer"}
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
