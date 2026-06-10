// Formatting helpers shared across pages (currency, dates).
export function useFormat() {
  const money = (value: number | null | undefined, currency = '₽') => {
    const n = Number(value || 0)
    return `${n.toLocaleString('ru-RU', { minimumFractionDigits: 0, maximumFractionDigits: 2 })} ${currency}`
  }

  const date = (value: string | null | undefined) => {
    if (!value) return '—'
    const d = new Date(value)
    if (Number.isNaN(d.getTime())) return String(value).slice(0, 10)
    return d.toLocaleDateString('ru-RU')
  }

  const dateTime = (value: string | null | undefined) => {
    if (!value) return '—'
    const d = new Date(value)
    if (Number.isNaN(d.getTime())) return String(value)
    return d.toLocaleString('ru-RU', { dateStyle: 'short', timeStyle: 'short' })
  }

  return { money, date, dateTime }
}
