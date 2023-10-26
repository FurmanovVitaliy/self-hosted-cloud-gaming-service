import {getRooms} from '../api/roomsApi'

class RoomsComponent extends HTMLElement {
    constructor(){
        super()

        const shadow = this.attachShadow({mode: 'open'})
        const wrapper = document.createElement('div')
        wrapper.setAttribute('class', 'rooms-block')

        const style = document.createElement('style')
        
        style.textContent = `
        
        .rooms-block{
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            width: 100%;
        }
        `
        this.selected = false;
        shadow.appendChild(style)
        shadow.appendChild(wrapper)
        
    }
    
    connectedCallback(){
        this.getRooms();
    }

    async getRooms(){
        const roomList = this.shadowRoot.querySelector('.rooms-block');
        
       
       getRooms().then((rooms)=>{
       rooms.forEach((room) => {
            console.log(room);
            const roomComponent = document.createElement('room-component');
            const title = roomComponent.shadowRoot.querySelector('.room-title');
            title.textContent = room.name;
            const id = roomComponent.shadowRoot.querySelector('.room-id');
            id.textContent = room.id;
            roomList.appendChild(roomComponent);
        })
       });
    }
}

customElements.define('rooms-list', RoomsComponent)