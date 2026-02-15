import { useAuth } from "@/context/auth-context"
import { Button } from "@/components/ui/button"
import { LogOut } from "lucide-react"

export function TopBar() {
  const { user, logout } = useAuth()
  if (!user) return null

  return (
    <header className="flex h-14 items-center justify-between border-b px-4">
      <span className="text-lg font-semibold md:hidden">Fambudg</span>
      <div className="hidden md:block" />
      <div className="flex items-center gap-3">
        <span className="text-sm text-muted-foreground">
          {user.name}
          <span className="ml-1.5 rounded-full bg-secondary px-2 py-0.5 text-xs">
            {user.role}
          </span>
        </span>
        <Button variant="ghost" size="icon" onClick={logout} aria-label="Logout">
          <LogOut className="h-4 w-4" />
        </Button>
      </div>
    </header>
  )
}
