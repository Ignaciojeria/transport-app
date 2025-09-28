import { ReactNode } from 'react'

interface CardProps {
  children: ReactNode
  className?: string
}

export function Card({ children, className = '' }: CardProps) {
  return (
    <div className={`bg-white rounded-xl shadow-lg border border-gray-200 p-10 ${className}`}>
      {children}
    </div>
  )
}

export function CardHeader({ children, className = '' }: CardProps) {
  return (
    <div className={`mb-8 ${className}`}>
      {children}
    </div>
  )
}

export function CardTitle({ children, className = '' }: CardProps) {
  return (
    <h2 className={`text-3xl font-bold text-gray-900 ${className}`}>
      {children}
    </h2>
  )
}

export function CardDescription({ children, className = '' }: CardProps) {
  return (
    <p className={`text-gray-600 mt-3 text-lg ${className}`}>
      {children}
    </p>
  )
}

export function CardContent({ children, className = '' }: CardProps) {
  return (
    <div className={className}>
      {children}
    </div>
  )
}

export function CardFooter({ children, className = '' }: CardProps) {
  return (
    <div className={`mt-8 ${className}`}>
      {children}
    </div>
  )
}
