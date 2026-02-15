import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import { renderWithProviders } from "@/test/test-utils"
import ReportsPage from "./reports"

// Mock recharts to avoid canvas rendering issues in jsdom
vi.mock("recharts", () => ({
  ResponsiveContainer: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  BarChart: ({ children }: { children: React.ReactNode }) => <div data-testid="bar-chart">{children}</div>,
  Bar: () => null,
  XAxis: () => null,
  YAxis: () => null,
  CartesianGrid: () => null,
  Tooltip: () => null,
  LineChart: ({ children }: { children: React.ReactNode }) => <div data-testid="line-chart">{children}</div>,
  Line: () => null,
  Legend: () => null,
}))

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

const monthlySummary = {
  month: 2,
  year: 2026,
  totalIncome: 500000,
  totalExpense: -200000,
  net: 300000,
}

const categorySpending = [
  { categoryId: "c1", categoryName: "Groceries", totalAmount: -120000, percentage: 60 },
  { categoryId: "c2", categoryName: "Utilities", totalAmount: -80000, percentage: 40 },
]

const trends = [
  { month: 1, year: 2026, totalIncome: 400000, totalExpense: -150000, net: 250000 },
  { month: 2, year: 2026, totalIncome: 500000, totalExpense: -200000, net: 300000 },
]

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url.startsWith("/reports/monthly")) return Promise.resolve({ data: monthlySummary })
    if (url.startsWith("/reports/by-category")) return Promise.resolve({ data: categorySpending })
    if (url.startsWith("/reports/trends")) return Promise.resolve({ data: trends })
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("ReportsPage", () => {
  it("renders monthly income, expense, and net", async () => {
    renderWithProviders(<ReportsPage />)

    await waitFor(() => {
      expect(screen.getByText(/5\.000,00/)).toBeInTheDocument()
    })
    expect(screen.getByText(/2\.000,00/)).toBeInTheDocument()
    expect(screen.getByText(/3\.000,00/)).toBeInTheDocument()
  })

  it("renders Reports heading", async () => {
    renderWithProviders(<ReportsPage />)

    await waitFor(() => {
      expect(screen.getByText("Reports")).toBeInTheDocument()
    })
  })

  it("renders tabs for Monthly, By Category, and Trends", async () => {
    renderWithProviders(<ReportsPage />)

    await waitFor(() => {
      expect(screen.getByRole("tab", { name: /monthly/i })).toBeInTheDocument()
    })
    expect(screen.getByRole("tab", { name: /by category/i })).toBeInTheDocument()
    expect(screen.getByRole("tab", { name: /trends/i })).toBeInTheDocument()
  })

  it("shows empty state when no monthly data", async () => {
    mockedApi.get.mockImplementation((url: string) => {
      if (url.startsWith("/reports/monthly")) return Promise.resolve({ data: null })
      if (url.startsWith("/reports/by-category")) return Promise.resolve({ data: [] })
      if (url.startsWith("/reports/trends")) return Promise.resolve({ data: [] })
      if (url === "/auth/me") {
        return Promise.resolve({
          data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
        })
      }
      return Promise.reject(new Error("not mocked"))
    })

    renderWithProviders(<ReportsPage />)

    await waitFor(() => {
      expect(screen.getByText(/no data for this period/i)).toBeInTheDocument()
    })
  })
})
