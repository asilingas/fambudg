import { NavLink } from "react-router-dom"
import {
  LayoutDashboard,
  ArrowLeftRight,
  Wallet,
  Tag,
  PiggyBank,
  BarChart3,
  Target,
  Receipt,
  ArrowRightLeft,
  Search,
  Users,
  type LucideIcon,
} from "lucide-react"
import { useAuth } from "@/context/auth-context"
import { getNavForRole } from "@/lib/navigation"
import { cn } from "@/lib/utils"

const iconMap: Record<string, LucideIcon> = {
  LayoutDashboard,
  ArrowLeftRight,
  Wallet,
  Tag,
  PiggyBank,
  BarChart3,
  Target,
  Receipt,
  ArrowRightLeft,
  Search,
  Users,
}

export function BottomTabs() {
  const { user } = useAuth()
  if (!user) return null

  const items = getNavForRole(user.role).slice(0, 5)

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 flex border-t bg-background md:hidden">
      {items.map((item) => {
        const Icon = iconMap[item.icon]
        return (
          <NavLink
            key={item.path}
            to={item.path}
            end={item.path === "/"}
            className={({ isActive }) =>
              cn(
                "flex flex-1 flex-col items-center gap-0.5 py-2 text-[10px]",
                isActive
                  ? "text-primary"
                  : "text-muted-foreground",
              )
            }
          >
            {Icon && <Icon className="h-5 w-5" />}
            {item.label}
          </NavLink>
        )
      })}
    </nav>
  )
}
