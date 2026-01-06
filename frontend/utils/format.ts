/**
 * Formata um valor numérico como moeda brasileira
 */
export const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(value)
}

/**
 * Formata um número com separadores de milhar
 */
export const formatNumber = (value: number): string => {
  return new Intl.NumberFormat('pt-BR').format(value)
}

/**
 * Formata um número como porcentagem
 */
export const formatPercentage = (value: number, decimals: number = 2): string => {
  return `${value.toFixed(decimals)}%`
}
