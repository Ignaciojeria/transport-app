import { useShape } from '@electric-sql/react'

function Component() {
  // Una línea para sacar el token del fragment
  const token = new URLSearchParams(window.location.hash.slice(1)).get('access_token') || 
               new URLSearchParams(window.location.hash.slice(1)).get('token') || 
               "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleS04NWQ1MWE4MTYxZDcwMDI1IiwidHlwIjoiSldUIn0.eyJzdWIiOiI1ZjUzMWU3My0wN2FkLTU1MDItYWQwMy1lMWNhZGRiYmMxMWYiLCJzY29wZXMiOlsib3JkZXJzOnJlYWQiXSwiY29udGV4dCI6e30sInRlbmFudCI6IjVmNTMxZTczLTA3YWQtNTUwMi1hZDAzLWUxY2FkZGJiYzExZi1DaGlsZSIsImlzcyI6InRyYW5zcG9ydC1hcHAiLCJhdWQiOlsienVwbG8tZ2F0ZXdheSJdLCJleHAiOjE3NTU0NTUyNTMsIm5iZiI6MTc1NTQ1MTY1MywiaWF0IjoxNzU1NDUxNjUzfQ.M-ccduULYOJJCg-EbRBKwaWty9nfR9F7pWym58Qmxq1lWgLBll1fuilCMhvgp2PcbsFqtQOVh2aQKuQg7Pw3P0QYYyaS0uk-gwMIqY0zV_7Vzqi3pEUH57guwWNEnbZfxpfbc01hsyVBLGMV1GN-aOtnEpkjsvHfK_avzsdPUFtPjQwl9YiziQoioil3076mWMsZoCaseM5GA9vCJMmq9rHOK2lOeOyWUHQnMoJKkH60CQgw68I-4qjQIb9hEeqzqBxKY6NqH6ilBN-6M5YDZDSRZbhcm4qAWFjiABkUKQ2WZGpHKjfLZsfKuMfZ3ooI9U6vt1vunHoOsSpx_11PWQ";

  const { data } = useShape({
    headers: {
      "X-Access-Token": `Bearer ${token}`
    },
    url: `https://einar-main-f0820bc.d2.zuplo.dev/electric/v1/shape?table=orders`
  })

  return (
    <pre>{JSON.stringify(data, (_key, value) => 
      typeof value === 'bigint' ? value.toString() : value, 2)}</pre>
  )
}

export default Component