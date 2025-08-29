import { useState, useEffect } from 'react';
import { uploadQueue, type UploadProgress } from '../utils/offline-upload-queue';
import { CheckCircle, XCircle } from 'lucide-react';

export function UploadProgressBar() {
  const [uploads, setUploads] = useState<UploadProgress[]>([]);
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [isCollapsed, setIsCollapsed] = useState(false);

  useEffect(() => {
    // Subscribe to upload progress
    const unsubscribeProgress = uploadQueue.onProgress(setUploads);
    
    // Subscribe to network changes
    const unsubscribeNetwork = uploadQueue.onNetworkChange(setIsOnline);

    return () => {
      unsubscribeProgress();
      unsubscribeNetwork();
    };
  }, []);

  const stats = uploads.reduce(
    (acc, upload) => {
      acc[upload.status]++;
      return acc;
    },
    { pending: 0, uploading: 0, completed: 0, failed: 0 }
  );

  const hasUploads = uploads.length > 0;
  const hasActiveUploads = stats.pending > 0 || stats.uploading > 0;

  if (!hasUploads) {
    return null;
  }

  return (
    <div className="fixed top-4 right-4 z-50">
      <div className={`bg-white rounded-lg shadow-lg border transition-all duration-300 ${
        isCollapsed ? 'w-16 h-16' : 'w-80 min-h-16'
      }`}>
        {/* Header */}
        <div 
          className="flex items-center justify-between p-4 cursor-pointer"
          onClick={() => setIsCollapsed(!isCollapsed)}
        >
          <div className="flex items-center space-x-2">
            {isOnline ? (
              <div className="w-4 h-4 bg-green-500 rounded-full"></div>
            ) : (
              <div className="w-4 h-4 bg-red-500 rounded-full"></div>
            )}
            
            {!isCollapsed && (
              <>
                <div className="w-4 h-4 bg-blue-500 rounded-sm"></div>
                <span className="text-sm font-medium text-gray-700">
                  Subidas de Evidencia
                </span>
              </>
            )}
          </div>

          {!isCollapsed && (
            <div className="flex items-center space-x-1 text-xs text-gray-500">
              {stats.uploading > 0 && (
                <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded-full">
                  {stats.uploading} subiendo
                </span>
              )}
              {stats.pending > 0 && (
                <span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full">
                  {stats.pending} pendiente
                </span>
              )}
              {stats.failed > 0 && (
                <span className="bg-red-100 text-red-800 px-2 py-1 rounded-full">
                  {stats.failed} fallido
                </span>
              )}
              {stats.completed > 0 && (
                <span className="bg-green-100 text-green-800 px-2 py-1 rounded-full">
                  {stats.completed} completado
                </span>
              )}
            </div>
          )}

          {isCollapsed && hasActiveUploads && (
            <div className="absolute -top-1 -right-1 w-3 h-3 bg-blue-500 rounded-full animate-pulse"></div>
          )}
        </div>

        {/* Upload List */}
        {!isCollapsed && (
          <div className="border-t max-h-64 overflow-y-auto">
            {!isOnline && (
              <div className="p-3 bg-yellow-50 border-b">
                <div className="flex items-center space-x-2 text-sm text-yellow-800">
                  <div className="w-4 h-4 bg-red-500 rounded-full"></div>
                  <span>Sin conexión - Las subidas continuarán cuando vuelva la red</span>
                </div>
              </div>
            )}

            {uploads.map((upload) => (
              <div key={upload.itemId} className="p-3 border-b last:border-b-0">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center space-x-2">
                    {upload.status === 'completed' && (
                      <CheckCircle className="w-4 h-4 text-green-500" />
                    )}
                    {upload.status === 'uploading' && (
                      <div className="w-4 h-4 bg-blue-500 rounded-sm animate-pulse"></div>
                    )}
                    {upload.status === 'pending' && (
                      <div className="w-4 h-4 bg-yellow-500 rounded-full"></div>
                    )}
                    {upload.status === 'failed' && (
                      <XCircle className="w-4 h-4 text-red-500" />
                    )}
                    
                    <span className="text-sm font-medium text-gray-700">
                      Evidencia #{upload.itemId.slice(-4)}
                    </span>
                  </div>

                  <span className="text-xs text-gray-500 capitalize">
                    {upload.status === 'pending' ? 'Pendiente' :
                     upload.status === 'uploading' ? 'Subiendo...' :
                     upload.status === 'completed' ? 'Completado' :
                     upload.status === 'failed' ? 'Falló' : upload.status}
                  </span>
                </div>

                {/* Progress Bar */}
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className={`h-2 rounded-full transition-all duration-300 ${
                      upload.status === 'completed' ? 'bg-green-500' :
                      upload.status === 'uploading' ? 'bg-blue-500' :
                      upload.status === 'failed' ? 'bg-red-500' :
                      'bg-gray-300'
                    }`}
                    style={{ width: `${upload.progress}%` }}
                  ></div>
                </div>

                {upload.error && (
                  <div className="mt-2 text-xs text-red-600">
                    {upload.error}
                  </div>
                )}
              </div>
            ))}

            {/* Action Buttons */}
            <div className="p-3 border-t bg-gray-50 space-y-2">
              {/* Manual Process Button */}
              {hasActiveUploads && (
                <button
                  onClick={() => {
                    console.log('🔄 Manual queue processing triggered');
                    uploadQueue.processQueueManually();
                  }}
                  className="w-full text-xs bg-blue-500 text-white px-3 py-2 rounded hover:bg-blue-600 transition-colors"
                >
                  Procesar Subidas Ahora
                </button>
              )}
              
              {/* Clear Button */}
              {(stats.completed > 0 || stats.failed > 0) && (
                <button
                  onClick={() => uploadQueue.clearQueue()}
                  className="w-full text-xs text-gray-600 hover:text-gray-800 transition-colors"
                >
                  Limpiar completados y fallidos
                </button>
              )}
              
              {/* Debug Info */}
              <div className="text-xs text-gray-500">
                Network: {isOnline ? 'Online' : 'Offline'} | 
                Total: {uploads.length}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}