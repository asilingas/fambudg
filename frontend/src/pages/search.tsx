import { useState } from "react"
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Badge } from "@/components/ui/badge"
import { Search } from "lucide-react"
import api from "@/lib/api"
import { formatCents } from "@/lib/format"
import type { Transaction, Category, Account } from "@/lib/types"
import { useEffect } from "react"
import { useLanguage } from "@/context/language-context"

export default function SearchPage() {
  const [description, setDescription] = useState("")
  const [startDate, setStartDate] = useState("")
  const [endDate, setEndDate] = useState("")
  const [categoryId, setCategoryId] = useState("")
  const [accountId, setAccountId] = useState("")
  const [minAmount, setMinAmount] = useState("")
  const [maxAmount, setMaxAmount] = useState("")
  const [tags, setTags] = useState("")

  const [results, setResults] = useState<Transaction[]>([])
  const [totalCount, setTotalCount] = useState(0)
  const [searched, setSearched] = useState(false)
  const [loading, setLoading] = useState(false)

  const [categories, setCategories] = useState<Category[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])

  const { t } = useLanguage()

  useEffect(() => {
    Promise.all([api.get("/categories"), api.get("/accounts")]).then(
      ([catsRes, acctsRes]) => {
        setCategories(catsRes.data ?? [])
        setAccounts(acctsRes.data ?? [])
      }
    )
  }, [])

  function categoryName(id: string) {
    return categories.find((c) => c.id === id)?.name ?? ""
  }

  function accountName(id: string) {
    return accounts.find((a) => a.id === id)?.name ?? ""
  }

  async function handleSearch() {
    setLoading(true)
    const params = new URLSearchParams()
    if (description) params.set("description", description)
    if (startDate) params.set("startDate", startDate)
    if (endDate) params.set("endDate", endDate)
    if (categoryId) params.set("categoryId", categoryId)
    if (accountId) params.set("accountId", accountId)
    if (minAmount) params.set("minAmount", String(Math.round(parseFloat(minAmount) * 100)))
    if (maxAmount) params.set("maxAmount", String(Math.round(parseFloat(maxAmount) * 100)))
    if (tags) params.set("tags", tags)

    try {
      const res = await api.get(`/search?${params.toString()}`)
      setResults(res.data.transactions ?? [])
      setTotalCount(res.data.totalCount ?? 0)
      setSearched(true)
    } catch {
      setResults([])
      setTotalCount(0)
    } finally {
      setLoading(false)
    }
  }

  function handleClear() {
    setDescription("")
    setStartDate("")
    setEndDate("")
    setCategoryId("")
    setAccountId("")
    setMinAmount("")
    setMaxAmount("")
    setTags("")
    setResults([])
    setTotalCount(0)
    setSearched(false)
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{t("search.title")}</h1>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div className="space-y-2">
          <Label htmlFor="search-desc">{t("search.description")}</Label>
          <Input
            id="search-desc"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="e.g. groceries"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="search-start">{t("search.startDate")}</Label>
          <Input
            id="search-start"
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="search-end">{t("search.endDate")}</Label>
          <Input
            id="search-end"
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
          />
        </div>
        <div className="space-y-2">
          <Label>{t("search.category")}</Label>
          <Select value={categoryId} onValueChange={setCategoryId}>
            <SelectTrigger>
              <SelectValue placeholder={t("search.all")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("search.all")}</SelectItem>
              {categories.map((c) => (
                <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label>{t("search.account")}</Label>
          <Select value={accountId} onValueChange={setAccountId}>
            <SelectTrigger>
              <SelectValue placeholder={t("search.all")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("search.all")}</SelectItem>
              {accounts.map((a) => (
                <SelectItem key={a.id} value={a.id}>{a.name}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="search-min">{t("search.minAmount")}</Label>
          <Input
            id="search-min"
            type="number"
            step="0.01"
            value={minAmount}
            onChange={(e) => setMinAmount(e.target.value)}
            placeholder="0.00"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="search-max">{t("search.maxAmount")}</Label>
          <Input
            id="search-max"
            type="number"
            step="0.01"
            value={maxAmount}
            onChange={(e) => setMaxAmount(e.target.value)}
            placeholder="0.00"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="search-tags">{t("search.tags")}</Label>
          <Input
            id="search-tags"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
            placeholder="tag1,tag2"
          />
        </div>
      </div>

      <div className="flex gap-2">
        <Button onClick={handleSearch} disabled={loading}>
          <Search className="mr-1 h-4 w-4" />
          {loading ? t("search.searching") : t("search.submit")}
        </Button>
        <Button variant="outline" onClick={handleClear}>
          {t("search.clear")}
        </Button>
      </div>

      {searched && (
        <div className="space-y-2">
          <p className="text-sm text-muted-foreground">
            {t("search.resultCount").replace("{count}", String(totalCount)).replace("{plural}", totalCount !== 1 ? "s" : "")}
          </p>

          {results.length > 0 ? (
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t("search.date")}</TableHead>
                    <TableHead>{t("search.description")}</TableHead>
                    <TableHead>{t("search.category")}</TableHead>
                    <TableHead>{t("search.account")}</TableHead>
                    <TableHead className="text-right">{t("search.amount")}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {results.map((t) => (
                    <TableRow key={t.id}>
                      <TableCell className="text-sm">{t.date.slice(0, 10)}</TableCell>
                      <TableCell>
                        <span className="text-sm">{t.description || "—"}</span>
                        {t.tags && t.tags.length > 0 && (
                          <div className="flex gap-1 mt-1">
                            {t.tags.map((tag) => (
                              <Badge key={tag} variant="outline" className="text-xs">
                                {tag}
                              </Badge>
                            ))}
                          </div>
                        )}
                      </TableCell>
                      <TableCell className="text-sm">{categoryName(t.categoryId)}</TableCell>
                      <TableCell className="text-sm">{accountName(t.accountId)}</TableCell>
                      <TableCell className={`text-right text-sm font-medium ${t.amount >= 0 ? "text-income" : "text-expense"}`}>
                        {formatCents(t.amount)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">{t("search.noResults")}</p>
          )}
        </div>
      )}
    </div>
  )
}
