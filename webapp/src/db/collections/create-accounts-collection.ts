import { useShape } from '@electric-sql/react'

// Tipo para la estructura que devuelve Electric
type ElectricAccountData = {
  id: string
}

export type { ElectricAccountData }

// Factory para crear la colección inyectando el token
export const createAccountsCollection = (token: string, email?: string) => {
  const url = (() => {
    const base = `https://einar-main-f0820bc.d2.zuplo.dev/electric-me/v1/shape?table=accounts&columns=id`
    return email ? `${base}&where=email='${email}'` : base
  })()
  
  return useShape({
    url,
    headers: {
      'X-Access-Token': `Bearer ${token}`,
    },
  })
}
