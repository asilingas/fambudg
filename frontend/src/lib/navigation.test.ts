import { describe, it, expect } from "vitest"
import { getNavForRole } from "./navigation"

describe("getNavForRole", () => {
  it("returns all items for admin", () => {
    const items = getNavForRole("admin")
    const labels = items.map((i) => i.label)
    expect(labels).toContain("Dashboard")
    expect(labels).toContain("Users")
    expect(labels).toContain("Categories")
    expect(labels).toContain("Budgets")
    expect(labels).toContain("Bills")
    expect(labels).toContain("Transfers")
    expect(labels).toContain("Allowances")
    expect(labels).toContain("Import/Export")
    expect(labels).toContain("Search")
  })

  it("returns member items (no Users, no Allowances)", () => {
    const items = getNavForRole("member")
    const labels = items.map((i) => i.label)
    expect(labels).toContain("Dashboard")
    expect(labels).toContain("Budgets")
    expect(labels).toContain("Transfers")
    expect(labels).toContain("Import/Export")
    expect(labels).not.toContain("Users")
    expect(labels).not.toContain("Allowances")
  })

  it("returns child items (limited access)", () => {
    const items = getNavForRole("child")
    const labels = items.map((i) => i.label)
    expect(labels).toContain("Dashboard")
    expect(labels).toContain("Transactions")
    expect(labels).toContain("Accounts")
    expect(labels).toContain("Categories")
    expect(labels).toContain("Reports")
    expect(labels).toContain("Search")
    expect(labels).toContain("Allowances")
    expect(labels).not.toContain("Budgets")
    expect(labels).not.toContain("Goals")
    expect(labels).not.toContain("Bills")
    expect(labels).not.toContain("Transfers")
    expect(labels).not.toContain("Import/Export")
    expect(labels).not.toContain("Users")
  })
})
