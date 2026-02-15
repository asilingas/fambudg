import { useState, useRef } from "react"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Download, Upload } from "lucide-react"
import { toast } from "sonner"
import api from "@/lib/api"

export default function ImportExportPage() {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [importing, setImporting] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [importResult, setImportResult] = useState<string | null>(null)

  async function handleExport() {
    setExporting(true)
    try {
      const res = await api.get("/export/csv", { responseType: "blob" })
      const url = window.URL.createObjectURL(new Blob([res.data]))
      const link = document.createElement("a")
      link.href = url
      link.setAttribute("download", "transactions.csv")
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
      toast.success("CSV exported")
    } catch {
      toast.error("Export failed")
    } finally {
      setExporting(false)
    }
  }

  async function handleImport() {
    const file = fileInputRef.current?.files?.[0]
    if (!file) {
      toast.error("Please select a CSV file")
      return
    }

    setImporting(true)
    setImportResult(null)
    try {
      const formData = new FormData()
      formData.append("file", file)
      const res = await api.post("/import/csv", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      })
      const count = res.data?.imported ?? res.data?.count ?? 0
      setImportResult(`Successfully imported ${count} transaction${count !== 1 ? "s" : ""}.`)
      toast.success("CSV imported")
      if (fileInputRef.current) fileInputRef.current.value = ""
    } catch {
      toast.error("Import failed")
    } finally {
      setImporting(false)
    }
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Import / Export</h1>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Export Card */}
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Export Transactions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Download all your transactions as a CSV file.
            </p>
            <Button onClick={handleExport} disabled={exporting}>
              <Download className="mr-1 h-4 w-4" />
              {exporting ? "Exporting..." : "Export CSV"}
            </Button>
          </CardContent>
        </Card>

        {/* Import Card */}
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">Import Transactions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Upload a CSV file to import transactions in bulk.
            </p>
            <div className="space-y-2">
              <Label htmlFor="csv-file">CSV File</Label>
              <Input id="csv-file" type="file" accept=".csv" ref={fileInputRef} />
            </div>
            <Button onClick={handleImport} disabled={importing}>
              <Upload className="mr-1 h-4 w-4" />
              {importing ? "Importing..." : "Import CSV"}
            </Button>
            {importResult && (
              <p className="text-sm text-income">{importResult}</p>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
