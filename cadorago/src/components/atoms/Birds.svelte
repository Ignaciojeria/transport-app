<script>
  import birds from '../../assets/birds.png';
  
  /** @type {string} */
  export let className = 'h-16 w-auto';
  /** @type {string} */
  export let alt = 'Birds';
</script>

<img src={birds} {alt} class={className} />
