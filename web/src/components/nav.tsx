import { gameElementsSearchCache as searchElements } from "@/common/cache";
import constants from "@/common/constants";


class MainNav extends HTMLElement {
	private html = (
		<null>
			<div class='header__inner'>

				<nav class='header__menu hiden-mobile' gamepad-focus-group>
					{/*}
					<a href={constants.routes.home} class='header__logo logo' gamepad-focus-element  inner-link>
						<img src='/images/logo.png' alt='cloud gaming logo' width='36' height='36' loading='lazy' />
					</a>*/}
					<a href={constants.routes.home} class='header__logo logo' gamepad-focus-element  inner-link>
						LOGO
						</a>
					<ul class='header__menu-list'>
						<li class='header__menu-item'>
							<a href={constants.routes.library} class='header__menu-link' gamepad-focus-element inner-link>
								LIBRARY
							</a>
						</li>
						<li class='header__menu-item'>
							<a href={constants.routes.rooms} class='header__menu-link' gamepad-focus-element inner-link>
								ROOMS
							</a>
						</li>
						<li class='header__menu-item'>
							<a href={constants.routes.faq} class='header__menu-link' gamepad-focus-element  inner-link>
								FAQ
							</a>
						</li>
					</ul>
				</nav>

				<div class="header__search search" gamepad-focus-group>

					<input type="text" class="search__input"
						placeholder="Search" gamepad-focus-element/>

					<ul class="search__result">
						<h2 class="visually-hidden">Search Results</h2>

						{/*
							<li class='game-tile game-tile--search search__item'>
						<a class='game-tile__poster' href='{this.url}' aria-label={this.name}>
								<img src={this.poster} alt={`Poster of ${this.name}`} />
						</a>
						<div class='game-tile__wrapper'>
							<a class='game-tile__header'>
								<h3>{this.name}</h3>
							</a>
							<div class='game-tile__details'>
								<span>{this.platform}</span>
								<span>{this.release}</span>
								<a href={this.url} aria-label={`More about ${this.name} on IGDB`} target='_blank' title={`More about ${this.name} on IGDB`}>
									<span class='accent'>IGDB</span>
								</a>
								<i class="icon--star"></i><span>{this.rating}</span>
							</div>
						</div>
					</li>
							//!--->this is how game.search.element looks like 
							//?--->@/components/game-element.tsx*/}

					</ul>
					<a href="/login" class="header_user-icon user-icon">
						<span>LOG IN</span>
						<img src="/images/user.png" alt="user icon" />
					</a>
				</div>
				<button class='button__burger-menu burger-button visible-mobile' type='button' onclick='mobileOverlay.showModal()'>
					<span class='visually-hidden'>Open navigation menu</span>
				</button>
			</div>

			<dialog class="mobile-overlay visible-mobile" id="mobileOverlay">
				<form class="mobile-overlay__close-button-wrapper" method="dialog">
					<button
						class="mobile-overlay__close-button cross-button"
						type="submit"
					>
						<span class="visually-hidden">Close navigation menu</span>
					</button>
				</form>
				<div class="mobile-overlay__body">
					<ul class="mobile-overlay__list">
						<li class="mobile-overlay__item">
							<a href={constants.routes.home} class='mobile-overlay__link' >
								Home
							</a>
						</li>
						<li class="mobile-overlay__item">
							<a href={constants.routes.library} class='mobile-overlay__link' >
								Library
							</a>
						</li>
						<li class="mobile-overlay__item">
							<a href={constants.routes.rooms} class='mobile-overlay__link'>
								Rooms
							</a>
						</li>
						<li class="mobile-overlay__item">
							<a href='/about' class='mobile-overlay__link'>
								About
							</a>
						</li>
						<li class="mobile-overlay__item">
							<a href={constants.routes.login} class='mobile-overlay__link' >
								Login or SingUp
							</a>
						</li>
					</ul>
				</div>
			</dialog>
		</null>
	)
	private searchInput: HTMLInputElement = this.html.querySelector(".search__input")!;
	private searchResultContainer: HTMLElement = this.html.querySelector(".search__result");
	private searchResultElements: HTMLElement[] | NodeListOf<HTMLAnchorElement> = [];

	constructor() {
		super();
		this.searchResultElements = searchElements.getArr();
		this.searchResultElements.forEach((element) => {
			this.searchResultContainer.appendChild(element.cloneNode(true))
		});

		this.searchInput.addEventListener("input", (e) => {
			const searchValue = this.searchInput.value.toLowerCase();
			const searchResult = this.searchResultContainer.querySelectorAll('.search__item');
			for (let i = 0; i < searchResult.length; i++) {
				this.searchInput.classList.add("has-value");
				const title = searchResult[i].querySelector("h3")!.textContent!.toLowerCase();
				if (title.includes(searchValue)) {
					(searchResult[i] as HTMLElement).style.display = "flex";
					(searchResult[i] as HTMLElement).classList.add("visible");
					this.searchResultContainer.style.display = "flex";
					this.searchResultContainer.style.visibility = "visible";
					this.searchResultContainer.style.opacity = "1";
				} else {
					(searchResult[i] as HTMLElement).style.display = "none";
					(searchResult[i] as HTMLElement).classList.remove("visible");
				}
				if (searchValue === "") {
					(searchResult[i] as HTMLElement).classList.remove("visible");
					this.searchInput.classList.remove("has-value");
					this.searchResultContainer.style.display = "none";
					this.searchResultContainer.style.visibility = "hidden";
					this.searchResultContainer.style.opacity = "0";

				}
			}
		});

	}

	connectedCallback() {
		this.appendChild(this.html);
	}
}

customElements.define("main-nav", MainNav);

export default MainNav;
