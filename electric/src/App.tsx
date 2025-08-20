import { useShape } from '@electric-sql/react'

function Component() {
  // Una línea para sacar el token del fragment
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || 
               "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS04NWQ1MWE4MTYxZDcwMDI1IiwidHlwIjoiSldUIn0.eyJzdWIiOiI1ZjUzMWU3My0wN2FkLTU1MDItYWQwMy1lMWNhZGRiYmMxMWYiLCJzY29wZXMiOlsib3JkZXJzOnJlYWQiXSwiY29udGV4dCI6e30sInRlbmFudCI6IjVmNTMxZTczLTA3YWQtNTUwMi1hZDAzLWUxY2FkZGJiYzExZi1DaGlsZSIsImlzcyI6InRyYW5zcG9ydC1hcHAiLCJhdWQiOlsienVwbG8tZ2F0ZXdheSJdLCJleHAiOjE3NTU2NTkxODIsIm5iZiI6MTc1NTY1NTU4MiwiaWF0IjoxNzU1NjU1NTgyfQ.KIAcQcPoPmYyXznP-Wnr_zb1VKtOzxuN4fG74wO5iO-Ckv2ii8L95PMXilMD7kw-NSkIUy6nPkhF-47B6WwzIZxuFyGSqqaR25e3SaXktmWC4qzAJ-R-g-eNwFNaCGvJlFOhDPLSq1RHM3K9WLVfLUhb3H0Yk22N3YwtdBlwzUw2U699_K540kNGcaACLN1F0KaVM3u4S4-1J211XB5SbDj2ong_PgWAK2kaHnAo4_dXlC_m2n0Wt38DKgnOVBEg9-Ekp9akhL57A8JonnExZ79PkvPELQzSBATvD4sgRVEx-L-sOF1gHSK8VIu9oUY0AdQ8mRB9b_spo_G-WMm9eQ";

  // Extraer reference_id de los query parameters
  const referenceId = new URLSearchParams(window.location.search).get('reference_id');

  const { data } = useShape({
    headers: {
      "X-Access-Token": `Bearer ${token}`
    },
    url: `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=routes&columns=id,raw${referenceId ? `&where=reference_id='${referenceId}'` : ''}`
  })

  return (
    <pre>{JSON.stringify(data, (_key, value) => 
      typeof value === 'bigint' ? value.toString() : value, 2)}</pre>
  )
}

export default Component