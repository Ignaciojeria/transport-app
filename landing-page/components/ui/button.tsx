import * as React from "react"
import { cn } from "@/lib/utils"

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "default" | "destructive" | "outline" | "secondary" | "ghost" | "link"
  size?: "default" | "sm" | "lg" | "icon"
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "default", size = "default", ...props }, ref) => {
    return (
      <button
        className={cn(
          "magic-button",
          {
            "magic-button-primary": variant === "default",
            "magic-button-secondary": variant === "secondary",
            "magic-button-outline": variant === "outline",
            "magic-button-ghost": variant === "ghost",
            "magic-button-link": variant === "link",
            "bg-destructive text-destructive-foreground hover:bg-destructive/90": variant === "destructive",
            "magic-button-sm": size === "sm",
            "magic-button-lg": size === "lg",
            "magic-button-icon": size === "icon",
            "h-10 px-4 py-2": size === "default",
          },
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button }
