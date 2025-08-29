/**
 * Cola de sincronización general para todas las interacciones del transportista
 * cuando está offline y se reconecta a internet
 */

export interface SyncQueueItem {
  id: string;
  type: 'route_start' | 'route_stop' | 'delivery_status' | 'delivery_evidence' | 'non_delivery_evidence' | 'license_set';
  routeId: string;
  data: any;
  timestamp: number;
  attempts: number;
  maxAttempts: number;
  status: 'pending' | 'syncing' | 'completed' | 'failed';
  lastError?: string;
  priority: number; // 1-10, donde 1 es máxima prioridad
}

export interface SyncProgress {
  itemId: string;
  type: string;
  status: 'pending' | 'syncing' | 'completed' | 'failed';
  error?: string;
}

class SyncQueue {
  private queue: SyncQueueItem[] = [];
  private isProcessing = false;
  private progressListeners: ((items: SyncProgress[]) => void)[] = [];
  private networkListeners: ((online: boolean) => void)[] = [];
  private retryInterval: number | null = null;

  constructor() {
    console.log('🔧 Initializing SyncQueue...');
    this.loadQueue();
    this.setupNetworkListener();
    this.startRetryLoop();
    
    // Process queue immediately if online
    setTimeout(() => {
      if (navigator.onLine && this.queue.length > 0) {
        console.log('🔄 Processing existing sync queue on initialization...');
        this.processQueue();
      }
    }, 1000);
  }

  /**
   * Add an action to the sync queue
   */
  addToQueue(item: Omit<SyncQueueItem, 'id' | 'timestamp' | 'attempts' | 'status'>): string {
    const queueItem: SyncQueueItem = {
      ...item,
      id: `sync_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      timestamp: Date.now(),
      attempts: 0,
      status: 'pending'
    };

    // Insert by priority (higher priority first)
    const insertIndex = this.queue.findIndex(existing => existing.priority > item.priority);
    if (insertIndex === -1) {
      this.queue.push(queueItem);
    } else {
      this.queue.splice(insertIndex, 0, queueItem);
    }

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
  getQueueStatus(): SyncProgress[] {
    return this.queue.map(item => ({
      itemId: item.id,
      type: item.type,
      status: item.status,
      error: item.lastError
    }));
  }

  /**
   * Subscribe to progress updates
   */
  onProgress(callback: (items: SyncProgress[]) => void): () => void {
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
   * Process the sync queue
   */
  private async processQueue(): Promise<void> {
    if (this.isProcessing) {
      console.log('⏸️ Sync queue already processing, skipping...');
      return;
    }

    if (!navigator.onLine) {
      console.log('📴 Network offline, skipping sync queue processing');
      return;
    }

    this.isProcessing = true;
    const pendingItems = this.queue
      .filter(item => item.status === 'pending' && item.attempts < item.maxAttempts)
      .sort((a, b) => a.priority - b.priority); // Process by priority

    console.log(`🔄 Processing ${pendingItems.length} pending items from sync queue`);

    for (const item of pendingItems) {
      if (!navigator.onLine) {
        console.log('📴 Network went offline during sync processing, stopping...');
        break;
      }

      await this.processItem(item);
    }

    this.isProcessing = false;
    this.removeCompletedItems();
    console.log('✅ Sync queue processing completed');
  }

  /**
   * Process a single queue item
   */
  private async processItem(item: SyncQueueItem): Promise<void> {
    item.status = 'syncing';
    item.attempts++;
    this.notifyProgressListeners();

    try {
      console.log(`🔄 Attempting sync ${item.attempts}/${item.maxAttempts} for item:`, item.id, item.type);
      
      // Simulate sync based on type - aquí iría la lógica real de sincronización con GunJS
      await this.syncToGun(item);

      item.status = 'completed';
      console.log('✅ Sync completed for item:', item.id);

    } catch (error) {
      console.warn(`❌ Sync failed for item ${item.id}:`, error);
      item.lastError = error instanceof Error ? error.message : 'Unknown error';
      
      if (item.attempts >= item.maxAttempts) {
        item.status = 'failed';
        console.error(`💀 Max attempts reached for sync item ${item.id}`);
      } else {
        item.status = 'pending'; // Will retry
        console.log(`🔄 Will retry sync item ${item.id} (${item.attempts}/${item.maxAttempts})`);
      }
    }

    this.saveQueue();
    this.notifyProgressListeners();
  }

  /**
   * Sync item to GunJS based on type
   */
  private async syncToGun(item: SyncQueueItem): Promise<void> {
    // Import GunJS functions dynamically to avoid circular dependencies
    const { setRouteStarted, setDeliveryStatus, setDeliveryEvidence, setNonDeliveryEvidence, setRouteLicense } = 
      await import('../db/driver-gun-state');

    switch (item.type) {
      case 'route_start':
        setRouteStarted(item.routeId, item.data.started);
        break;

      case 'route_stop':
        setRouteStarted(item.routeId, false);
        break;

      case 'delivery_status':
        setDeliveryStatus(
          item.routeId, 
          item.data.visitIndex, 
          item.data.orderIndex, 
          item.data.unitIndex, 
          item.data.status
        );
        break;

      case 'delivery_evidence':
        await setDeliveryEvidence(
          item.routeId,
          item.data.visitIndex,
          item.data.orderIndex,
          item.data.unitIndex,
          item.data.evidence,
          item.data.routeData
        );
        break;

      case 'non_delivery_evidence':
        await setNonDeliveryEvidence(
          item.routeId,
          item.data.visitIndex,
          item.data.orderIndex,
          item.data.unitIndex,
          item.data.evidence,
          item.data.routeData
        );
        break;

      case 'license_set':
        setRouteLicense(item.routeId, item.data.license);
        break;

      default:
        throw new Error(`Unknown sync type: ${item.type}`);
    }

    // Simular delay de red
    await new Promise(resolve => setTimeout(resolve, 500));
  }

  /**
   * Setup network status listener
   */
  private setupNetworkListener(): void {
    const handleOnline = () => {
      console.log('🌐 Network back online, processing sync queue...');
      this.networkListeners.forEach(callback => callback(true));
      
      // Process queue after a short delay to ensure network is stable
      setTimeout(() => {
        this.processQueue();
      }, 2000);
    };

    const handleOffline = () => {
      console.log('📴 Network offline, sync operations will retry when connection returns');
      this.networkListeners.forEach(callback => callback(false));
    };

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);
    
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
          console.log('🔄 Retry loop: Processing pending sync operations...');
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
      localStorage.setItem('sync-queue', JSON.stringify(this.queue));
    } catch (error) {
      console.warn('Failed to save sync queue:', error);
    }
  }

  /**
   * Load queue from localStorage
   */
  private loadQueue(): void {
    try {
      const saved = localStorage.getItem('sync-queue');
      if (saved) {
        this.queue = JSON.parse(saved);
        console.log(`📥 Loaded ${this.queue.length} items from sync queue:`, this.queue);
        
        // Reset syncing items to pending on app restart
        this.queue.forEach(item => {
          if (item.status === 'syncing') {
            item.status = 'pending';
            console.log(`🔄 Reset sync item ${item.id} from syncing to pending`);
          }
        });
        
        this.notifyProgressListeners();
      } else {
        console.log('📥 No sync queue found in localStorage');
      }
    } catch (error) {
      console.warn('❌ Failed to load sync queue:', error);
      this.queue = [];
    }
  }

  /**
   * Clear completed and failed items
   */
  clearQueue(): void {
    this.queue = this.queue.filter(item => 
      item.status === 'pending' || item.status === 'syncing'
    );
    this.saveQueue();
    this.notifyProgressListeners();
  }

  /**
   * Get queue statistics
   */
  getStats() {
    return {
      total: this.queue.length,
      pending: this.queue.filter(i => i.status === 'pending').length,
      syncing: this.queue.filter(i => i.status === 'syncing').length,
      completed: this.queue.filter(i => i.status === 'completed').length,
      failed: this.queue.filter(i => i.status === 'failed').length,
    };
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
}

// Create singleton instance
export const syncQueue = new SyncQueue();