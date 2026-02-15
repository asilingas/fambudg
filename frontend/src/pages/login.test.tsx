import { describe, it, expect, vi, beforeEach } from "vitest"
import { screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { renderWithProviders } from "@/test/test-utils"
import LoginPage from "./login"

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

beforeEach(() => {
  vi.clearAllMocks()
  localStorage.clear()
  mockedApi.get.mockRejectedValue(new Error("no token"))
})

describe("LoginPage", () => {
  it("renders email and password fields", () => {
    renderWithProviders(<LoginPage />)
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /sign in/i })).toBeInTheDocument()
  })

  it("shows error on failed login", async () => {
    const user = userEvent.setup()
    mockedApi.post.mockRejectedValueOnce({
      response: { data: { error: "invalid credentials" } },
    })

    renderWithProviders(<LoginPage />)

    await user.type(screen.getByLabelText(/email/i), "bad@example.com")
    await user.type(screen.getByLabelText(/password/i), "wrong")
    await user.click(screen.getByRole("button", { name: /sign in/i }))

    await waitFor(() => {
      expect(screen.getByRole("alert")).toHaveTextContent("invalid credentials")
    })
  })

  it("stores token and navigates on successful login", async () => {
    const user = userEvent.setup()
    mockedApi.post.mockResolvedValueOnce({
      data: {
        token: "jwt-token-123",
        user: { id: "1", email: "admin@test.com", name: "Admin", role: "admin" },
      },
    })

    renderWithProviders(<LoginPage />, {
      routerProps: { initialEntries: ["/login"] },
    })

    await user.type(screen.getByLabelText(/email/i), "admin@test.com")
    await user.type(screen.getByLabelText(/password/i), "password123")
    await user.click(screen.getByRole("button", { name: /sign in/i }))

    await waitFor(() => {
      expect(localStorage.getItem("token")).toBe("jwt-token-123")
    })
  })
})
