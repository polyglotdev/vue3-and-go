<template>
  <Header />
  <div>
    <router-view />
  </div>
  <Footer />
</template>

<script>
  import Header from './components/Header.vue'
  import Body from './components/Body.vue'
  import Footer from './components/Footer.vue'

  export default {
    name: 'Frontend',
    components: {
      Header,
      Body,
      Footer
    }
  }
</script>

<style scoped></style>
