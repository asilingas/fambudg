import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "@/test/test-utils"
import SearchPage from "./search"

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

const categories = [
  { id: "c1", name: "Groceries", type: "expense", sortOrder: 0 },
]

const accounts = [
  { id: "a1", userId: "u1", name: "Checking", type: "checking", currency: "EUR", balance: 100000, createdAt: "2026-01-01" },
]

const searchResults = {
  transactions: [
    {
      id: "t1",
      userId: "u1",
      accountId: "a1",
      categoryId: "c1",
      amount: -4599,
      type: "expense",
      description: "Weekly groceries",
      date: "2026-02-10",
      isShared: false,
      isRecurring: false,
      tags: ["food"],
      createdAt: "2026-02-10T10:00:00Z",
      updatedAt: "2026-02-10T10:00:00Z",
    },
  ],
  totalCount: 1,
}

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url === "/categories") return Promise.resolve({ data: categories })
    if (url === "/accounts") return Promise.resolve({ data: accounts })
    if (url.startsWith("/search")) return Promise.resolve({ data: searchResults })
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("SearchPage", () => {
  it("renders search form with filters", async () => {
    renderWithProviders(<SearchPage />)

    await waitFor(() => {
      expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    })
    expect(screen.getByLabelText(/start date/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/end date/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/min amount/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/max amount/i)).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /search/i })).toBeInTheDocument()
  })

  it("displays search results after clicking Search", async () => {
    const user = userEvent.setup()
    renderWithProviders(<SearchPage />)

    await waitFor(() => {
      expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    })

    await user.type(screen.getByLabelText(/description/i), "groceries")
    await user.click(screen.getByRole("button", { name: /search/i }))

    await waitFor(() => {
      expect(screen.getByText("Weekly groceries")).toBeInTheDocument()
    })
    expect(screen.getByText("1 result found")).toBeInTheDocument()
    expect(screen.getByText("food")).toBeInTheDocument()
  })

  it("shows no results message when empty", async () => {
    mockedApi.get.mockImplementation((url: string) => {
      if (url === "/categories") return Promise.resolve({ data: categories })
      if (url === "/accounts") return Promise.resolve({ data: accounts })
      if (url.startsWith("/search")) return Promise.resolve({ data: { transactions: [], totalCount: 0 } })
      if (url === "/auth/me") {
        return Promise.resolve({
          data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
        })
      }
      return Promise.reject(new Error("not mocked"))
    })

    const user = userEvent.setup()
    renderWithProviders(<SearchPage />)

    await waitFor(() => {
      expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    })

    await user.click(screen.getByRole("button", { name: /search/i }))

    await waitFor(() => {
      expect(screen.getByText(/no transactions found/i)).toBeInTheDocument()
    })
  })

  it("clears results on Clear button click", async () => {
    const user = userEvent.setup()
    renderWithProviders(<SearchPage />)

    await waitFor(() => {
      expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    })

    await user.type(screen.getByLabelText(/description/i), "groceries")
    await user.click(screen.getByRole("button", { name: /search/i }))

    await waitFor(() => {
      expect(screen.getByText("Weekly groceries")).toBeInTheDocument()
    })

    await user.click(screen.getByRole("button", { name: /clear/i }))

    expect(screen.queryByText("Weekly groceries")).not.toBeInTheDocument()
    expect(screen.queryByText(/result/)).not.toBeInTheDocument()
  })
})
