import Page from "@comp/templates/page";
import { gameElementsCardsCache as gameCards } from "@/common/cache";
import { enableGamepadNavigation } from "@/components/templates/input";

class LibraryPage extends Page {
	//html template
	private html = (
		<null>
			<header class='header' animated-upd-if="main">
				<main-nav></main-nav>
			</header>
			<main>
				<section class='library' >
					<h1 class='visually-hidden'>Game Library</h1>
					<div class='filters interact'animated-upd="library" static-if-next="library" static-if-prev="library"> 
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
									<span class='custom-select__selected year-select' value=''>
										Select an option
									</span>
								</div>
								<div class='custom-select__options-container '>
									<div class='custom-select__option' value=''>
										All time
									</div>
									<div class='custom-select__option' value='2024'>
										2024
									</div>
									<div class='custom-select__option' value='2023'>
										2023
									</div>
									<div class='custom-select__option' value='2022'>
										2022
									</div>
									<div class='custom-select__option' value='2021'>
										2021
									</div>
									<div class='custom-select__option' value='2020'>
										2020
									</div>
								</div>
							</div>
						</div>
					</div>
					<div class="main-inner" animated-upd="library">
						<div class='sort'>
							<span>Store/Library</span>
							<h2 class='visually-hidden'>Sort</h2>
							<div class='custom-select '>
								<div class='custom-select__select-box custom-select__select-box--wide'>
									<i class='icon--sort'></i>&nbsp;<span>Sorted</span>&nbsp;
									<span class='custom-select__selected sort-select' value=''>
										{" "}
										by default
									</span>
								</div>
								<div class='custom-select__options-container custom-select__options-container--wide'>
									<div class='custom-select__option' value='asc'>
										from A-Z
									</div>
									<div class='custom-select__option' value='desc'>
										from Z-A
									</div>
									<div class='custom-select__option' value='pop'>
										by popularity
									</div>
								</div>
							</div>
						</div>
						<div class='games'>
							{/*
						<div class='game-tile game-tile--lib'>
						<a class='game-tile__poster'>
								<img src={this.poster} alt={`Poster of ${this.name}`} />
						</a>
						<div class='game-tile__wrapper'>
							
							<a class='game-tile__header'>
								<h3>{this.name}</h3>
							</a>
							<span class='game-tile__publisher'>{this.publisher}</span>
							<div class='game-tile__tags details'  data-tags=''>
								{tags}
							</div>
							<div class='game-tile__footer'>
								<span class='game-tile__rating'>
								<i class="icon--star"></i> ({this.rating})
								</span>
								<div class="game-tile__actions">
									<a class="link-button link-button--like" href="#"></a>
								</div>
							</div>								
						</div>
					</div>
							//!--->this is how game.tile.element looks like 
							//?--->@/components/game-element.tsx*/}
						</div>
					</div>
				</section>
			</main>
		</null>
	);
	//declare main fields
	public title: string = "Library | Pixel Cloud";
	private elements_game_cards: HTMLElement[] | NodeListOf<HTMLDivElement>;

	private element_check_genre: NodeListOf<HTMLInputElement>;
	private element_check_rating: NodeListOf<HTMLInputElement>;

	private element_options: NodeListOf<HTMLDivElement>;
	private elements_select_boxes: NodeListOf<HTMLDivElement>;
	private element_selected_year: HTMLSpanElement;
	private element_selected_sort: HTMLSpanElement;

	constructor() {
		super();
		//fill page main fields
		this.elements_game_cards = gameCards.getArr();
		this.element_check_genre = this.html.querySelectorAll<HTMLInputElement>("input[type='checkbox'][name='genre']");
		this.element_check_rating = this.html.querySelectorAll<HTMLInputElement>("input[type='checkbox'][name='rating']");

		this.element_options = this.html.querySelectorAll<HTMLDivElement>(".custom-select__option");
		this.elements_select_boxes = this.html.querySelectorAll<HTMLDivElement>(".custom-select__select-box");
		this.element_selected_year = this.html.querySelector<HTMLSpanElement>(".year-select")!;
		this.element_selected_sort = this.html.querySelector<HTMLSpanElement>(".sort-select")!;

		//add content to page
		this.loadTemplate(this.html);
		this.loadLibrary();
		//local ecents
		this.eventProcessFilters();
		this.eventProcessSort();
		this.eventToggleDropdown();
		this.eventSelectOption();

		//this.selectOption(this.html.querySelector(".option") as HTMLElement);
		//local functions
		//connectedCallback
		this.addConnectedCallback(enableGamepadNavigation);
	}

	//! all appeal to html documents as this.container

	private loadLibrary() {
		this.container.querySelector(".games")!.append(...this.elements_game_cards);
	}

	private eventProcessFilters() {
		this.container.addEventListener("change", () => {
			const selectedGenres = Array.from(this.element_check_genre)
				.filter((checkbox) => checkbox.checked)
				.map((checkbox) => checkbox.value.toLowerCase().trim());

			const selectedRatings = Array.from(this.element_check_rating)
				.filter((checkbox) => checkbox.checked)
				.map((checkbox) => parseFloat(checkbox.value));

			// Изменяем способ получения года
			let selectedYear = this.element_selected_year.getAttribute("value");
			this.elements_game_cards.forEach((tile) => {
				const tileGenres = tile
					.getAttribute("data-genre")!
					.split(",")
					.map((genre) => genre.toLowerCase().trim());
				const tileYear = tile.getAttribute("data-year");
				const tileRating = parseFloat(tile.getAttribute("data-rating") || "0");

				const genreMatch = selectedGenres.length === 0 || selectedGenres.some((genre) => tileGenres.includes(genre));
				const yearMatch = selectedYear === "" || tileYear === selectedYear;

				// Проверка рейтинга в диапазоне
				const ratingMatch =
					selectedRatings.length === 0 ||
					selectedRatings.some(
						(rating) => tileRating >= rating && tileRating < rating + 1, // Отбираем от rating до rating+1
					);

				if (genreMatch && yearMatch && ratingMatch) {
					tile.classList.remove("visually-hidden");
				} else {
					tile.classList.add("visually-hidden");
				}
			});
		});
	}

	private eventProcessSort() {
		this.container.addEventListener("change", () => {
			const selectedSort = this.element_selected_sort.getAttribute("value");
			const games = Array.from(this.elements_game_cards) as HTMLDivElement[];
			if (selectedSort === "asc") {
				games.sort((a, b) => {
					const aTitle = a.getAttribute("data-title")!;
					const bTitle = b.getAttribute("data-title")!;
					return aTitle.localeCompare(bTitle);
				});
			} else if (selectedSort === "desc") {
				games.sort((a, b) => {
					const aTitle = a.getAttribute("data-title")!;
					const bTitle = b.getAttribute("data-title")!;
					return bTitle.localeCompare(aTitle);
				});
			} else if (selectedSort === "pop") {
				games.sort((a, b) => {
					const aRating = parseFloat(a.getAttribute("data-rating")!);
					const bRating = parseFloat(b.getAttribute("data-rating")!);
					return bRating - aRating;
				});
			}
			for (let i = 0; i < games.length; i++) {
				games[i].style.order = i.toString();
			}
		});
	}

	private eventToggleDropdown() {
		this.elements_select_boxes.forEach((element) => {
			element.addEventListener("click", () => {
				const optionsContainer = element.nextElementSibling as HTMLElement;

				// Если элемент скрыт, показываем его с анимацией
				if (!optionsContainer.classList.contains("toggle")) {
					optionsContainer.classList.remove("hidden"); // Убираем скрытие
					requestAnimationFrame(() => {
						optionsContainer.classList.add("toggle");
						element.classList.add("open"); // Добавляем класс для анимации выезжания
					});
				} else {
					// Убираем класс анимации и скрываем после завершения анимации
					optionsContainer.classList.remove("toggle");
					element.classList.remove("open");
					optionsContainer.addEventListener(
						"transitionend",
						() => {
							optionsContainer.classList.add("hidden");
							// Полное скрытие после анимации
						},
						{ once: false },
					);
				}
			});
		});
	}

	private eventSelectOption() {
		this.element_options.forEach((element) => {
			element.addEventListener("click", () => {
				// Удаляем класс 'active' у всех элементов
				this.element_options.forEach((option) => {
					option.classList.remove("custom-select__option--active");
				});

				// Добавляем класс 'active' к выбранному элементу
				element.classList.add("custom-select__option--active");

				// Обновляем текст и атрибут выбранного элемента
				const optionsContainer = element
					.closest(".custom-select")!
					.querySelector(".custom-select__options-container") as HTMLElement;
				const selectBox = element.closest(".custom-select")!.querySelector(".custom-select__select-box") as HTMLElement;
				const selected = element.closest(".custom-select")!.querySelector(".custom-select__selected") as HTMLSpanElement;

				selected.textContent = element.textContent;
				selected.setAttribute("value", element.getAttribute("value")!);

				optionsContainer.classList.remove("toggle");
				selectBox.classList.remove("open");

				// Создаём и триггерим событие 'change'
				const changeEvent = new Event("change", {
					bubbles: true, // позволяет событию всплывать
					cancelable: true, // позволяет отменять событие
				});
				selected.dispatchEvent(changeEvent);
			});
		});
	}

	async render() {
		return this.container;
	}
}

export default LibraryPage;
