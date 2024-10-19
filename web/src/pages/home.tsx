import Page from "@comp/templates/page";
import log from "@/common/log";
import { v4 } from "uuid";
import roomsApi from "@/api/rooms";
import { enableGamepadNavigation } from "@/components/templates/input";

class HomePage extends Page {
    public title: string = "Pixel Cloud | Home";
    private currentContentIndex: number = 0;
    private contentArray = [
        {
            title: 'Pixel Cloud',
            description: 'Cloud gaming is a type of online gaming that allows direct and on-demand streaming of games onto devices such as computers, consoles.',
            rating: '10.0',
            details: ['3+', 'CLOUD', '2024'],
            videoSrc: '/trailers.webm'
        },
        {
            title: 'Witcher 3: Wild Hunt',
            description: 'This is a new description for the game content.',
            rating: '8.2',
            details: ['16+', 'RPG', 'ACTION' ],
            videoSrc: '/witcher.webm'
        },
        {
            title: 'Last of Us 2',
            description: 'This is a new description for the game content.',
            rating: '7.5',
            details: ['16+', 'HORROR', 'ACTION' ],
            videoSrc: '/last-of-us.webm'
        },
        {
            title: 'Elden Ring',
            description: 'This is a new description for the game content.',
            rating: '7.9',
            details: ['18+', 'RPG' ],
            videoSrc: '/elden-ring.webm'
        }
    ];

    private autoUpdateInterval: number | undefined;

    private html = (
		<null>
        <div class='background'>
        <video autoplay playsinline loop muted class='bg background__video--frame'>
            <source src='/trailers.webm' type='video/webm'></source>
        </video>
        </div>
        <div class='screen-box container'>
            <video autoplay playsinline loop muted class='bg background__video'>
                <source src='/trailers.webm' type='video/webm'></source>
            </video>
            <header class='header' animated-upd="main">
                <main-nav></main-nav>
            </header>
            <main animated-upd="main">
                <section class="hero">
                    <header class="hero__header">
                        <h1 class="hero__title">Cloud Gaming</h1>
                    </header>                    
                    <div class="hero__body">
                        <h2 class="visually-hidden">Game description</h2>
                        <p class="hero__description"></p>
                        <div class="hero__raiting">
                            <h3 class="visually-hidden">Game rating</h3>
                            <span>RAITING </span><span class="hero__raiting--accent"></span>&nbsp;<span>/ 10</span>
                        </div>
                        <div class="hero__details details">
                            <h4 class="visually-hidden">Extra info about the game</h4>
                        </div>
                        <div class="hero__btns" gamepad-focus-group>
                            <a class="link-button link-button--main link-button--play" gamepad-focus-element href="#" data-gameID="">
                                <span class="link-button__text">PLAY NOW</span>
                            </a>
                            <a class="link-button link-button--main  link-button--like" href="#">
                                <span class="link-button__text--like" gamepad-focus-element data-gameID="">ADD TO LIST</span>
                            </a>
                        </div>
                        <div class="hero__pagination pagination"></div> 
                    </div>
                </section>
            </main>
        </div>
		</null>
	);

    constructor() {
        super();
        this.loadTemplate(this.html);
        this.addEventListener(".link-button--play", "click", this.handlePlayButtonClick);
        this.addConnectedCallback(enableGamepadNavigation);
        // Динамически создать кружочки пагинации
        this.createPaginationDots();
        this.updateContent();
        this.startAutoUpdate(); // Начать автоматическое обновление контента
    }

  
    // Метод для создания пагинации (кружочки)
    createPaginationDots() {
        const paginationContainer = this.container.querySelector(".pagination") as HTMLDivElement;
        paginationContainer.innerHTML = ''; // Очистить существующие кружочки
        this.contentArray.forEach((_, index) => {
            const dot = document.createElement("div");
            dot.classList.add("pagination__dot");
            if (index === this.currentContentIndex) {
                dot.classList.add("pagination__dot--active"); // Активный кружочек
            }
            dot.setAttribute("data-index", index.toString());
            dot.addEventListener("click", () => this.handleDotClick(index)); // Обработка клика по кружочку
            paginationContainer.appendChild(dot);
        });
    }

    // Обработчик клика по кружочку
    handleDotClick(index: number) {
        this.currentContentIndex = index;
        this.updateContent();
        this.resetAutoUpdate(); // Сбросить и перезапустить таймер
    }

    // Метод для обновления активного состояния кружочков
    updateActiveDot() {
        const dots = this.container.querySelectorAll(".pagination__dot");
        dots.forEach((dot, index) => {
            if (index === this.currentContentIndex) {
                dot.classList.add("pagination__dot--active");
            } else {
                dot.classList.remove("pagination__dot--active");
            }
        });
    }

    // Метод для обновления контента
    updateContent() {

        const content = this.contentArray[this.currentContentIndex];
        const title = this.container.querySelector(".hero__title") as HTMLHeadingElement;
        const description = this.container.querySelector(".hero__description") as HTMLParagraphElement;
        const rating = this.container.querySelector(".hero__raiting .hero__raiting--accent") as HTMLSpanElement;
        const details = this.container.querySelector(".hero__details") as HTMLDivElement;
        const video = this.container.querySelectorAll(".bg") as NodeListOf<HTMLVideoElement>;

        // Применить эффект исчезновения
        title.classList.add("fade-out");
        description.classList.add("fade-out");
        rating.classList.add("fade-out");
        details.classList.add("fade-out");
        video[0].classList.add("fade-o");
        video[1].classList.add("fade-o");

        setTimeout(() => {
            // Обновить контент
            title.textContent = content.title;
            description.textContent = content.description;
            rating.textContent = content.rating;
            details.innerHTML = "";
            content.details.forEach((detail) => {
                let tag = document.createElement("span");
                tag.classList.add("details__item", "details__item--tag");
                tag.textContent = detail;
                details.appendChild(tag);
            });

            // Обновить видео
            const source = video[0].querySelector("source") as HTMLSourceElement;
            const source2 = video[1].querySelector("source") as HTMLSourceElement;

            source.src = content.videoSrc;
            source2.src = content.videoSrc;
            video[0].load();
            video[1].load();
            video[0].play();
            video[1].play();

            // Удалить эффект исчезновения и применить эффект появления
            title.classList.remove("fade-out");
            description.classList.remove("fade-out");
            rating.classList.remove("fade-out");
            details.classList.remove("fade-out");
            video[0].classList.remove("fade-o");
            video[1].classList.remove("fade-o");

            title.classList.add("fade-in");
            description.classList.add("fade-in");
            rating.classList.add("fade-in");
            details.classList.add("fade-in");
            video[0].classList.add("fade-i");
            video[1].classList.add("fade-i");

            setTimeout(() => {
                title.classList.remove("fade-in");
                description.classList.remove("fade-in");
                rating.classList.remove("fade-in");
                details.classList.remove("fade-in");
                video[0].classList.remove("fade-i");
                video[1].classList.remove("fade-i");
            }, 1000);
        }, 1000);

        // Обновить активный кружочек
        this.updateActiveDot();
        this.currentContentIndex = (this.currentContentIndex + 1) % this.contentArray.length;
    }

    // Начать автообновление контента каждые 30 секунд
    startAutoUpdate() {
        this.autoUpdateInterval = window.setInterval(() => this.updateContent(),  1*60*1000);
    }

    // Сбросить таймер автообновления
    resetAutoUpdate() {
        if (this.autoUpdateInterval) {
            clearInterval(this.autoUpdateInterval);
        }
        this.startAutoUpdate();
    }

    playVideoBackground() {
        const videos = this.container.querySelectorAll(".bg") as NodeListOf<HTMLVideoElement>;
        videos.forEach((video) => {
            if (!video.muted) {
                video.muted = true;
            }
    
            video.play().catch((err) => {
                // Полностью игнорируем ошибку NotAllowedError
                if (err.name !== "NotAllowedError") {
                    console.error(err);
                }
            });
        });
    }

    handlePlayButtonClick(event: Event) {
        const link = event.currentTarget as HTMLElement;
        const gameID = link.getAttribute("data-gameID");
        if (link || gameID) {
            const uuid = v4();
            roomsApi
                .create(uuid, gameID!, {})
                .then((res) => {
                    if (res) {
                        window.location.href = "/rooms/" + uuid;
                    }
                })
                .catch((err) => {
                    log.error(err);
                });
            return;
        }
    }

    async render() {
        this.playVideoBackground();
        return this.container;
    }
}

export default HomePage;
