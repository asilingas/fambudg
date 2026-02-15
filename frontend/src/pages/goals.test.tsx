import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "@/test/test-utils"
import GoalsPage from "./goals"

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

const goals = [
  {
    id: "g1",
    name: "Vacation Fund",
    targetAmount: 200000,
    currentAmount: 80000,
    targetDate: "2026-06-01T00:00:00Z",
    priority: 1,
    status: "active",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-02-01T00:00:00Z",
  },
  {
    id: "g2",
    name: "Emergency Fund",
    targetAmount: 100000,
    currentAmount: 100000,
    priority: 2,
    status: "completed",
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-02-01T00:00:00Z",
  },
]

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url === "/saving-goals") return Promise.resolve({ data: goals })
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("GoalsPage", () => {
  it("renders saving goals with progress", async () => {
    renderWithProviders(<GoalsPage />)

    await waitFor(() => {
      expect(screen.getByText("Vacation Fund")).toBeInTheDocument()
    })
    expect(screen.getByText("Emergency Fund")).toBeInTheDocument()
    expect(screen.getByText(/40% reached/)).toBeInTheDocument()
  })

  it("shows completed badge for completed goals", async () => {
    renderWithProviders(<GoalsPage />)

    await waitFor(() => {
      expect(screen.getByText("Completed")).toBeInTheDocument()
    })
  })

  it("shows contribute button for active goals (admin)", async () => {
    renderWithProviders(<GoalsPage />)

    await waitFor(() => {
      expect(screen.getByText("Vacation Fund")).toBeInTheDocument()
    })
    expect(screen.getByRole("button", { name: /contribute/i })).toBeInTheDocument()
  })

  it("opens contribute dialog on click", async () => {
    const user = userEvent.setup()
    renderWithProviders(<GoalsPage />)

    await waitFor(() => {
      expect(screen.getByText("Vacation Fund")).toBeInTheDocument()
    })

    await user.click(screen.getByRole("button", { name: /contribute/i }))
    expect(screen.getByText(/Contribute to Vacation Fund/)).toBeInTheDocument()
  })

  it("shows empty state when no goals", async () => {
    mockedApi.get.mockImplementation((url: string) => {
      if (url === "/saving-goals") return Promise.resolve({ data: [] })
      if (url === "/auth/me") {
        return Promise.resolve({
          data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
        })
      }
      return Promise.reject(new Error("not mocked"))
    })

    renderWithProviders(<GoalsPage />)

    await waitFor(() => {
      expect(screen.getByText(/no saving goals yet/i)).toBeInTheDocument()
    })
  })
})
