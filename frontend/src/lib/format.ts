export function formatCents(cents: number): string {
  const euros = cents / 100
  return new Intl.NumberFormat("de-DE", {
    style: "currency",
    currency: "EUR",
  }).format(euros)
}

export function centsToInput(cents: number): string {
  return (cents / 100).toFixed(2)
}

export function inputToCents(value: string): number {
  return Math.round(parseFloat(value) * 100)
}
