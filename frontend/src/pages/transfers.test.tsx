import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "@/test/test-utils"
import TransfersPage from "./transfers"

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

const mockedApi = vi.mocked(api, true)

const accounts = [
  { id: "a1", userId: "u1", name: "Checking", type: "checking", currency: "EUR", balance: 150000, createdAt: "2026-01-01" },
  { id: "a2", userId: "u1", name: "Savings", type: "savings", currency: "EUR", balance: 500000, createdAt: "2026-01-01" },
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
  mockedApi.post.mockResolvedValue({ data: {} })
})

describe("TransfersPage", () => {
  it("renders transfer form with account selects", async () => {
    renderWithProviders(<TransfersPage />)

    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Transfer" })).toBeInTheDocument()
    })
    expect(screen.getByLabelText(/amount/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/date/i)).toBeInTheDocument()
  })

  it("renders account select dropdowns", async () => {
    renderWithProviders(<TransfersPage />)

    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Transfer" })).toBeInTheDocument()
    })

    // Two combobox selects (From and To)
    const selects = screen.getAllByRole("combobox")
    expect(selects.length).toBe(2)
  })

  it("shows description and date fields", async () => {
    renderWithProviders(<TransfersPage />)

    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Transfer" })).toBeInTheDocument()
    })

    expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/date/i)).toBeInTheDocument()
  })
})
