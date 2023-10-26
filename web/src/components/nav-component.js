import appConstants from '../common/constants'
class NavComponent extends HTMLElement {
    constructor(){
        super()
        const shadow = this.attachShadow({mode: 'open'})
        const wrapper = document.createElement('div')
        this.searchType = appConstants.search.types.post

        wrapper.setAttribute('class', 'main-menu')
        this.links = [
            {href: appConstants.routes.index, name: 'Home', class: 'home-link'},
            {href: appConstants.routes.games, name: 'Games', class: 'games-link'},
            {href: appConstants.routes.rooms, name: 'Rooms', class: 'rooms-link'},
        ]

        const style = document.createElement('style')

        style.textContent = `
           .main-menu {
                z-index: 2;
                position: sticky;
                display: flex;
               align-items: center;
               padding: 4px 4px;
            

           }

           .global-search {
               border: 1px solid #aaa;
              
              
               max-haight: 10vh;
               width: 100wv;
            
           }

           .global-search:placeholder{
               color: #aaa;
           }
           
        `

        shadow.appendChild(style)
        shadow.appendChild(wrapper)

        this.links.forEach(link => {
            const l = document.createElement('nav-link')
            l.setAttribute('class', `main-link ${link.class}`)
            l.setAttribute('href', link.href)
            l.setAttribute('text', link.name)
            wrapper.appendChild(l)
        })

        const search = document.createElement('input')
        search.setAttribute('class', 'global-search')
        search.addEventListener('keyup', (e) => {
            e.stopPropagation()
            if(e.key === 'Enter') {
                e.preventDefault()
                const text = e.target.value
                console.log('search', text)
            }
        })

        wrapper.appendChild(search)

    }

    updateSearch() {
        const shadow = this.shadowRoot
        const input = shadow.querySelector('input')
        const search = this.getAttribute('search')
        input.value = search
        if(this.searchType === appConstants.search.types.post){
            input.setAttribute('placeholder', 'Search post...')
        } else if(this.searchType === appConstants.search.types.user){
            input.setAttribute('placeholder', 'Search user...')
        }
        
    }

    connectedCallback(){
        const shadow = this.shadowRoot;
        const searchText = this.getAttribute('search')
        this.searchType = this.getAttribute('type') ? this.getAttribute('type') : appConstants.search.types.post

        if(searchText){
            const input = shadow.querySelector('input')
            input.value = searchText
        }

        const {pathname: path} = new URL(window.location.href)
        const link = this.links.find((l) => l.href === path)

        if(link) {
            const linkElement = shadow.querySelector(`.${link.class}`)
            linkElement.setAttribute('selected', 'true')
        }
    }

    
    static get observedAttributes(){
        return ['search', 'type']
    }

    attributeChangedCallback(name, oldValue, newValue){
        if(name === 'search'){
            this.updateSearch()
        }
        if(name === 'type'){
            this.searchType = newValue
            this.updateSearch()
        }
    }
}

customElements.define('main-nav', NavComponent)