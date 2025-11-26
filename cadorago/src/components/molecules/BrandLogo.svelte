<script>
  import Logo from '../atoms/Logo.svelte';
  import Birds from '../atoms/Birds.svelte';
  import LogoText from '../atoms/LogoText.svelte';
  
  /** @type {'micro' | 'navbar' | 'footer' | 'hero' | 'xl' | 'xxl' | 'custom'} */
  export let size = 'hero';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
  
  // Size configurations with proportional scaling
  const sizeConfigs = {
    micro: {
      logoHeight: 'h-12',
      logoMargin: 'ml-1.5',
      birdsHeight: 'h-12',
      birdsPosition: '-top-3 -right-3',
      textWidth: 'w-16',
      textMargin: '-mt-[1.9rem]'
    },
    navbar: {
      logoHeight: 'h-20',
      logoMargin: 'ml-3',
      birdsHeight: 'h-20',
      birdsPosition: '-top-6 -right-6',
      textWidth: 'w-28',
      textMargin: '-mt-[3.2rem]'
    },
    footer: {
      logoHeight: 'h-32',
      logoMargin: 'ml-4',
      birdsHeight: 'h-32',
      birdsPosition: '-top-9 -right-9',
      textWidth: 'w-44',
      textMargin: '-mt-[5rem]'
    },
    hero: {
      logoHeight: 'h-56',
      logoMargin: 'ml-7',
      birdsHeight: 'h-56',
      birdsPosition: '-top-16 -right-16',
      textWidth: 'w-80',
      textMargin: '-mt-[9rem]'
    },
    xl: {
      logoHeight: 'h-[24rem]',
      logoMargin: 'ml-14',
      birdsHeight: 'h-[24rem]',
      birdsPosition: '-top-32 -right-32',
      textWidth: 'w-[36rem]',
      textMargin: '-mt-[15.5rem]'
    },
    xxl: {
      logoHeight: 'h-[32rem]',
      logoMargin: 'ml-16',
      birdsHeight: 'h-[32rem]',
      birdsPosition: '-top-36 -right-48',
      textWidth: 'w-[42rem]',
      textMargin: '-mt-[20.5rem]'
    }
  };
  
  $: config = sizeConfigs[size] || sizeConfigs.hero;
</script>

<div class={`relative inline-block ${className}`}>
  <!-- Logo with Birds positioned in top right corner -->
  <div class="relative"> 
    <Logo className="{config.logoHeight} w-auto {config.logoMargin}" />
    <div class="absolute {config.birdsPosition}">
      <Birds className="{config.birdsHeight} w-auto" />
    </div>
  </div>
  
  <!-- LogoText below the logo with negative margin to overlap -->
  <div class="{config.textMargin} relative z-10">
    <LogoText className="{config.textWidth}" />
  </div>
</div>
