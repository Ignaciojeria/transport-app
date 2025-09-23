import NextAuth from 'next-auth'
import GoogleProvider from 'next-auth/providers/google'
import type { NextAuthOptions } from 'next-auth'

export const authOptions: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
      authorization: {
        params: {
          prompt: "consent",
          access_type: "offline",
          response_type: "code",
          hd: process.env.GOOGLE_HD_DOMAIN // Restricir a dominios corporativos específicos
        }
      }
    })
  ],
  pages: {
    signIn: '/',
    error: '/auth/error',
  },
  callbacks: {
    async jwt({ token, user, account }) {
      if (account && user) {
        token.accessToken = account.access_token
        token.role = user.role || 'user'
      }
      return token
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken
      session.user.role = token.role
      return session
    },
    async signIn({ user, account, profile }) {
      // Validar dominio de email para Google
      if (account?.provider === 'google') {
        const email = user.email || ''
        const allowedDomains = process.env.ALLOWED_DOMAINS?.split(',') || []
        
        if (allowedDomains.length > 0) {
          const domain = email.split('@')[1]
          if (!allowedDomains.includes(domain)) {
            return false
          }
        }
      }
      return true
    }
  },
  session: {
    strategy: 'jwt',
    maxAge: 30 * 24 * 60 * 60, // 30 días
  },
  secret: process.env.NEXTAUTH_SECRET,
}

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST }
