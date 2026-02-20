/**
 * Token de demo para el video de Remotion.
 * Google bloquea OAuth en iframes. Pasas el token vía postMessage desde Remotion:
 *   { clickIniciarSesion: true, accessToken: "eyJ...", refreshToken?: "..." }
 */
export function getDemoAuthData(accessToken: string, refreshToken = '') {
  return {
    access_token: accessToken,
    token_type: 'Bearer' as const,
    expires_in: 3600,
    refresh_token: refreshToken,
    user: {
      id: 'demo-user-remotion',
      email: 'demo@micartapro.com',
      verified_email: true,
      name: 'Demo Remotion',
      given_name: 'Demo',
      family_name: 'Remotion',
      picture: '',
      locale: 'es',
    },
    timestamp: Date.now(),
    provider: 'supabase' as const,
  }
}
