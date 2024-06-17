<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <h1 class="mt-5">Login</h1>
        <hr />
        <form-tag @myevent="submitHandler" name="myform" event="myevent">
          <text-input
            v-model="email"
            label="Email"
            type="email"
            name="email"
            required="true"
            aria-describedby="emailHelp"
          >
          </text-input>
          <div id="emailHelp" aria-live="polite">
            <h6>We'll never share your email with anyone else.</h6>
          </div>

          <text-input
            v-model="password"
            label="Password"
            type="password"
            name="password"
            required="true"
          >
          </text-input>
          <check-input label="Remember me" name="remember"></check-input>
          <hr />
          <button class="btn btn-primary" aria-label="Submit button">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              width="24"
              height="24"
              fill="currentColor"
            >
              <path
                d="M12 2a10 10 0 100 20 10 10 0 000-20zm0 18a8 8 0 110-16 8 8 0 010 16zm-.75-11.47l5 5a.75.75 0 01-1.06 1.06l-4.22-4.22-4.22 4.22a.75.75 0 11-1.06-1.06l5-5a.75.75 0 011.06 0z"
              ></path>
            </svg>
            Submit
          </button>
        </form-tag>
      </div>
    </div>
  </div>
</template>

<script>
  import FormTag from './forms/FormTag.vue'
  import TextInput from './forms/TextInput.vue'
  import CheckInput from './forms/CheckInput.vue'

  export default {
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'login',
    components: {
      FormTag,
      TextInput,
      CheckInput
    },
    data() {
      return {
        email: '',
        password: ''
      }
    },
    methods: {
      submitHandler() {
        console.log('submitHandler called - success!')

        const payload = {
          email: this.email,
          password: this.password
        }

        const requestOptions = {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(payload)
        }

        fetch('http://localhost:8081/users/login', requestOptions)
          .then((response) => response.json())
          .then((data) => {
            if (data.error) {
              console.log('Error:', data.message)
            } else {
              console.log(data)
            }
          })
      }
    }
  }
</script>

<style scoped>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 10px 20px;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
  }

  .btn svg {
    margin-right: 8px;
  }

  .btn:hover {
    background-color: #0056b3;
  }
</style>
