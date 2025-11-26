<script>
  import logo from '../../assets/logo.png';
  
  /** @type {string} */
  export let className = 'h-16 w-auto';
  /** @type {string} */
  export let alt = 'Company Logo';
</script>

<img src={logo} {alt} class={className} />
