import { useShape } from '@electric-sql/react'
import type { ElectricAccountData } from '../db/collections/create-accounts-collection'

export const useAccountByEmail = (token: string | null, email: string | null) => {
  // Construir la URL siempre, pero usar parámetros vacíos si no hay token/email
  const url = (() => {
    if (!token || !email) {
      // URL que no devuelve resultados cuando no hay token/email
      return `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id&where=id='__INVALID__'`
    }
    const base = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id`
    return `${base}&where=email='${email}'`
  })()

  // Siempre llamar useShape (hooks deben llamarse incondicionalmente)
  const shapeResult = useShape({
    url,
    headers: token ? {
      'X-Access-Token': `Bearer ${token}`,
    } : {},
  })

  const data = shapeResult?.data as ElectricAccountData[] | undefined
  const isLoading = shapeResult?.isLoading ?? false
  const error = shapeResult?.error

  // Si no hay token/email, considerar como no existe cuenta
  const accountExists = token && email && data && data.length > 0
  const accountId = accountExists ? data[0].id : null

  return {
    accountExists,
    accountId,
    isLoading: token && email ? isLoading : false,
    error: token && email ? error : null,
    data
  }
}
