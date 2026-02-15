import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "@/test/test-utils"
import UsersPage from "./users"

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

const users = [
  { id: "u1", email: "admin@test.com", name: "Admin", role: "admin", createdAt: "2026-01-01T00:00:00Z", updatedAt: "2026-01-01T00:00:00Z" },
  { id: "u2", email: "member@test.com", name: "Member", role: "member", createdAt: "2026-01-15T00:00:00Z", updatedAt: "2026-01-15T00:00:00Z" },
  { id: "u3", email: "child@test.com", name: "Kid", role: "child", createdAt: "2026-02-01T00:00:00Z", updatedAt: "2026-02-01T00:00:00Z" },
]

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.setItem("token", "test-token")
  mockedApi.get.mockImplementation((url: string) => {
    if (url === "/users") return Promise.resolve({ data: users })
    if (url === "/auth/me") {
      return Promise.resolve({
        data: { id: "u1", email: "admin@test.com", name: "Admin", role: "admin" },
      })
    }
    return Promise.reject(new Error("not mocked"))
  })
})

describe("UsersPage", () => {
  it("renders users table with roles", async () => {
    renderWithProviders(<UsersPage />)

    await waitFor(() => {
      expect(screen.getByText("Admin")).toBeInTheDocument()
    })
    expect(screen.getByText("Member")).toBeInTheDocument()
    expect(screen.getByText("Kid")).toBeInTheDocument()
    expect(screen.getByText("admin")).toBeInTheDocument()
    expect(screen.getByText("member")).toBeInTheDocument()
    expect(screen.getByText("child")).toBeInTheDocument()
  })

  it("shows add user button", async () => {
    renderWithProviders(<UsersPage />)

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /add user/i })).toBeInTheDocument()
    })
  })

  it("opens create dialog on add button click", async () => {
    const user = userEvent.setup()
    renderWithProviders(<UsersPage />)

    await waitFor(() => {
      expect(screen.getByText("Admin")).toBeInTheDocument()
    })

    await user.click(screen.getByRole("button", { name: /add user/i }))
    expect(screen.getByText("Create User")).toBeInTheDocument()
  })

  it("shows edit and delete buttons per user row", async () => {
    renderWithProviders(<UsersPage />)

    await waitFor(() => {
      expect(screen.getByText("Admin")).toBeInTheDocument()
    })
    // 3 users × 2 buttons each = 6 icon buttons (edit + delete)
    const buttons = screen.getAllByRole("button", { name: "" })
    // Filter to only the edit/delete icon buttons in the table
    expect(buttons.length).toBeGreaterThanOrEqual(6)
  })
})
