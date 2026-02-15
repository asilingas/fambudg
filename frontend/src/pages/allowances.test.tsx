import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import { renderWithProviders } from "@/test/test-utils"
import AllowancesPage from "./allowances"

vi.mock("@/lib/api", () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  },
}))

import api from "@/lib/api"

const mockedApi = vi.mocked(api)

const allowances = [
  {
    id: "al1",
    userId: "u3",
    amount: 5000,
    spent: 2000,
    remaining: 3000,
    periodStart: "2026-02-01T00:00:00Z",
    createdAt: "2026-02-01T00:00:00Z",
    updatedAt: "2026-02-10T00:00:00Z",
  },
]

const users = [
  { id: "u1", email: "admin@test.com", name: "Admin", role: "admin", createdAt: "2026-01-01", updatedAt: "2026-01-01" },
  { id: "u3", email: "child@test.com", name: "Kid", role: "child", createdAt: "2026-01-01", updatedAt: "2026-01-01" },
]

describe("AllowancesPage", () => {
  describe("as admin", () => {
    beforeEach(() => {
      vi.clearAllMocks()
      localStorage.setItem("token", "test-token")
      mockedApi.get.mockImplementation((url: string) => {
        if (url === "/allowances") return Promise.resolve({ data: allowances })
        if (url === "/users") return Promise.resolve({ data: users })
        if (url === "/auth/me") {
          return Promise.resolve({
            data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
          })
        }
        return Promise.reject(new Error("not mocked"))
      })
    })

    it("renders allowance cards with spending info", async () => {
      renderWithProviders(<AllowancesPage />)

      await waitFor(() => {
        expect(screen.getByText("Kid")).toBeInTheDocument()
      })
      expect(screen.getByText(/Spent:/)).toBeInTheDocument()
      expect(screen.getByText(/Remaining:/)).toBeInTheDocument()
    })

    it("shows set allowance button for admin", async () => {
      renderWithProviders(<AllowancesPage />)

      await waitFor(() => {
        expect(screen.getByRole("button", { name: /set allowance/i })).toBeInTheDocument()
      })
    })

    it("shows edit button on allowance card", async () => {
      renderWithProviders(<AllowancesPage />)

      await waitFor(() => {
        expect(screen.getByText("Kid")).toBeInTheDocument()
      })
      // The edit icon button
      const buttons = screen.getAllByRole("button")
      expect(buttons.length).toBeGreaterThanOrEqual(2)
    })
  })

  describe("as child", () => {
    beforeEach(() => {
      vi.clearAllMocks()
      localStorage.setItem("token", "test-token")
      mockedApi.get.mockImplementation((url: string) => {
        if (url === "/allowances") return Promise.resolve({ data: allowances })
        if (url === "/auth/me") {
          return Promise.resolve({
            data: { id: "u3", email: "child@test.com", name: "Kid", role: "child" },
          })
        }
        return Promise.reject(new Error("not mocked"))
      })
    })

    it("shows 'My Allowance' for child role", async () => {
      renderWithProviders(<AllowancesPage />)

      await waitFor(() => {
        expect(screen.getByText("My Allowance")).toBeInTheDocument()
      })
    })

    it("does not show set allowance button for child", async () => {
      renderWithProviders(<AllowancesPage />)

      await waitFor(() => {
        expect(screen.getByText("My Allowance")).toBeInTheDocument()
      })
      expect(screen.queryByRole("button", { name: /set allowance/i })).not.toBeInTheDocument()
    })
  })
})
