import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import { renderWithProviders } from "@/test/test-utils"
import BillsPage from "./bills"

vi.mock("@/lib/api", () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  },
}))

import api from "@/lib/api"

const mockedApi = vi.mocked(api, true)

const bills = [
  {
    id: "br1",
    name: "Rent",
    amount: 80000,
    dueDay: 1,
    frequency: "monthly",
    categoryId: "c1",
    accountId: "a1",
    isActive: true,
    nextDueDate: "2026-03-01T00:00:00Z",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-02-01T00:00:00Z",
  },
  {
    id: "br2",
    name: "Insurance",
    amount: 15000,
    dueDay: 15,
    frequency: "monthly",
    isActive: true,
    nextDueDate: "2026-02-10T00:00:00Z",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-02-01T00:00:00Z",
  },
]

const upcoming = [bills[1]] // Insurance is upcoming (past due)

const categories = [{ id: "c1", name: "Housing", type: "expense", sortOrder: 0 }]
const accounts = [{ id: "a1", userId: "u1", name: "Checking", type: "checking", currency: "EUR", balance: 100000, createdAt: "2026-01-01" }]

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url === "/bill-reminders") return Promise.resolve({ data: bills })
    if (url === "/bill-reminders/upcoming") return Promise.resolve({ data: upcoming })
    if (url === "/categories") return Promise.resolve({ data: categories })
    if (url === "/accounts") return Promise.resolve({ data: accounts })
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("BillsPage", () => {
  it("renders bill reminders table", async () => {
    renderWithProviders(<BillsPage />)

    await waitFor(() => {
      expect(screen.getByText("Rent")).toBeInTheDocument()
    })
    // Insurance appears in both upcoming cards and all bills table
    expect(screen.getAllByText("Insurance").length).toBeGreaterThanOrEqual(1)
  })

  it("shows upcoming bills section", async () => {
    renderWithProviders(<BillsPage />)

    await waitFor(() => {
      expect(screen.getByText("Upcoming")).toBeInTheDocument()
    })
  })

  it("shows overdue badge for past-due bills", async () => {
    renderWithProviders(<BillsPage />)

    await waitFor(() => {
      expect(screen.getByText("Overdue")).toBeInTheDocument()
    })
  })

  it("shows mark as paid button", async () => {
    renderWithProviders(<BillsPage />)

    await waitFor(() => {
      expect(screen.getByText("Rent")).toBeInTheDocument()
    })
    expect(screen.getByText("Mark as Paid")).toBeInTheDocument()
  })

  it("shows empty state when no bills", async () => {
    mockedApi.get.mockImplementation((url: string) => {
      if (url === "/bill-reminders") return Promise.resolve({ data: [] })
      if (url === "/bill-reminders/upcoming") return Promise.resolve({ data: [] })
      if (url === "/categories") return Promise.resolve({ data: categories })
      if (url === "/accounts") return Promise.resolve({ data: accounts })
      if (url === "/auth/me") {
        return Promise.resolve({
          data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
        })
      }
      return Promise.reject(new Error("not mocked"))
    })

    renderWithProviders(<BillsPage />)

    await waitFor(() => {
      expect(screen.getByText(/no bill reminders yet/i)).toBeInTheDocument()
    })
  })
})
