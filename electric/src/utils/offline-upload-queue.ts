/**
 * Offline upload queue system for handling photo uploads with retries
 */

export interface UploadQueueItem {
  id: string;
  routeId: string;
  visitIndex: number;
  orderIndex: number;
  unitIndex: number;
  photoDataUrl: string;
  uploadUrl: string;
  type: 'delivery' | 'non-delivery';
  timestamp: number;
  attempts: number;
  maxAttempts: number;
  status: 'pending' | 'uploading' | 'completed' | 'failed';
  lastError?: string;
}

export interface UploadProgress {
  itemId: string;
  status: 'pending' | 'uploading' | 'completed' | 'failed';
  progress: number;
  error?: string;
}

class OfflineUploadQueue {
  private queue: UploadQueueItem[] = [];
  private isProcessing = false;
  private progressListeners: ((items: UploadProgress[]) => void)[] = [];
  private networkListeners: ((online: boolean) => void)[] = [];
  private retryInterval: number | null = null;

  constructor() {
    console.log('🔧 Initializing OfflineUploadQueue...');
    this.loadQueue();
    this.setupNetworkListener();
    this.startRetryLoop();
    
    // Process queue immediately if online
    setTimeout(() => {
      if (navigator.onLine && this.queue.length > 0) {
        console.log('🔄 Processing existing queue on initialization...');
        this.processQueue();
      }
    }, 1000); // Give some time for app to fully load
  }

  /**
   * Add a photo to the upload queue
   */
  addToQueue(item: Omit<UploadQueueItem, 'id' | 'timestamp' | 'attempts' | 'status'>): string {
    const queueItem: UploadQueueItem = {
      ...item,
      id: `upload_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      timestamp: Date.now(),
      attempts: 0,
      status: 'pending'
    };

    this.queue.push(queueItem);
    this.saveQueue();
    this.notifyProgressListeners();

    // Start processing if online
    if (navigator.onLine) {
      this.processQueue();
    }

    return queueItem.id;
  }

  /**
   * Get current queue status
   */
  getQueueStatus(): UploadProgress[] {
    return this.queue.map(item => ({
      itemId: item.id,
      status: item.status,
      progress: item.status === 'completed' ? 100 : 
                item.status === 'uploading' ? 50 : 0,
      error: item.lastError
    }));
  }

  /**
   * Subscribe to progress updates
   */
  onProgress(callback: (items: UploadProgress[]) => void): () => void {
    this.progressListeners.push(callback);
    return () => {
      this.progressListeners = this.progressListeners.filter(cb => cb !== callback);
    };
  }

  /**
   * Subscribe to network status changes
   */
  onNetworkChange(callback: (online: boolean) => void): () => void {
    this.networkListeners.push(callback);
    return () => {
      this.networkListeners = this.networkListeners.filter(cb => cb !== callback);
    };
  }

  /**
   * Process the upload queue
   */
  private async processQueue(): Promise<void> {
    if (this.isProcessing) {
      console.log('⏸️ Queue already processing, skipping...');
      return;
    }

    if (!navigator.onLine) {
      console.log('📴 Network offline, skipping queue processing');
      return;
    }

    this.isProcessing = true;
    const pendingItems = this.queue.filter(item => 
      item.status === 'pending' && item.attempts < item.maxAttempts
    );

    console.log(`🔄 Processing ${pendingItems.length} pending items from queue`);

    for (const item of pendingItems) {
      if (!navigator.onLine) {
        console.log('📴 Network went offline during processing, stopping...');
        break; // Stop processing if we go offline
      }

      await this.processItem(item);
    }

    this.isProcessing = false;
    this.removeCompletedItems();
    console.log('✅ Queue processing completed');
  }

  /**
   * Manually trigger queue processing (for debugging)
   */
  processQueueManually(): void {
    console.log('🔧 Manual queue processing requested');
    this.processQueue();
  }

  /**
   * Process a single queue item
   */
  private async processItem(item: UploadQueueItem): Promise<void> {
    item.status = 'uploading';
    item.attempts++;
    this.notifyProgressListeners();

    try {
      console.log(`🔄 Attempting upload ${item.attempts}/${item.maxAttempts} for item:`, item.id);
      
      // Convert base64 to WebP blob
      const response = await fetch(item.photoDataUrl);
      const originalBlob = await response.blob();
      const webpBlob = await this.convertToWebP(originalBlob, 0.8);

      // Upload to server
      const uploadResponse = await fetch(item.uploadUrl, {
        method: 'PUT',
        body: webpBlob,
        headers: {
          'Content-Type': 'image/webp',
        },
      });

      if (uploadResponse.ok) {
        item.status = 'completed';
        console.log('✅ Upload completed for item:', item.id);
      } else {
        throw new Error(`Upload failed: ${uploadResponse.status} ${uploadResponse.statusText}`);
      }

    } catch (error) {
      console.warn(`❌ Upload failed for item ${item.id}:`, error);
      item.lastError = error instanceof Error ? error.message : 'Unknown error';
      
      if (item.attempts >= item.maxAttempts) {
        item.status = 'failed';
        console.error(`💀 Max attempts reached for item ${item.id}`);
      } else {
        item.status = 'pending'; // Will retry
        console.log(`🔄 Will retry item ${item.id} (${item.attempts}/${item.maxAttempts})`);
      }
    }

    this.saveQueue();
    this.notifyProgressListeners();
  }

  /**
   * Convert image to WebP format
   */
  private convertToWebP(imageBlob: Blob, quality: number = 0.8): Promise<Blob> {
    return new Promise((resolve, reject) => {
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      const img = new Image();
      
      img.onload = () => {
        canvas.width = img.width;
        canvas.height = img.height;
        
        ctx?.drawImage(img, 0, 0);
        
        canvas.toBlob(
          (blob) => {
            if (blob) {
              resolve(blob);
            } else {
              reject(new Error('Failed to convert image to WebP'));
            }
          },
          'image/webp',
          quality
        );
      };
      
      img.onerror = () => reject(new Error('Failed to load image'));
      img.src = URL.createObjectURL(imageBlob);
    });
  }

  /**
   * Setup network status listener
   */
  private setupNetworkListener(): void {
    const handleOnline = () => {
      console.log('🌐 Network back online, processing upload queue...');
      this.networkListeners.forEach(callback => callback(true));
      
      // Process queue after a short delay to ensure network is stable
      setTimeout(() => {
        this.processQueue();
      }, 2000);
    };

    const handleOffline = () => {
      console.log('📴 Network offline, uploads will retry when connection returns');
      this.networkListeners.forEach(callback => callback(false));
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);
    
    // Log initial network status
    console.log(`🌐 Initial network status: ${navigator.onLine ? 'Online' : 'Offline'}`);
  }

  /**
   * Start retry loop for failed items
   */
  private startRetryLoop(): void {
    this.retryInterval = setInterval(() => {
      if (navigator.onLine && !this.isProcessing) {
        const hasRetryableItems = this.queue.some(item => 
          item.status === 'pending' && item.attempts < item.maxAttempts
        );
        
        if (hasRetryableItems) {
          console.log('🔄 Retry loop: Processing pending uploads...');
          this.processQueue();
        }
      }
    }, 30000); // Check every 30 seconds
  }

  /**
   * Remove completed items from queue
   */
  private removeCompletedItems(): void {
    const originalLength = this.queue.length;
    this.queue = this.queue.filter(item => 
      item.status !== 'completed' || Date.now() - item.timestamp < 300000 // Keep completed items for 5 minutes
    );
    
    if (this.queue.length !== originalLength) {
      this.saveQueue();
      this.notifyProgressListeners();
    }
  }

  /**
   * Notify progress listeners
   */
  private notifyProgressListeners(): void {
    const progress = this.getQueueStatus();
    this.progressListeners.forEach(callback => callback(progress));
  }

  /**
   * Save queue to localStorage
   */
  private saveQueue(): void {
    try {
      localStorage.setItem('upload-queue', JSON.stringify(this.queue));
    } catch (error) {
      console.warn('Failed to save upload queue:', error);
    }
  }

  /**
   * Load queue from localStorage
   */
  private loadQueue(): void {
    try {
      const saved = localStorage.getItem('upload-queue');
      if (saved) {
        this.queue = JSON.parse(saved);
        console.log(`📥 Loaded ${this.queue.length} items from localStorage:`, this.queue);
        
        // Reset uploading items to pending on app restart
        this.queue.forEach(item => {
          if (item.status === 'uploading') {
            item.status = 'pending';
            console.log(`🔄 Reset item ${item.id} from uploading to pending`);
          }
        });
        
        // Notify listeners immediately
        this.notifyProgressListeners();
      } else {
        console.log('📥 No upload queue found in localStorage');
      }
    } catch (error) {
      console.warn('❌ Failed to load upload queue:', error);
      this.queue = [];
    }
  }

  /**
   * Clear all completed and failed items
   */
  clearQueue(): void {
    this.queue = this.queue.filter(item => 
      item.status === 'pending' || item.status === 'uploading'
    );
    this.saveQueue();
    this.notifyProgressListeners();
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    if (this.retryInterval) {
      clearInterval(this.retryInterval);
      this.retryInterval = null;
    }
  }

  /**
   * Get queue statistics
   */
  getStats() {
    return {
      total: this.queue.length,
      pending: this.queue.filter(i => i.status === 'pending').length,
      uploading: this.queue.filter(i => i.status === 'uploading').length,
      completed: this.queue.filter(i => i.status === 'completed').length,
      failed: this.queue.filter(i => i.status === 'failed').length,
    };
  }
}

// Create singleton instance
export const uploadQueue = new OfflineUploadQueue();