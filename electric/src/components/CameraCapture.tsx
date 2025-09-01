import { useRef, useState, useEffect } from 'react'
import Webcam from 'react-webcam'

interface CameraCaptureProps {
  onCapture: (photoDataUrl: string) => void
  title?: string
  buttonText?: string
  className?: string
}

export function CameraCapture({ 
  onCapture, 
  title = "Capturar foto",
  buttonText = "Activar cámara",
  className = ""
}: CameraCaptureProps) {
  const [usingCamera, setUsingCamera] = useState(false)
  const [cameraError, setCameraError] = useState<string | null>(null)
  const [photoDataUrl, setPhotoDataUrl] = useState<string | null>(null)
  const [flashActive, setFlashActive] = useState(false)
  const webcamRef = useRef<any>(null)
  const cameraContainerRef = useRef<HTMLDivElement | null>(null)

  const stopWebcam = () => {
    try {
      const stream: MediaStream | undefined = (webcamRef.current as any)?.stream
      if (stream) {
        stream.getTracks().forEach((t) => t.stop())
      }
    } catch {}
    setUsingCamera(false)
  }

  const captureFromWebcam = () => {
    try {
      const imgSrc = webcamRef.current?.getScreenshot?.()
      if (imgSrc) {
        setPhotoDataUrl(imgSrc)
        try { (navigator as any)?.vibrate?.(60) } catch {}
        setFlashActive(true)
        setTimeout(() => setFlashActive(false), 140)
        stopWebcam()
        // Llamar automáticamente a onCapture cuando se captura la foto
        onCapture(imgSrc)
      }
    } catch (e) {
      setCameraError('No se pudo capturar la imagen.')
    }
  }

  const handleActivateCamera = () => {
    try { (document.activeElement as any)?.blur?.() } catch {}
    setCameraError(null)
    setUsingCamera(true)
  }

  const handleRetakePhoto = () => {
    setPhotoDataUrl(null)
    setUsingCamera(true)
    setCameraError(null)
  }



  useEffect(() => {
    // Cuando se activa la cámara, desplazar el modal para centrar la vista de cámara
    if (usingCamera) {
      setTimeout(() => {
        try { cameraContainerRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' }) } catch {}
      }, 80)
    }
  }, [usingCamera])

  // Limpiar cámara al desmontar
  useEffect(() => {
    return () => {
      stopWebcam()
    }
  }, [])

  return (
    <>
      <div className={`space-y-3 ${className}`}>
        <div>
          <label className="block text-xs text-gray-600 mb-1">{title}</label>
          
          {/* Cámara con react-webcam: activar bajo demanda */}
          <div className="mb-2">
            {!usingCamera ? (
              <div>
                <button
                  type="button"
                  onClick={handleActivateCamera}
                  className="px-3 py-2 text-sm rounded-md border bg-white hover:bg-gray-50 transition-colors"
                >
                  {buttonText}
                </button>
                {cameraError && <p className="text-xs text-red-600 mt-1">{cameraError}</p>}
              </div>
            ) : (
              <div>
                <div
                  className="relative w-full h-[60vh] sm:h-96 rounded-md overflow-hidden border bg-black cursor-pointer select-none"
                  ref={cameraContainerRef}
                  onClick={captureFromWebcam}
                  title="Toca para capturar"
                >
                  <Webcam
                    ref={webcamRef}
                    audio={false}
                    screenshotFormat="image/jpeg"
                    className="w-full h-full object-cover"
                    videoConstraints={{
                      facingMode: { ideal: 'environment' },
                      width: { ideal: 1280 },
                      height: { ideal: 720 },
                    }}
                    onUserMediaError={() => setCameraError('No se pudo acceder a la cámara. Revisa permisos.')}
                  />
                  <div className="absolute bottom-0 left-0 right-0 bg-black/40 text-white text-xs text-center py-1">
                    Toca para capturar
                  </div>
                </div>
                <div className="flex items-center gap-2 mt-2">
                  <button 
                    type="button" 
                    onClick={stopWebcam} 
                    className="px-3 py-2 text-sm rounded-md border bg-white hover:bg-gray-50 transition-colors"
                  >
                    Cerrar cámara
                  </button>
                </div>
                {cameraError && <p className="text-xs text-red-600 mt-1">{cameraError}</p>}
              </div>
            )}
          </div>
          
          {/* Vista previa de la foto capturada */}
          {photoDataUrl && (
            <div className="mt-2 flex items-center gap-3">
              <img 
                src={photoDataUrl} 
                alt="Foto capturada" 
                className="w-24 h-24 object-cover rounded-md border" 
              />
              <button 
                type="button" 
                onClick={handleRetakePhoto} 
                className="px-3 py-2 text-sm rounded-md border bg-white hover:bg-gray-50 transition-colors"
              >
                Cambiar foto
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Efecto de flash */}
      {flashActive && (
        <div className="fixed inset-0 z-[100000] pointer-events-none bg-white opacity-70"></div>
      )}
    </>
  )
}
