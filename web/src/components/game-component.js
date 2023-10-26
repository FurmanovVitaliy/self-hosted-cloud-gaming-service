
class GameComponent extends HTMLElement {
    constructor(){
        super()
        const shadow = this.attachShadow({mode: 'open'})
        this.selected = false;

        const poster = document.createElement('div')
        poster.setAttribute('class', 'game-poster')

        const name = document.createElement('div')
        name.setAttribute('class', 'game-name')
        const button = document.createElement('button')
        button.setAttribute('class', 'create-room-button')
        button.textContent = 'Create Room'
        name.appendChild(button)

      

        const style = document.createElement('style')
        style.textContent = `
        .game-poster{
            border-radius: 20px;
            width: 100%;
	        height: 85%;
	        background-size: cover;
	        position: center;
	        transition: transform 0.5s;
	        will-change: transform;
	
        }
        .game-name{
            color: white;
            text-decoration: none;
            display: block;
            text-align: center;
        }
        .create-room-button{
            background-color: #4CAF50;
            border: none;
            color: white;
            padding: 15px 32px;
            margin: 4px 2px;
            cursor: pointer;
            border-radius: 20px;
        }
      
        `



        shadow.appendChild(style)
        shadow.appendChild(poster)
        shadow.appendChild(name)

    }

    connectedCallback(){
        const shadow = this.shadowRoot;
        const href = this.getAttribute('id')
    }
    
}
customElements.define('game-title', GameComponent )