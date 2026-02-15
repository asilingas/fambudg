import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "@/test/test-utils"
import AccountsPage from "./accounts"

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

const accounts = [
  { id: "1", userId: "u1", name: "Checking", type: "checking", currency: "EUR", balance: 150000, createdAt: "2026-01-01" },
  { id: "2", userId: "u1", name: "Savings", type: "savings", currency: "EUR", balance: 500000, createdAt: "2026-01-01" },
]

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url === "/accounts") return Promise.resolve({ data: accounts })
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("AccountsPage", () => {
  it("renders account cards with balances", async () => {
    renderWithProviders(<AccountsPage />)

    await waitFor(() => {
      expect(screen.getByText("Checking")).toBeInTheDocument()
    })
    expect(screen.getByText("Savings")).toBeInTheDocument()
    expect(screen.getByText(/1\.500,00/)).toBeInTheDocument()
    expect(screen.getByText(/5\.000,00/)).toBeInTheDocument()
  })

  it("opens create dialog on add button click", async () => {
    const user = userEvent.setup()
    renderWithProviders(<AccountsPage />)

    await waitFor(() => {
      expect(screen.getByText("Checking")).toBeInTheDocument()
    })

    await user.click(screen.getByRole("button", { name: /add account/i }))
    expect(screen.getByText("New Account")).toBeInTheDocument()
  })

  it("shows empty state when no accounts", async () => {
    mockedApi.get.mockImplementation((url: string) => {
      if (url === "/accounts") return Promise.resolve({ data: [] })
      return Promise.reject(new Error("not mocked"))
    })

    renderWithProviders(<AccountsPage />)

    await waitFor(() => {
      expect(screen.getByText(/no accounts yet/i)).toBeInTheDocument()
    })
  })
})
