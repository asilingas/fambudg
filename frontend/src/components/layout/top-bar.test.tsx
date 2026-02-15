import { render, screen, fireEvent } from "@testing-library/react"
import { LanguageProvider } from "@/context/language-context"

const mockLogout = vi.fn()
const mockToggleTheme = vi.fn()
let mockTheme = "light"
let mockUser: { name: string; role: string } | null = { name: "Test User", role: "admin" }

vi.mock("@/context/auth-context", () => ({
  useAuth: () => ({ user: mockUser, logout: mockLogout }),
}))

vi.mock("@/hooks/use-theme", () => ({
  useTheme: () => ({ theme: mockTheme, toggleTheme: mockToggleTheme }),
}))

import { TopBar } from "./top-bar"

function renderTopBar() {
  return render(
    <LanguageProvider>
      <TopBar />
    </LanguageProvider>,
  )
}

describe("TopBar", () => {
  beforeEach(() => {
    mockTheme = "light"
    mockUser = { name: "Test User", role: "admin" }
    vi.clearAllMocks()
    localStorage.clear()
  })

  it("renders user name and role", () => {
    renderTopBar()
    expect(screen.getByText("Test User")).toBeInTheDocument()
    expect(screen.getByText("admin")).toBeInTheDocument()
  })

  it("renders theme toggle button", () => {
    renderTopBar()
    expect(screen.getByLabelText("Toggle theme")).toBeInTheDocument()
  })

  it("renders language toggle button", () => {
    renderTopBar()
    expect(screen.getByLabelText("Toggle language")).toBeInTheDocument()
    expect(screen.getByText("LT")).toBeInTheDocument()
  })

  it("calls toggleTheme on theme button click", () => {
    renderTopBar()
    fireEvent.click(screen.getByLabelText("Toggle theme"))
    expect(mockToggleTheme).toHaveBeenCalled()
  })

  it("toggles language on language button click", () => {
    renderTopBar()
    expect(screen.getByText("LT")).toBeInTheDocument()
    fireEvent.click(screen.getByLabelText("Toggle language"))
    expect(screen.getByText("EN")).toBeInTheDocument()
  })

  it("calls logout on logout button click", () => {
    renderTopBar()
    fireEvent.click(screen.getByLabelText("Logout"))
    expect(mockLogout).toHaveBeenCalled()
  })

  it("returns null when no user", () => {
    mockUser = null
    const { container } = renderTopBar()
    expect(container.firstChild).toBeNull()
  })
})
