import { format, parseISO, addDays, isAfter, differenceInDays } from 'date-fns'
import { ptBR } from 'date-fns/locale'

/**
 * Formata uma data ISO para o formato brasileiro
 */
export const formatDate = (dateString: string): string => {
  const date = parseISO(dateString)
  return format(date, "dd 'de' MMMM 'de' yyyy", { locale: ptBR })
}

/**
 * Formata uma data ISO para o formato curto (dd/MM/yyyy)
 */
export const formatDateShort = (dateString: string): string => {
  const date = parseISO(dateString)
  return format(date, 'dd/MM/yyyy', { locale: ptBR })
}

/**
 * Adiciona dias a uma data
 */
export const addDaysToDate = (dateString: string, days: number): Date => {
  const date = parseISO(dateString)
  return addDays(date, days)
}

/**
 * Verifica se uma data é futura
 */
export const isFutureDate = (dateString: string): boolean => {
  const date = parseISO(dateString)
  return isAfter(date, new Date())
}

/**
 * Calcula a diferença em dias entre duas datas
 */
export const getDaysDifference = (
  startDate: string,
  endDate: string
): number => {
  const start = parseISO(startDate)
  const end = parseISO(endDate)
  return differenceInDays(end, start)
}
