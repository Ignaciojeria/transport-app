import { useShape } from '@electric-sql/react'
import { useParams } from '@tanstack/react-router'

// Componente para la ruta raíz
export function HomeComponent() {
  // Una línea para sacar el token del fragment
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || 
               "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS04NWQ1MWE4MTYxZDcwMDI1IiwidHlwIjoiSldUIn0.eyJzdWIiOiI1ZjUzMWU3My0wN2FkLTU1MDItYWQwMy1lMWNhZGRiYmMxMWYiLCJzY29wZXMiOlsib3JkZXJzOnJlYWQiXSwiY29udGV4dCI6e30sInRlbmFudCI6IjVmNTMxZTczLTA3YWQtNTUwMi1hZDAzLWUxY2FkZGJiYzExZiIsImlzcyI6InRyYW5zcG9ydC1hcHAiLCJhdWQiOlsienVwbG8tZ2F0ZXdheSJdLCJleHAiOjE3NTU2NTkxODIsIm5iZiI6MTc1NTY1NTU4MiwiaWF0IjoxNzU1NjU1NTgyfQ.KIAcQcPoPmYyXznP-Wnr_zb1VKtOzxuN4fG74wO5iO-Ckv2ii8L95PMXilMD7kw-NSkIUy6nPkhF-47B6WwzIZxuFyGSqqaR25e3SaXktmWC4qzAJ-R-g-eNwFNaCGvJlFOhDPLSq1RHM3K9WLVfLUhb3H0Yk22N3YwtdBlwzUw2U699_K540kNGcaACLN1F0KaVM3u4S4-1J211XB5SbDj2ong_PgWAK2kaHnAo4_dXlC_m2n0Wt38DKgnOVBEg9-Ekp9akhL57A8JonnExZ79PkvPELQzSBATvD4sgRVEx-L-sOF1gHSK8VIu9oUY0AdQ8mRB9b_spo_G-WMm9eQ";

  const { data } = useShape({
    headers: {
      "X-Access-Token": `Bearer ${token}`
    },
    url: `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=routes&columns=id,raw`
  })

  return (
    <div>
      <h1>Rutas Disponibles</h1>
      <p>Accede a una ruta específica usando: /driver/routes/&#123;routeId&#125;#access_token=&#123;token&#125;</p>
      <pre>{JSON.stringify(data, (_key, value) => 
        typeof value === 'bigint' ? value.toString() : value, 2)}</pre>
    </div>
  )
}

// Componente para rutas específicas del driver
export function RouteComponent() {
  // Obtener el routeId de los parámetros de la ruta usando TanStack Router
  const { routeId } = useParams({ from: '/driver/routes/$routeId' })
  
  // Una línea para sacar el token del fragment
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || 
               "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS04NWQ1MWE4MTYxZDcwMDI1IiwidHlwIjoiSldUIn0.eyJzdWIiOiI1ZjUzMWU3My0wN2FkLTU1MDItYWQwMy1lMWNhZGRiYmMxMWYiLCJzY29wZXMiOlsib3JkZXJzOnJlYWQiXSwiY29udGV4dCI6e30sInRlbmFudCI6IjVmNTMxZTczLTA3YWQtNTUwMi1hZDAzLWUxY2FkZGJiYzExZiIsImlzcyI6InRyYW5zcG9ydC1hcHAiLCJhdWQiOlsienVwbG8tZ2F0ZXdheSJdLCJleHAiOjE3NTU2NTkxODIsIm5iZiI6MTc1NTY1NTU4MiwiaWF0IjoxNzU1NjU1NTgyfQ.KIAcQcPoPmYyXznP-Wnr_zb1VKtOzxuN4fG74wO5iO-Ckv2ii8L95PMXilMD7kw-NSkIUy6nPkhF-47B6WwzIZxuFyGSqqaR25e3SaXktmWC4qzAJ-R-g-eNwFNaCGvJlFOhDPLSq1RHM3K9WLVfLUhb3H0Yk22N3YwtdBlwzUw2U699_K540kNGcaACLN1F0KaVM3u4S4-1J211XB5SbDj2ong_PgWAK2kaHnAo4_dXlC_m2n0Wt38DKgnOVBEg9-Ekp9akhL57A8JonnExZ79PkvPELQzSBATvD4sgRVEx-L-sOF1gHSK8VIu9oUY0AdQ8mRB9b_spo_G-WMm9eQ";

  const { data } = useShape({
    headers: {
      "X-Access-Token": `Bearer ${token}`
    },
    url: `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=routes&columns=id,raw${routeId ? `&where=reference_id='${routeId}'` : ''}`
  })

  return (
    <div>
      <h2>Driver - Ruta ID: {routeId}</h2>
      <p>Token: {token ? '✅ Presente' : '❌ No encontrado'}</p>
      <pre>{JSON.stringify(data, (_key, value) => 
        typeof value === 'bigint' ? value.toString() : value, 2)}</pre>
    </div>
  )
}