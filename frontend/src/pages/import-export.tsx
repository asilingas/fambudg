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
import { useLanguage } from "@/context/language-context"

export default function ImportExportPage() {
  const { t } = useLanguage()
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
      toast.success(t("importExport.exported"))
    } catch {
      toast.error(t("importExport.exportFailed"))
    } finally {
      setExporting(false)
    }
  }

  async function handleImport() {
    const file = fileInputRef.current?.files?.[0]
    if (!file) {
      toast.error(t("importExport.selectFile"))
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
      setImportResult(t("importExport.importSuccess").replace("{count}", String(count)).replace("{plural}", count !== 1 ? "s" : ""))
      toast.success(t("importExport.imported"))
      if (fileInputRef.current) fileInputRef.current.value = ""
    } catch {
      toast.error(t("importExport.importFailed"))
    } finally {
      setImporting(false)
    }
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{t("importExport.title")}</h1>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Export Card */}
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">{t("importExport.exportTitle")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              {t("importExport.exportDescription")}
            </p>
            <Button onClick={handleExport} disabled={exporting}>
              <Download className="mr-1 h-4 w-4" />
              {exporting ? t("importExport.exporting") : t("importExport.exportCsv")}
            </Button>
          </CardContent>
        </Card>

        {/* Import Card */}
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-medium">{t("importExport.importTitle")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              {t("importExport.importDescription")}
            </p>
            <div className="space-y-2">
              <Label htmlFor="csv-file">{t("importExport.csvFile")}</Label>
              <Input id="csv-file" type="file" accept=".csv" ref={fileInputRef} />
            </div>
            <Button onClick={handleImport} disabled={importing}>
              <Upload className="mr-1 h-4 w-4" />
              {importing ? t("importExport.importing") : t("importExport.importCsv")}
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
