import type { Role } from "@/context/auth-context"
import type { TranslationKey } from "@/lib/i18n"

export interface NavItem {
  labelKey: TranslationKey
  path: string
  icon: string
  roles: Role[]
}

export const navItems: NavItem[] = [
  { labelKey: "nav.dashboard", path: "/", icon: "LayoutDashboard", roles: ["admin", "member", "child"] },
  { labelKey: "nav.transactions", path: "/transactions", icon: "ArrowLeftRight", roles: ["admin", "member", "child"] },
  { labelKey: "nav.accounts", path: "/accounts", icon: "Wallet", roles: ["admin", "member", "child"] },
  { labelKey: "nav.categories", path: "/categories", icon: "Tag", roles: ["admin", "member", "child"] },
  { labelKey: "nav.budgets", path: "/budgets", icon: "PiggyBank", roles: ["admin", "member"] },
  { labelKey: "nav.reports", path: "/reports", icon: "BarChart3", roles: ["admin", "member", "child"] },
  { labelKey: "nav.goals", path: "/goals", icon: "Target", roles: ["admin", "member"] },
  { labelKey: "nav.bills", path: "/bills", icon: "Receipt", roles: ["admin", "member"] },
  { labelKey: "nav.transfers", path: "/transfers", icon: "ArrowRightLeft", roles: ["admin", "member"] },
  { labelKey: "nav.allowances", path: "/allowances", icon: "Coins", roles: ["admin", "child"] },
  { labelKey: "nav.importExport", path: "/import-export", icon: "FileSpreadsheet", roles: ["admin", "member"] },
  { labelKey: "nav.search", path: "/search", icon: "Search", roles: ["admin", "member", "child"] },
  { labelKey: "nav.users", path: "/users", icon: "Users", roles: ["admin"] },
]

export function getNavForRole(role: Role): NavItem[] {
  return navItems.filter((item) => item.roles.includes(role))
}
