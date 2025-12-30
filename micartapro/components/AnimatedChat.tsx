'use client'

import { useState, useEffect } from 'react'
import { useLanguage } from '@/lib/useLanguage'

export function AnimatedChat() {
  const { t, isLoading } = useLanguage()
  const [currentMessageIndex, setCurrentMessageIndex] = useState(0)
  const [displayedText, setDisplayedText] = useState('')
  const [isTyping, setIsTyping] = useState(false)

  const messages = isLoading ? [] : t.chatExamples.messages

  useEffect(() => {
    if (messages.length === 0 || isLoading) return

    // Reset displayed text when message changes
    setDisplayedText('')
    setIsTyping(true)

    const currentMessage = messages[currentMessageIndex]
    let charIndex = 0

    // Typing animation
    const typingInterval = setInterval(() => {
      if (charIndex < currentMessage.length) {
        setDisplayedText(currentMessage.substring(0, charIndex + 1))
        charIndex++
      } else {
        setIsTyping(false)
        clearInterval(typingInterval)
        
        // Wait before showing next message
        setTimeout(() => {
          setCurrentMessageIndex((prev) => (prev + 1) % messages.length)
        }, 4000) // Show complete message for 4 seconds
      }
    }, 40) // Typing speed: 40ms per character

    return () => clearInterval(typingInterval)
  }, [currentMessageIndex, messages, isLoading])

  if (isLoading) {
    return null
  }

  return (
    <div className="w-full max-w-3xl mx-auto">
      {/* Large Search Input with Gradient Border */}
      <div className="relative">
        <div className="absolute inset-0 bg-gradient-to-r from-pink-500 via-purple-500 to-blue-500 rounded-2xl blur-sm opacity-30"></div>
        <div className="relative bg-white rounded-2xl border-2 border-transparent p-1">
          <div className="flex items-center px-6 py-5 bg-white rounded-xl">
            <div className="flex-1 min-w-0">
              <span className="text-gray-400 text-lg">
                {displayedText || t.chatExamples.placeholder}
              </span>
              {isTyping && (
                <span className="inline-block w-1 h-5 bg-blue-600 ml-1 animate-pulse">|</span>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Response indicator below input */}
      {!isTyping && displayedText.length > 0 && (
        <div className="mt-4 flex items-center justify-center animate-fade-in">
          <div className="flex items-center space-x-2 bg-green-50 border border-green-200 rounded-full px-4 py-2">
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <p className="text-sm text-green-700 font-medium">
              {t.chatExamples.response}
            </p>
          </div>
        </div>
      )}

      {/* Example suggestions */}
      <div className="mt-8 flex flex-wrap gap-3 justify-center">
        {messages.slice(0, 3).map((message, index) => (
          <button
            key={index}
            className="px-4 py-2 bg-gray-50 hover:bg-gray-100 border border-gray-200 rounded-full text-sm text-gray-600 transition-colors"
            disabled
          >
            {message.length > 40 ? message.substring(0, 40) + '...' : message}
          </button>
        ))}
      </div>
    </div>
  )
}

