import { describe, it, expect } from "vitest"
import { getNavForRole } from "./navigation"

describe("getNavForRole", () => {
  it("returns all items for admin", () => {
    const items = getNavForRole("admin")
    const keys = items.map((i) => i.labelKey)
    expect(keys).toContain("nav.dashboard")
    expect(keys).toContain("nav.users")
    expect(keys).toContain("nav.categories")
    expect(keys).toContain("nav.budgets")
    expect(keys).toContain("nav.bills")
    expect(keys).toContain("nav.transfers")
    expect(keys).toContain("nav.allowances")
    expect(keys).toContain("nav.importExport")
    expect(keys).toContain("nav.search")
  })

  it("returns member items (no Users, no Allowances)", () => {
    const items = getNavForRole("member")
    const keys = items.map((i) => i.labelKey)
    expect(keys).toContain("nav.dashboard")
    expect(keys).toContain("nav.budgets")
    expect(keys).toContain("nav.transfers")
    expect(keys).toContain("nav.importExport")
    expect(keys).not.toContain("nav.users")
    expect(keys).not.toContain("nav.allowances")
  })

  it("returns child items (limited access)", () => {
    const items = getNavForRole("child")
    const keys = items.map((i) => i.labelKey)
    expect(keys).toContain("nav.dashboard")
    expect(keys).toContain("nav.transactions")
    expect(keys).toContain("nav.accounts")
    expect(keys).toContain("nav.categories")
    expect(keys).toContain("nav.reports")
    expect(keys).toContain("nav.search")
    expect(keys).toContain("nav.allowances")
    expect(keys).not.toContain("nav.budgets")
    expect(keys).not.toContain("nav.goals")
    expect(keys).not.toContain("nav.bills")
    expect(keys).not.toContain("nav.transfers")
    expect(keys).not.toContain("nav.importExport")
    expect(keys).not.toContain("nav.users")
  })
})
