import {navigateTo, routes} from '../router/'

class RoomComponent extends HTMLElement {
    constructor(){
        super()

        const shadow = this.attachShadow({mode: 'open'})
        const wrapper = document.createElement('div')
        wrapper.setAttribute('class', 'room-block')
        wrapper.innerHTML = `
        <div class="room-title"></div>
        <div class="room-description"></div>
        <div class="room-players"></div>
        <div class="room-uuid"></div>
        `
        const style = document.createElement('style')
        
        style.textContent = `
        
        .room-block{
            height: 15vh;
            width: 20vw;
            display: flex;
            justify-content: center;
            align-items: center;
            flex-direction: column;
            padding: 10px;
            background-color: #020508;
            border: 1px solid #ccc;
            border-radius: 8px;
            margin: 10px;
        }
        `
        this.selected = false;
        shadow.appendChild(style)
        shadow.appendChild(wrapper)
        
    }
    
    connectedCallback(){
        const shadow = this.shadowRoot;
        const id = shadow.querySelector('.room-uuid');
        id.addEventListener('click', (e) => {navigateTo(routes.Room.reverse({uuid: e.target.textContent})); console.log(e.target.textContent)})
    }
}

customElements.define('room-component', RoomComponent)