import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'
import { changeLanguage } from './lib/useLanguage'

// Modo demo (iframe de Remotion): limpiar storage, agrandar texto y asegurar scroll
if (typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('demo') === '1') {
  try {
    localStorage.clear()
    sessionStorage.clear()
    document.cookie.split(';').forEach((c) => {
      const name = c.trim().split('=')[0]
      if (name) document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/`
    })
    // Zoom 40% más grande para el video
    document.documentElement.style.zoom = '1.4'
    // Asegurar que el documento sea scrollable dentro del iframe
    document.documentElement.style.overflowY = 'auto'
    document.documentElement.style.height = 'auto'
    document.body.style.overflowY = 'auto'
    document.body.style.minHeight = '100%'
  } catch (_) { /* ignore */ }
}

// Listen for commands from parent (e.g. Remotion video)
const isDemo = typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('demo') === '1'
window.addEventListener('message', (e) => {
  if (!isDemo) {
    const allowedOrigins = ['http://localhost:', 'http://127.0.0.1:', 'https://localhost:', 'https://127.0.0.1:']
    if (!allowedOrigins.some(o => e.origin?.startsWith(o))) return
  }
  if (typeof e.data?.scrollY === 'number') {
    const y = e.data.scrollY
    const doScroll = () => {
      const scrollEl = document.scrollingElement || document.documentElement
      if (scrollEl) scrollEl.scrollTo({ top: y, behavior: 'auto' })
      else window.scrollTo({ top: y, behavior: 'auto' })
    }
    doScroll()
    setTimeout(doScroll, 150)
    setTimeout(doScroll, 400)
  }
  if (e.data?.changeLanguage && ['EN', 'ES', 'PT'].includes(e.data.changeLanguage)) {
    changeLanguage(e.data.changeLanguage)
  }
  if (e.data?.scrollToElement) {
    const selector = e.data.scrollToElement
    const run = () => {
      const el = typeof selector === 'string'
        ? Array.from(document.querySelectorAll('[data-item-title], .chef-card, [class*="menu-item"], [id]'))
            .find(n => n.textContent?.toLowerCase().includes(selector.toLowerCase()))
        : document.querySelector(selector)
      if (el) el.scrollIntoView({ behavior: 'auto', block: 'center' })
    }
    run()
    // Reintentar si el DOM aún no está listo (ej. menú cargando)
    setTimeout(run, 200)
    setTimeout(run, 600)
  }
  if (e.data?.clickAddToCart) {
    let btn
    if (typeof e.data.clickAddToCart === 'string') {
      const title = e.data.clickAddToCart.toLowerCase().trim()
      const searchTerms = title.split(/\s+/).filter(Boolean)
      // Demo: selector directo por data-remotion-add (más fiable)
      btn = Array.from(document.querySelectorAll('button[data-remotion-add]'))
        .find(x => (x.getAttribute('data-remotion-add') || '').toLowerCase().includes(title)) || null
      if (!btn) {
        btn = Array.from(document.querySelectorAll('button[data-item-title], button.cta-btn, [aria-label="Agregar al carrito"]'))
          .find(b => (b.dataset?.itemTitle || '').toLowerCase().includes(title))
      }
      if (!btn) {
        btn = Array.from(document.querySelectorAll('button.cta-btn, [aria-label="Agregar al carrito"]'))
          .find(b => {
            const cardText = b.closest('.chef-card, [class*="menu-item"]')?.textContent?.toLowerCase() || ''
            return searchTerms.every(t => cardText.includes(t))
          })
      }
      if (!btn) {
        btn = Array.from(document.querySelectorAll('button.cta-btn, [aria-label="Agregar al carrito"]'))
          .find(b => (b.closest('.chef-card, [class*="menu-item"]')?.textContent || '').toLowerCase().includes(title))
      }
      // Fallback: botón "Agregar" / "Add" dentro de una card que contiene el nombre del producto
      if (!btn) {
        btn = Array.from(document.querySelectorAll('button'))
          .find(b => {
            const txt = (b.textContent || '').toLowerCase()
            const isAddBtn = txt.includes('agregar') || txt.includes('add') || txt.includes('adicionar')
            const card = b.closest('.chef-card, [class*="menu-item"]')
            const cardText = (card?.textContent || '').toLowerCase()
            return isAddBtn && cardText.includes(title)
          })
      }
      // Fallback: buscar la card por título y tomar el botón .cta-btn dentro
      if (!btn) {
        const card = Array.from(document.querySelectorAll('.chef-card, .menu-item'))
          .find(c => (c.textContent || '').toLowerCase().includes(title))
        if (card) btn = card.querySelector('button.cta-btn, button[data-item-title]')
      }
    } else {
      const index = (e.data.clickAddToCart === true || e.data.clickAddToCart === 1) ? 0 : (e.data.clickAddToCart - 1)
      const btns = document.querySelectorAll('button.cta-btn, [aria-label="Agregar al carrito"]')
      btn = btns[index]
    }
    if (btn) {
      btn.scrollIntoView({ behavior: 'auto', block: 'center' })
      setTimeout(() => btn.click(), 80)
    } else if (typeof e.data.clickAddToCart === 'string') {
      const title = e.data.clickAddToCart
      const tryClick = () => {
        let retryBtn = Array.from(document.querySelectorAll('button[data-remotion-add]'))
          .find(x => (x.getAttribute('data-remotion-add') || '').toLowerCase().includes(title.toLowerCase()))
        if (!retryBtn) {
          retryBtn = Array.from(document.querySelectorAll('button.cta-btn, button[data-item-title]'))
            .find(b => (b.dataset?.itemTitle || b.closest('.chef-card')?.textContent || '').toLowerCase().includes(title.toLowerCase()))
        }
        if (retryBtn) {
          retryBtn.scrollIntoView({ behavior: 'auto', block: 'center' })
          setTimeout(() => retryBtn.click(), 50)
        }
      }
      setTimeout(tryClick, 300)
      setTimeout(tryClick, 600)
      setTimeout(tryClick, 1000)
    }
  }
  if (e.data?.selectAndAddFromSheet) {
    let retries = 0
    const useSides = e.data.selectAndAddFromSheet === 'sides'
    const run = () => {
      const findSheet = () => document.querySelector('.acompanamiento-sheet--visible') || document.querySelector('.acompanamiento-sheet')
      const sheet = findSheet()
      if (!sheet && retries < 40) {
        retries++
        setTimeout(run, 100)
        return
      }
      if (!sheet) return
      const addBtn = () => document.querySelector('.acompanamiento-sheet-add-btn:not([disabled])')
      const sides = sheet.querySelectorAll('.acompanamiento-sheet-option')
      const groups = sheet.querySelectorAll('.selectable-group')
      let delay = 0
      if (useSides && sides.length > 0) {
        sides[0].scrollIntoView({ behavior: 'auto', block: 'center' })
        sides[0].click()
        delay = 250
      }
      setTimeout(() => {
        const sheet2 = findSheet()
        if (!sheet2) return
        const groups2 = sheet2.querySelectorAll('.selectable-group')
        groups2.forEach((group, i) => {
          setTimeout(() => {
            const opt = group.querySelector('.selectable-option:not(.selectable-option--disabled)')
            if (opt) {
              opt.scrollIntoView({ behavior: 'auto', block: 'center' })
              const input = opt.querySelector('input')
              if (input) input.click()
              else opt.click()
            }
          }, i * 150)
        })
        const addDelay = Math.max(groups2.length * 150 + 200, 500)
        setTimeout(() => {
          const btn = addBtn()
          if (btn) {
            btn.scrollIntoView({ behavior: 'auto', block: 'center' })
            btn.click()
          }
        }, addDelay)
      }, delay)
    }
    setTimeout(run, 450)
  }
  if (e.data?.clickRealizarPedido) {
    const btn = document.querySelector('[aria-label="Realizar pedido"]')
    if (btn) btn.click()
  }
  if (e.data?.selectDelivery) {
    let retries = 0
    const run = () => {
      const btn = Array.from(document.querySelectorAll('button')).find(b =>
        b.textContent?.includes('📦') || b.textContent?.toLowerCase().includes('domicilio') || b.textContent?.toLowerCase().includes('delivery')
      )
      if (btn) {
        btn.click()
      } else if (retries < 15) {
        retries++
        setTimeout(run, 100)
      }
    }
    setTimeout(run, 200)
  }
  if (typeof e.data?.typeDeliveryAddress === 'string') {
    let retries = 0
    const run = () => {
      const input = document.querySelector('#delivery-address')
      if (input) {
        input.value = e.data.typeDeliveryAddress
        input.dispatchEvent(new Event('input', { bubbles: true }))
      } else if (retries < 30) {
        retries++
        setTimeout(run, 100)
      }
    }
    setTimeout(run, 300)
  }
  if (e.data?.selectFirstAddressSuggestion) {
    let retries = 0
    const run = () => {
      const btn = document.querySelector('button[data-address-suggestion]')
      if (btn) {
        btn.click()
      } else if (retries < 40) {
        retries++
        setTimeout(run, 150)
      }
    }
    setTimeout(run, 500)
  }
  if (e.data?.clickDeliveryNextStep) {
    let retries = 0
    const run = () => {
      const btn = document.querySelector('button[data-delivery-next]:not([disabled])')
      if (btn) {
        btn.click()
      } else if (retries < 30) {
        retries++
        setTimeout(run, 150)
      }
    }
    setTimeout(run, 400)
  }
  if (e.data?.fillDeliveryContactForm) {
    const data = e.data.fillDeliveryContactForm
    let retries = 0
    const run = () => {
      const nameInput = document.querySelector('#nombre-retiro')
      const phoneInput = document.querySelector('#phone-contact-delivery')
      const emailInput = document.querySelector('#email-contact-delivery')
      const addressNumInput = document.querySelector('#address-number')
      if (nameInput && phoneInput && emailInput) {
        if (data?.name) { nameInput.value = data.name; nameInput.dispatchEvent(new Event('input', { bubbles: true })) }
        if (data?.phone) { phoneInput.value = data.phone; phoneInput.dispatchEvent(new Event('input', { bubbles: true })) }
        if (data?.email) { emailInput.value = data.email; emailInput.dispatchEvent(new Event('input', { bubbles: true })) }
        if (data?.addressNumber && addressNumInput) { addressNumInput.value = data.addressNumber; addressNumInput.dispatchEvent(new Event('input', { bubbles: true })) }
      } else if (retries < 40) {
        retries++
        setTimeout(run, 100)
      }
    }
    setTimeout(run, 500)
  }
  if (e.data?.clickPlaceOrder) {
    let retries = 0
    const run = () => {
      const btn = document.querySelector('button[data-place-order]:not([disabled])')
      if (btn) {
        btn.click()
      } else if (retries < 30) {
        retries++
        setTimeout(run, 150)
      }
    }
    setTimeout(run, 400)
  }
})

const app = mount(App, {
  target: document.getElementById('app'),
})

export default app
