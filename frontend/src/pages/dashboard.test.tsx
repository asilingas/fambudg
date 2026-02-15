import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import { renderWithProviders } from "@/test/test-utils"
import DashboardPage from "./dashboard"

vi.mock("@/lib/api", () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  },
}))

import api from "@/lib/api"

const mockedApi = vi.mocked(api)

const dashboardData = {
  accounts: [
    { id: "1", userId: "u1", name: "Checking", type: "checking", currency: "EUR", balance: 150000 },
  ],
  monthSummary: {
    month: 2,
    year: 2026,
    totalIncome: 500000,
    totalExpense: -200000,
    net: 300000,
  },
  recentTransactions: [
    {
      id: "t1",
      userId: "u1",
      accountId: "1",
      categoryId: "c1",
      amount: -4599,
      type: "expense",
      description: "Groceries",
      date: "2026-02-10",
      isShared: true,
      isRecurring: false,
      createdAt: "2026-02-10T10:00:00Z",
      updatedAt: "2026-02-10T10:00:00Z",
    },
  ],
}

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (typeof url === "string" && url.startsWith("/reports/dashboard")) {
      return Promise.resolve({ data: dashboardData })
    }
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("DashboardPage", () => {
  it("renders income, expense, and net summary", async () => {
    renderWithProviders(<DashboardPage />)

    await waitFor(() => {
      expect(screen.getByText(/5\.000,00/)).toBeInTheDocument()
    })
    expect(screen.getByText(/2\.000,00/)).toBeInTheDocument()
    expect(screen.getByText(/3\.000,00/)).toBeInTheDocument()
  })

  it("renders account cards", async () => {
    renderWithProviders(<DashboardPage />)

    await waitFor(() => {
      expect(screen.getByText("Checking")).toBeInTheDocument()
    })
    expect(screen.getByText(/1\.500,00/)).toBeInTheDocument()
  })

  it("renders recent transactions", async () => {
    renderWithProviders(<DashboardPage />)

    await waitFor(() => {
      expect(screen.getByText("Groceries")).toBeInTheDocument()
    })
    expect(screen.getByText(/-45,99/)).toBeInTheDocument()
  })

  it("shows empty state when no accounts", async () => {
    mockedApi.get.mockImplementation((url: string) => {
      if (typeof url === "string" && url.startsWith("/reports/dashboard")) {
        return Promise.resolve({
          data: {
            accounts: [],
            monthSummary: { month: 2, year: 2026, totalIncome: 0, totalExpense: 0, net: 0 },
            recentTransactions: [],
          },
        })
      }
      return Promise.reject(new Error("not mocked"))
    })

    renderWithProviders(<DashboardPage />)

    await waitFor(() => {
      expect(screen.getByText(/no accounts yet/i)).toBeInTheDocument()
    })
  })
})
