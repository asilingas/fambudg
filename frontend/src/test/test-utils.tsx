import type { ReactNode } from "react"
import { render, type RenderOptions } from "@testing-library/react"
import { MemoryRouter, type MemoryRouterProps } from "react-router-dom"
import { AuthProvider } from "@/context/auth-context"
import { LanguageProvider } from "@/context/language-context"

interface WrapperProps {
  children: ReactNode
}

interface CustomRenderOptions extends RenderOptions {
  routerProps?: MemoryRouterProps
}

export function renderWithProviders(
  ui: ReactNode,
  options: CustomRenderOptions = {},
) {
  const { routerProps, ...renderOptions } = options

  function Wrapper({ children }: WrapperProps) {
    return (
      <LanguageProvider>
        <MemoryRouter {...routerProps}>
          <AuthProvider>{children}</AuthProvider>
        </MemoryRouter>
      </LanguageProvider>
    )
  }

  return render(ui, { wrapper: Wrapper, ...renderOptions })
}
