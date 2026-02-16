/**
 * Store para navegación - guarda la última ruta del menú visitada
 */
import { writable } from 'svelte/store';

export const lastMenuPath = writable('/');
