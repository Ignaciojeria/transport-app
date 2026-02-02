<script>
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import { language, changeLanguage } from '../../lib/useLanguage';
  
  const restaurantData = $derived(restaurantDataStore.value);
  const businessName = $derived(restaurantData?.businessInfo?.businessName || '');
  const whatsapp = $derived(restaurantData?.businessInfo?.whatsapp || '');
  const businessHours = $derived(restaurantData?.businessInfo?.businessHours || []);
  const currentLang = $derived($language);
  
  // Formatear horarios para mostrar en el header
  const formattedHours = $derived(
    businessHours.length === 0 ? '' : businessHours.join(' | ')
  );
</script>

<!-- Topbar completa sticky al hacer scroll -->
<header class="modern-header">
  <div class="header-content">
    {#if businessName}
      <h1 class="header-title">{businessName}</h1>
    {/if}
    <div class="header-contact">
      {#if whatsapp}
        <a href={`https://wa.me/${whatsapp.replace(/\D/g, '')}`} target="_blank" rel="noopener">
          📱 {whatsapp}
        </a>
      {/if}
      {#if formattedHours}
        <span>⏰ {formattedHours}</span>
      {/if}
      <div class="language-selector">
        <button 
          class="lang-btn" 
          class:active={currentLang === 'PT'}
          onclick={() => changeLanguage('PT')}
          title="Português"
        >
          🇧🇷
        </button>
        <button 
          class="lang-btn" 
          class:active={currentLang === 'ES'}
          onclick={() => changeLanguage('ES')}
          title="Español"
        >
          🇪🇸
        </button>
        <button 
          class="lang-btn" 
          class:active={currentLang === 'EN'}
          onclick={() => changeLanguage('EN')}
          title="English"
        >
          🇺🇸
        </button>
      </div>
    </div>
  </div>
</header>

<style>
  .modern-header {
    background-color: var(--white);
    border-bottom: 1px solid var(--border-light);
    padding: 0.4rem 0.6rem;
    position: sticky;
    top: 0;
    z-index: 100;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  }

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .header-title {
    font-family: 'Playfair Display', serif;
    font-size: 1.5rem;
    font-weight: 800;
    margin: 0;
    color: var(--text-dark);
    white-space: nowrap;
  }

  .header-contact {
    margin: 0;
    display: flex;
    justify-content: flex-end;
    gap: 0.6rem;
    flex-wrap: wrap;
    font-size: 0.95rem;
  }

  .header-contact a {
    color: var(--text-light);
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 0.4rem;
    transition: color 0.3s;
  }

  .header-contact a:hover {
    color: var(--primary);
  }

  .language-selector {
    display: flex;
    gap: 0.3rem;
    align-items: center;
  }

  .lang-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1.25rem;
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
    transition: background-color 0.3s;
  }

  .lang-btn:hover {
    background-color: var(--secondary);
  }

  .lang-btn.active {
    background-color: var(--primary);
    color: var(--white);
  }

  @media (max-width: 768px) {
    .modern-header {
      padding: 0.35rem 0.5rem;
    }

    .header-content {
      flex-direction: column;
      gap: 0.4rem;
    }

    .header-title {
      font-size: 1.35rem;
    }

    .header-contact {
      justify-content: center;
      gap: 0.5rem;
      font-size: 0.9rem;
    }

    .lang-btn {
      font-size: 1.15rem;
      padding: 0.15rem 0.35rem;
    }
  }
</style>
