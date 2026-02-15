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
  })

  it("returns member items (no Users)", () => {
    const items = getNavForRole("member")
    const labels = items.map((i) => i.label)
    expect(labels).toContain("Dashboard")
    expect(labels).toContain("Budgets")
    expect(labels).not.toContain("Users")
  })

  it("returns child items (no Budgets, Goals, Bills, Users)", () => {
    const items = getNavForRole("child")
    const labels = items.map((i) => i.label)
    expect(labels).toContain("Dashboard")
    expect(labels).toContain("Transactions")
    expect(labels).toContain("Accounts")
    expect(labels).toContain("Categories")
    expect(labels).toContain("Reports")
    expect(labels).not.toContain("Budgets")
    expect(labels).not.toContain("Goals")
    expect(labels).not.toContain("Bills")
    expect(labels).not.toContain("Users")
  })
})
