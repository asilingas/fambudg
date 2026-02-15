import type { Role } from "@/context/auth-context"

export interface NavItem {
  label: string
  path: string
  icon: string
  roles: Role[]
}

export const navItems: NavItem[] = [
  { label: "Dashboard", path: "/", icon: "LayoutDashboard", roles: ["admin", "member", "child"] },
  { label: "Transactions", path: "/transactions", icon: "ArrowLeftRight", roles: ["admin", "member", "child"] },
  { label: "Accounts", path: "/accounts", icon: "Wallet", roles: ["admin", "member", "child"] },
  { label: "Categories", path: "/categories", icon: "Tag", roles: ["admin", "member", "child"] },
  { label: "Budgets", path: "/budgets", icon: "PiggyBank", roles: ["admin", "member"] },
  { label: "Reports", path: "/reports", icon: "BarChart3", roles: ["admin", "member", "child"] },
  { label: "Goals", path: "/goals", icon: "Target", roles: ["admin", "member"] },
  { label: "Bills", path: "/bills", icon: "Receipt", roles: ["admin", "member"] },
  { label: "Transfers", path: "/transfers", icon: "ArrowRightLeft", roles: ["admin", "member"] },
  { label: "Search", path: "/search", icon: "Search", roles: ["admin", "member", "child"] },
  { label: "Users", path: "/users", icon: "Users", roles: ["admin"] },
]

export function getNavForRole(role: Role): NavItem[] {
  return navItems.filter((item) => item.roles.includes(role))
}
