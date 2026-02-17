<script>
  import { getMultilingualText } from '../../lib/multilingual';
  
  let activeCategory = $state('all');
  
  const { categories = [], onCategoryChange } = $props();
  
  function handleCategoryClick(category) {
    activeCategory = category;
    onCategoryChange?.(category);
  }
</script>

<nav class="categories-nav">
  <div class="categories-container">
    <button 
      class="category-btn category-btn--dark" 
      class:active={activeCategory === 'all'}
      onclick={() => handleCategoryClick('all')}
    >
      <span>Todos</span>
    </button>
    {#each categories as category}
      <button 
        class="category-btn" 
        class:active={activeCategory === category.id}
        onclick={() => handleCategoryClick(category.id)}
        style={category.image ? `background-image: url('${category.image}');` : ''}
      >
        <span>{getMultilingualText(category.title)}</span>
      </button>
    {/each}
  </div>
</nav>

<style>
  .categories-nav {
    background-color: var(--white);
    padding: 1rem;
    border-bottom: 1px solid var(--border-light);
    overflow-x: auto;
    position: sticky;
    top: 3.25rem; /* Debajo del header sticky (aprox. altura de la topbar) */
    z-index: 90;
    margin-top: 0;
  }

  .categories-container {
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    gap: 0.75rem;
    justify-content: flex-start;
    flex-wrap: nowrap;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    /* Ocultar scrollbar manteniendo scroll (touch, rueda, trackpad) */
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .categories-container::-webkit-scrollbar {
    display: none;
  }

  .category-btn {
    padding: 0.85rem 1.25rem;
    border: none;
    background-size: cover;
    background-position: center;
    color: var(--white);
    border-radius: 12px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.9rem;
    transition: all 0.3s ease;
    white-space: nowrap;
    position: relative;
    overflow: hidden;
    min-height: 72px;
    min-width: 120px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  .category-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.4);
    z-index: 1;
    transition: background-color 0.3s ease;
  }

  .category-btn span {
    position: relative;
    z-index: 2;
  }

  .category-btn:hover,
  .category-btn.active {
    transform: translateY(-4px);
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  }

  .category-btn:hover::before,
  .category-btn.active::before {
    background-color: rgba(0, 0, 0, 0.5);
  }

  /* Todos: oscuro, sin imagen de fondo */
  .category-btn--dark {
    background-color: #374151;
    background-image: none;
  }

  .category-btn--dark::before {
    background-color: rgba(0, 0, 0, 0.3);
  }

  .category-btn--dark:hover::before,
  .category-btn--dark.active::before {
    background-color: rgba(0, 0, 0, 0.5);
  }

  @media (max-width: 768px) {
    .categories-nav {
      padding: 0.75rem 0.5rem;
      top: 4.75rem; /* Debajo de la topbar en móvil (evitar que tape las secciones) */
    }

    .categories-container {
      gap: 0.5rem;
    }

    .category-btn {
      padding: 0.65rem 1rem;
      font-size: 0.8rem;
      min-height: 58px;
      min-width: 100px;
    }
  }
</style>
