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
      class="category-btn" 
      class:active={activeCategory === 'all'}
      onclick={() => handleCategoryClick('all')}
      style="background-image: url('https://storage.googleapis.com/micartapro-images/70895fae-85aa-47ea-a14a-dd76ea379f2e/1769720140-887099-empanadas-de-pino-e34e2d7e.png');"
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
    padding: 2rem 1rem;
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
    gap: 1rem;
    justify-content: center;
    flex-wrap: wrap;
  }

  .category-btn {
    padding: 1rem 1.5rem;
    border: none;
    background-size: cover;
    background-position: center;
    color: var(--white);
    border-radius: 12px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.95rem;
    transition: all 0.3s ease;
    white-space: nowrap;
    position: relative;
    overflow: hidden;
    min-height: 80px;
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

  @media (max-width: 768px) {
    .categories-nav {
      padding: 1rem 0.5rem;
      top: 4.75rem; /* Debajo de la topbar en móvil (evitar que tape las secciones) */
    }

    .categories-container {
      gap: 0.4rem;
    }

    .category-btn {
      padding: 0.5rem 0.8rem;
      font-size: 0.75rem;
      min-height: 50px;
    }
  }
</style>
