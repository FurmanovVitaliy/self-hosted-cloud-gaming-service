import { log } from "../common/log.js";
import { authApi } from "../api";

class FakeForm extends HTMLElement {
  constructor() {

    super();
    const shadow = this.attachShadow({ mode: 'open' });
    const wrapper = document.createElement('div')
    wrapper.setAttribute('class', 'form-block');
    wrapper.innerHTML = `
    <div class="form">
      <h4 class="title"></h4>
      <input placeholder="User name" class="input-text">
      <input type="password" placeholder="password" class="input-password">
      <div class="buttons">
          <button class='ok-button'>Ok</button>
          <button class='cancel-button'>Cancel</button>
      </div>
    </div>
    `


    const style = document.createElement('style');
    style.textContent = `
    .wrapper{
      width: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 20px;
      text-align: center;
    }
    .form{
      display: flex;
      justify-content: center;
      align-items: center;
      flex-direction: column;
      padding: 10px;
      border: 1px solid #ccc;
      
      border-radius: 8px;
      margin: 10px;
      font-size: 10px;
      font-family: arial;
    }

    h4{
        text-align: center;
    }
    input{
      margin-bottom: 5px;
    }
    `

    shadow.appendChild(style);
    shadow.appendChild(wrapper);

  }

  connectedCallback() {

    const shadow = this.shadowRoot;
  

    const wrapper = shadow.querySelector('.form-block')
    const title = shadow.querySelector('.title')

    const okButton = shadow.querySelector('.ok-button')
    const cancelButton = shadow.querySelector('.cancel-button')

    const inputText = shadow.querySelector('.input-text')
    const inputPassword = shadow.querySelector('.input-password')

    okButton.addEventListener('click', (e) => {
 
        if (inputText && inputPassword){
            authApi.login(inputText.value,inputPassword.value).then((response) => {
              console.log(response.headers)
              //(response.headers.get('Set-Cookie')) ? log.debug("Coockie present in response") : log.debug("Coockie not present in response")
                //console.bebug(response)
                response.json().then((data) => {
                    console.log(data)
                    if (data){
                        console.log("Auth OK")
                        //localStorage.setItem('token', data.token)
                        //localStorage.setItem('user', data.user)
                        window.location.href = '/games'
                    }
                })
            })
        }

    })

    cancelButton.addEventListener('click', (e) => {
        console.log("Auth Canceled")

    })
    }

authWithLoginAndPassword(){
    login = shadow.querySelector('.input-text');
    password = shadow.querySelector('.input-password');
    if (login && password){
        authApi.login(login,password).then((response) => {
          
            console.log(response)
        })
    }
}
}
customElements.define('login-form', FakeForm)

