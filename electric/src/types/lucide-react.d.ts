declare module 'lucide-react' {
  import * as React from 'react'

  export interface IconProps extends React.SVGProps<SVGSVGElement> {
    size?: number | string
  }

  export const CheckCircle: React.FC<IconProps>
  export const XCircle: React.FC<IconProps>
  export const Play: React.FC<IconProps>
  export const Package: React.FC<IconProps>
  export const Phone: React.FC<IconProps>
  export const User: React.FC<IconProps>
  export const MapPin: React.FC<IconProps>
  export const Maximize2: React.FC<IconProps>
  export const Minimize2: React.FC<IconProps>
  export const Crosshair: React.FC<IconProps>
  
}


