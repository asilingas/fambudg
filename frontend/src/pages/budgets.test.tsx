import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "@/test/test-utils"
import BudgetsPage from "./budgets"

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

const budgets = [
  { id: "b1", categoryId: "c1", amount: 50000, month: 2, year: 2026, createdAt: "2026-02-01" },
]

const summary = [
  {
    categoryId: "c1",
    categoryName: "Groceries",
    budgetAmount: 50000,
    actualAmount: -35000,
    remaining: 15000,
  },
]

const categories = [
  { id: "c1", name: "Groceries", type: "expense", sortOrder: 0 },
  { id: "c2", name: "Salary", type: "income", sortOrder: 0 },
]

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url.startsWith("/budgets/summary")) return Promise.resolve({ data: summary })
    if (url.startsWith("/budgets")) return Promise.resolve({ data: budgets })
    if (url === "/categories") return Promise.resolve({ data: categories })
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("BudgetsPage", () => {
  it("renders budget cards with progress", async () => {
    renderWithProviders(<BudgetsPage />)

    await waitFor(() => {
      expect(screen.getByText("Groceries")).toBeInTheDocument()
    })
    expect(screen.getByText(/350,00.*spent/)).toBeInTheDocument()
    expect(screen.getByText(/500,00/)).toBeInTheDocument()
    expect(screen.getByText(/150,00.*remaining/)).toBeInTheDocument()
  })

  it("shows overspent warning when over budget", async () => {
    const overspentSummary = [
      {
        categoryId: "c1",
        categoryName: "Groceries",
        budgetAmount: 30000,
        actualAmount: -45000,
        remaining: -15000,
      },
    ]
    mockedApi.get.mockImplementation((url: string) => {
      if (url.startsWith("/budgets/summary")) return Promise.resolve({ data: overspentSummary })
      if (url.startsWith("/budgets")) return Promise.resolve({ data: budgets })
      if (url === "/categories") return Promise.resolve({ data: categories })
      if (url === "/auth/me") {
        return Promise.resolve({
          data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
        })
      }
      return Promise.reject(new Error("not mocked"))
    })

    renderWithProviders(<BudgetsPage />)

    await waitFor(() => {
      expect(screen.getByText(/Overspent by/)).toBeInTheDocument()
    })
  })

  it("shows add button for admin", async () => {
    renderWithProviders(<BudgetsPage />)

    await waitFor(() => {
      expect(screen.getByText("Groceries")).toBeInTheDocument()
    })
    expect(screen.getByRole("button", { name: /add budget/i })).toBeInTheDocument()
  })

  it("opens create dialog on add button click", async () => {
    const user = userEvent.setup()
    renderWithProviders(<BudgetsPage />)

    await waitFor(() => {
      expect(screen.getByText("Groceries")).toBeInTheDocument()
    })

    await user.click(screen.getByRole("button", { name: /add budget/i }))
    expect(screen.getByText("New Budget")).toBeInTheDocument()
  })

  it("shows empty state when no budgets", async () => {
    mockedApi.get.mockImplementation((url: string) => {
      if (url.startsWith("/budgets/summary")) return Promise.resolve({ data: [] })
      if (url.startsWith("/budgets")) return Promise.resolve({ data: [] })
      if (url === "/categories") return Promise.resolve({ data: categories })
      if (url === "/auth/me") {
        return Promise.resolve({
          data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
        })
      }
      return Promise.reject(new Error("not mocked"))
    })

    renderWithProviders(<BudgetsPage />)

    await waitFor(() => {
      expect(screen.getByText(/no budgets set/i)).toBeInTheDocument()
    })
  })
})
