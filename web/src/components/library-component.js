import {log} from '../common/log.js';
import Swiper from 'swiper';
import 'swiper/css';
import {Keyboard, Parallax,Controller} from 'swiper/modules'
import { getGames } from '../api/gamesApi';





class LibraryComponent extends HTMLElement {
    constructor(){
        super()
    
        const sliderMain = document.createElement('div')
        sliderMain.setAttribute('class', 'slider_main slider')
        const sliderBack = document.createElement('div')
        sliderBack.setAttribute('class', 'slider_bg slider')

        sliderMain.innerHTML = `<div class="swiper-wrapper slider__wrapper"></div>`
        sliderBack.innerHTML = `<div class="swiper-wrapper slider__wrapper"></div>`
        

        const description = document.createElement('div')
        description.innerHTML =`
        <div class="description">
        <div class="logo">Games</div>
        <p>Lorem ipsum dolor  maxime ipsa temporibus obcaecati totam recusandae illo minus. Placeat nulla iusto illum nemo voluptas doloremque ab similique cupiditate rerum.</p>
        </div>`
        const link = document.createElement('link');
        link.setAttribute('rel', 'stylesheet');
        link.setAttribute('href', '/src/styles/library.css');
        
        this.appendChild(link)
        this.appendChild(description)
        this.appendChild(sliderMain)
        this.appendChild(sliderBack)
              
    }
    
    connectedCallback() {
            this.initSwiper();
            log.info('LibraryComponent connected');
            this.addSlides();
        
        }
        
        
       
        initSwiper() {
           
        const sliderBg = new Swiper(document.querySelector('.slider_bg'), {
            modules: [Parallax],
            centeredSlides: true,
            parallax: true,
            spaceBetween: 60,
            slidesPerView: 3.5
        })
        const sliderMain = new Swiper(document.querySelector('.slider_main'), {
            modules: [Keyboard,Parallax,Controller],
            controller: {
                control: sliderBg,
            },
            parallax: true,
            centeredSlides: true,
            mousewheel: true,
            keyboard:{enabled: true,},
            breakpoints: {
                0: {
                    slidesPerView: 2,
                    spaceBetween: 5
                },
                680: {
                    slidesPerView: 3.5,
                    spaceBetween: 10
                }
            }
        })
       

    }

    addSlides() {
        const slider = document.querySelector('.slider_main .swiper-wrapper');
        const sliderBack = document.querySelector('.slider_bg .swiper-wrapper');
        getGames().then((games) => {
            games.forEach((game) => {
                const slide = document.createElement('game-component');
                slide.setAttribute('class', 'swiper-slide slide')
                const poster = slide.shadowRoot.querySelector('.game-poster');
                poster.setAttribute('style', `background-image: url(${game.logo})`);
                const name = slide.shadowRoot.querySelector('.game-name');
                name.appendChild(document.createTextNode(game.name));
                name.setAttribute('game-id', game.id);
              
             

                slider.appendChild(slide);  
            });
        });
        getGames().then((games) => {
            games.forEach((game) => {
                const slide = document.createElement('div');
                slide.setAttribute('class', 'swiper-slide slide');
                const poster = slide.appendChild(document.createElement('div'))
                poster.setAttribute('data-swiper-parallax', '100');
                poster.setAttribute('class', 'game-poster');
                poster.setAttribute('style', `background-image: url(${game.logo})`);

                sliderBack.appendChild(slide);
        

            });
        });
        
    }
}
customElements.define('game-library', LibraryComponent)