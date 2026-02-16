/**
 * Store de trackings de pedidos
 * Persiste en localStorage para acceso rápido. Soporta múltiples pedidos.
 */
import { writable, get } from 'svelte/store';

const STORAGE_KEY = 'cadorago_trackings';
const OLD_STORAGE_KEY = 'cadorago_current_tracking';
const MAX_TRACKINGS = 20; // Límite para no llenar localStorage

function normalizeEntry(entry) {
  if (typeof entry === 'string' && entry.trim()) {
    return { id: entry.trim(), menuId: null, createdAt: null, isDelivered: false };
  }
  if (entry && typeof entry === 'object' && entry.id) {
    return {
      id: String(entry.id).trim(),
      menuId: entry.menuId || null,
      createdAt: entry.createdAt || null,
      isDelivered: entry.isDelivered === true,
    };
  }
  return null;
}

function loadFromStorage() {
  try {
    if (typeof window === 'undefined' || !window.localStorage) return [];
    let saved = localStorage.getItem(STORAGE_KEY);
    // Migrar desde formato antiguo (clave y valor único)
    if (!saved) {
      const old = localStorage.getItem(OLD_STORAGE_KEY);
      if (old && old.trim()) {
        const migrated = [{ id: old.trim(), createdAt: null }];
        localStorage.setItem(STORAGE_KEY, JSON.stringify(migrated));
        localStorage.removeItem(OLD_STORAGE_KEY);
        return migrated;
      }
      return [];
    }
    const parsed = JSON.parse(saved);
    if (typeof parsed === 'string' && parsed.trim()) {
      return [normalizeEntry(parsed)];
    }
    if (!Array.isArray(parsed)) return [];
    return parsed.map(normalizeEntry).filter(Boolean);
  } catch (error) {
    console.warn('Error al cargar trackings:', error);
    return [];
  }
}

function saveToStorage(ids) {
  try {
    if (typeof window === 'undefined' || !window.localStorage) return;
    if (ids.length > 0) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(ids));
    } else {
      localStorage.removeItem(STORAGE_KEY);
    }
  } catch (error) {
    console.warn('Error al guardar trackings:', error);
  }
}

function createTrackingStore() {
  const w = writable(loadFromStorage());
  const { subscribe, set } = w;

  function addTracking(trackingId, createdAt = null, menuId = null) {
    const id = (trackingId && typeof trackingId === 'string' ? trackingId : String(trackingId || '')).trim();
    if (!id) return;
    const current = get(w) || [];
    const filtered = current.filter((x) => (typeof x === 'string' ? x : x.id) !== id);
    const entry = {
      id,
      menuId: menuId || (typeof trackingId === 'object' && trackingId?.menuId) || null,
      createdAt: createdAt || (typeof trackingId === 'object' && trackingId?.createdAt) || null,
      isDelivered: false,
    };
    const next = [entry, ...filtered.map((x) => (typeof x === 'object' ? x : normalizeEntry(x)))].slice(0, MAX_TRACKINGS);
    set(next);
    saveToStorage(next);
  }

  function updateTracking(trackingId, updates) {
    const current = get(w) || [];
    const next = current.map((e) => {
      const eid = typeof e === 'string' ? e : e.id;
      if (eid !== trackingId) return typeof e === 'object' ? e : normalizeEntry(e);
      const base = typeof e === 'object' ? e : { id: e, menuId: null, createdAt: null };
      return { ...base, ...updates };
    });
    set(next);
    saveToStorage(next);
  }

  function removeTracking(trackingId) {
    const current = get(w) || [];
    const next = current.filter((x) => (typeof x === 'string' ? x : x.id) !== trackingId);
    set(next);
    saveToStorage(next);
  }

  /** Obtiene trackings opcionalmente filtrados por menuId */
  function getTrackings(menuId = null) {
    const list = get(w) || [];
    const normalized = list.map((x) => (typeof x === 'object' ? x : normalizeEntry(x)));
    if (!menuId) return normalized;
    return normalized.filter((e) => e.menuId === menuId);
  }

  return {
    subscribe,
    addTracking,
    setTracking: addTracking,
    updateTracking,
    removeTracking,
    getTrackings,
    clear: () => {
      set([]);
      saveToStorage([]);
    },
  };
}

export const trackingStore = createTrackingStore();
