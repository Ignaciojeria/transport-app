import { NextRequest, NextResponse } from 'next/server'

export async function POST(request: NextRequest) {
  try {
    const body = await request.text()
    
    // Reenviar la request al backend
    const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL || 'https://einar-main-f0820bc.d2.zuplo.dev'
    
    const response = await fetch(`${backendUrl}/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: body
    })

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Error refreshing token' },
        { status: response.status }
      )
    }

    const data = await response.json()
    return NextResponse.json(data)

  } catch (error) {
    console.error('Refresh token error:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
}
