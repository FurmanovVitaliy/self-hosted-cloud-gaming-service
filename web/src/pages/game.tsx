import Page from "@comp/templates/page";
import { gameElementsFullscreenCache as GameCards } from "@/common/cache";

import { enableGamepadNavigation } from "@/components/templates/input";

class GamePage extends Page {
	//html template
	private html = (
		<null>
			<header class='header' animated-upd-if="main">
				<main-nav></main-nav>
			</header>
			<main>
				<section class='library'>
					<h1 class='visually-hidden'>Game Library</h1>
					<div class='filters not-interact' animated-upd="library" static-if-next="library" static-if-prev="library">
						<h2 class='visually-hidden'>Filters</h2>
						<div class='filters__item'>
							<h3>Categories</h3>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='action' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Action</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='adventure' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Adventure</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='rpg' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>RPG</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='shooter' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Shooter</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='strategy' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Strategy</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='simulation' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Simulation</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='racing' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Racing</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='horror' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Horror</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='fighting' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Fighting</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='genre' value='platform' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>Platform</span>
							</label>
						</div>

						<div class='filters__item'>
							<h3>Rating</h3>
							<label class='custom-checkbox'>
								<input type='checkbox' name='rating' value='5' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
								</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='rating' value='4' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star-gray'></i>
								</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='rating' value='3' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star-gray'></i>
									<i class='icon icon--star-gray'></i>
								</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='rating' value='2' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>
									<i class='icon icon--star'></i>
									<i class='icon icon--star'></i>
									<i class='icon icon--star-gray'></i>
									<i class='icon icon--star-gray'></i>
									<i class='icon icon--star-gray'></i>
								</span>
							</label>
							<label class='custom-checkbox'>
								<input type='checkbox' name='rating' value='1' />
								<span class='custom-checkbox__checkmark'></span>
								<span class='custom-checkbox__field'>
									<i class='icon icon--star'></i>
									<i class='icon icon--star-gray'></i>
									<i class='icon icon--star-gray'></i>
									<i class='icon icon--star-gray'></i>
									<i class='icon icon--star-gray'></i>
								</span>
							</label>
						</div>

						<div class='filters__item'>
							<h3>Year</h3>

							<div class='custom-select'>
								<div class='custom-select__select-box custom-select__select-box--arrow'>
									<span class='custom-select__selected year-select' value="">Select an option</span>
								</div>
								<div class='custom-select__options-container '>
									<div class='custom-select__option' value="">
										All time
									</div>
									<div class='custom-select__option' value="2024">
										2024
									</div>
									<div class='custom-select__option' value="2023">
										2023
									</div>
									<div class='custom-select__option' value="2022">
										2022
									</div>
									<div class='custom-select__option' value="2021">
										2021
									</div>
									<div class='custom-select__option' value="2020">
										2020
									</div>
								</div>
							</div>
						</div>
					</div>
					<div class='game' animated-upd="library" >
						{/* 
						//! Game card template 	
						<div class='game-card game-card--fullscreen'>
						<h2 class='visually-hidden'>Game card</h2>
						<a class="link-button game-card__go-to" amepad-focus-element href='/library' inner-link>
							<span><i class="icon icon--go-left"></i></span>
						</a>
						<div class='game-card__banner'>
							<img src={bannerImgUrl} />
						</div>
						<div class='game-card__wrapper'>
							<div class='game-card__poster'>
								<img src={this.poster} alt={`Poster of ${this.name}`} />
							</div>
							<div class='game-card__wrapper game-card__wrapper--column'>
								<div class='game-card__rating'>
									<span>RAITING </span>
									<span class='accent'>{(this.rating / 20).toFixed(1)}</span>
									&nbsp;
									<span>/ 10</span>
								</div>
								<div class='game-card__wrapper game-card__wrapper--column'>
									<div class='game-card__details details'>
										{ageTag}
										{genresTags}
									</div>
									<div class='game-card__header' href={constants.routes.library + "/" + this.id} inner-link>
										<h3>{this.name}</h3>
									</div>
									<div class='game-card__summary'>
										<p>{this.summary}</p>
									</div>
								</div>
								<div class='game-card__actions'>
									<a class='link-button link-button--main link-button--play' gamepad-focus-element data-gameID={this.id}>
										<span class='link-button__text'>PLAY NOW</span>
									</a>
									<a class='link-button link-button--main link-button--like' href='#'>
										<span class='link-button__text--like' gamepad-focus-element data-gameID=''>
											ADD TO LIST
										</span>
									</a>
								</div>
							</div>
						</div>
						{media}
					</div>*/}
					</div>
				</section>
			</main>
		</null>
	);
	//declare main fields
	public title: string = "Game | Pixel Cloud";

	constructor() {
		super();
		//fill page main fields
		//add content to page
		this.loadTemplate(this.html);
		//add event listeners
		//local functions
		//connectedCallback
		this.addConnectedCallback(enableGamepadNavigation);
	}

	//! all appeal to html documents as this.container

	private insertGameCard() {
		const uuid = this.params.uuid;
		const gameCard = GameCards.get(uuid);
		if (!gameCard) return;
		const gameContainer = this.container.querySelector(".game");
		if (gameContainer) {
			gameContainer.innerHTML = "";
			gameContainer.appendChild(gameCard);
		} else {
			console.error("Game container not found");
		}
	}

	async render() {
		this.insertGameCard();
		return this.container;
	}
}

export default GamePage;
