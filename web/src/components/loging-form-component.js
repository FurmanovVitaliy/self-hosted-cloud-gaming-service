import { log } from "../common/log.js";
import { authApi } from "../api/index.js";
async function loadIonicons() {
    await import('https://unpkg.com/ionicons@7.1.0/dist/ionicons/ionicons.esm.js');
}

class LoginForm extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });
    const wrapper = document.createElement("section");
    wrapper.innerHTML = `
    <div class="login-box">
   <form action="">
      <h2>Login</h2>
      <div class="input-box">
         <span class="icon">
            <ion-icon name="mail"></ion-icon>
         </span>
         <input type="email" required>
         <label>Email</label>
      </div>
      <div class="input-box">
         <span class="icon">
            <ion-icon name="lock-closed"></ion-icon>
         </span>
         <input type="password" required>
         <label>Password</label>
      </div>
      <div class="remember-forgot">
         <label><input type="checkbox">
         Remember me</label>
         <a href="#">Forget Password</a>
      </div>
      <button type="submit">Login</button>
      <div class="register-link">
         <p>Don't have an account? <a href="#">SingUp</a></p>
      </div>
   </form>
    </div>
    
    `;

                const style = document.createElement("style");
    style.textContent = `
    section{
      input:-webkit-autofill,
      input:-webkit-autofill:hover, 
      input:-webkit-autofill:focus, 
      input:-webkit-autofill:active {
        -webkit-text-fill-color: white !important;
        -webkit-background-clip: text !important;

      }

  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100%;

  background: url('src/img/bg.jpg') no-repeat center;
  background-size: cover;

}


.login-box{
  position: relative;
  width: 400px;
  height: 500px;
  border-radius: 20px;
  background: transparent;
  border: 2px solid rgb(255, 255, 255, .5);
  display: flex;
  backdrop-filter: blur(15px);
  justify-content: center;
  align-items: center;
}
h2{
  font-size: 2em;
  color: white;
  text-align: center;
}
.input-box{
  position: relative;
  width: 310px;
  margin: 30px 0;
  border-bottom: 2px solid white;
}
.input-box label{
  position: absolute;
  color: white;
  top: 50%;
  left: 5px;
  transform: translateY(-50%);
  font-size: 1em;
  pointer-events: none;
  transition: .5s;
}
.input-box input:focus ~ label,
.input-box input:valid ~ label{
  top: -5px;
  color: white;
}

.input-box input.has-content:invalid ~ label{
  top: -5px;
  color: red;
}


.input-box input{
  width: 100%;
  height: 50px;
  padding: 0 35px 0 5px;
  font-size: 1em;
  border: none;
  outline: none;
  background: transparent;
  color: white;
}
.input-box .icon{
  position: absolute;
  right:8px;
  font-size: 1.2em;
  line-height: 57px;
  color: white;
}
.remember-forgot{
  font-size: .9em;
  color: white;
  margin: -15px 0 15px;
  display: flex;
  justify-content: space-between;
}
.remremember-forgot label input{
  margin: 3px;  
}
.remember-forgot a{
  color: white;
  text-decoration: none;
}
.remember-forgot a:hover{
  text-decoration: underline;
}
button{
  width: 100%;
  height: 40px;
  border: none;
  outline: none;
  border-radius: 40px;
  background: white;
  color:black;
  font-size: 1em;
  cursor: pointer;
  font-weight: 500;
}
.register-link{
  font-size: .9em;
  color: white;
  margin: 15px 0 0;
  text-align: center;
  margin: 25px 0 10px;
}
.register-link p a{
  color: white;
  text-decoration: none;
  font-weight: 600;
}
.register-link p a:hover{
  text-decoration: underline;
}
@media (max-width: 360px) {
  .login-box{
      width: 100%;
      height: 100%;
      border: none;
      border-radius: 0;
  }
  .input-box{
      width: 290px;
  }
  
}
`;

    shadow.appendChild(style);
    shadow.appendChild(wrapper);
  }

  connectedCallback() {
    loadIonicons();
    const shadow = this.shadowRoot;


    const wrapper = shadow.querySelector(".login-box");
    const title = shadow.querySelector(".login-box h2");

    const submitButton = shadow.querySelector(".login-box button[type='submit']");
  

    const inputLogin = shadow.querySelector(".input-box input[type='email']")
    inputLogin.addEventListener("input", (e) => { 
      if (e.target.value) {
        e.target.classList.add('has-content');
      }
      else {
        e.target.classList.remove('has-content');
      }
    }
  );

    const inputPassword = shadow.querySelector(".input-box input[type='password']");

    submitButton.addEventListener("click", (e) => {
      e.preventDefault();

      if (inputLogin && inputPassword) {
        authApi
          .login(inputLogin.value, inputPassword.value)
          .then((response) => {
            response
              .json()
              .then((data) => {
                if (data) {
                  localStorage.setItem("username", data.username);
                  localStorage.setItem("userId", data.id);
                  window.location.href = "/games";
                }
              })
              .catch((error) => {
                console.log("error");
              });
          });
      }
    });
  }
}
customElements.define('login-form', LoginForm)

