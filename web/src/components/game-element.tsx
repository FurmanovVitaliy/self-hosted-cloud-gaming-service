import { Game } from "@/types/types";
import constants from "@/common/constants";
import { isUserLoggedIn } from "@/common/utils";
import roomsApi from "@/api/rooms";
import ResponseError from "@/api/err";
import eventBus from "@/common/event-bus";
import { v4 as uuidv4 } from "uuid";

type createRoomResponce = { uuid: string; game_name: string};

export enum GameElementType {
	fullscreen = "img",
	search = "search",
	tile = "slide",
}

interface IGameElementBuilder {
	addData(game: Game): this;
	addFx(fx: Function): this;
	build(): HTMLElement;
}
export class GameElementBuilder implements IGameElementBuilder {
	private type: GameElementType;
	private id: string = "";
	private name: string = "";
	private url: string = "";
	private poster: string = "";
	private rating: number = 0;
	private release: number = 0;
	private ageRating: string = "";
	private summary: string = "";
	private videos: string[] = [];
	private images: string[] = [];
	private developer: string = "";
	private platform: string = "";
	private fx: Function[] = [];
	private publisher: string = "";
	private tags: string[] = [];

	constructor(type: GameElementType) {
		this.type = type;
	}

	addData(game: Game): this {
		this.id = game.id;
		this.name = game.name;
		this.url = game.url || "";
		this.poster = game.poster || "";
		this.rating = game.rating ? Math.round(game.rating * 10) / 10 : 0;
		this.release = game.release || 0;
		this.ageRating = game.age_rating || "";
		this.summary = game.summary || "";
		this.videos = game.videos || [];
		this.images = game.images || [];
		this.platform = game.platform || "";
		this.developer = game.developer || "";
		this.publisher = game.publisher || "";
		this.tags = game.genres || [];
		return this;
	}

	addFx(fx: Function): this {
		this.fx.push(fx);
		return this;
	}

	build(): HTMLElement {
		let el: HTMLElement;
		let media = document.createElement("div");
		media.classList.add("game-card__media");
		let trailers = this.videos.forEach((videoUrl, index) => {
			const embedUrl = videoUrl.replace(
				/https:\/\/www\.youtube\.com\/watch\?v=([a-zA-Z0-9_-]+)/,
				"https://www.youtube.com/embed/$1",
			);

			const iframe = document.createElement("iframe");
			iframe.width = "100%";
			iframe.height = "auto";
			iframe.src = embedUrl;
			iframe.frameBorder = "0";
			iframe.title = `${this.name} trailer-video ${index + 1}`;
			iframe.allow = "accelerometer; clipboard-write; encrypted-media; gyroscope; picture-in-picture";
			iframe.allowFullscreen = true;
			media.appendChild(iframe);
		});
		let bannerImgUrl = this.images[Math.floor(Math.random() * this.images.length)].replace("/t_thumb/", "/t_1080p/");
		let genresTags = this.tags.map((detail) => {
			const tag = document.createElement("span");
			tag.classList.add("details__item", "details__item--tag");
			tag.textContent =
				detail === "Role-playing (RPG)" ? "RPG" : detail.charAt(0).toUpperCase() + detail.slice(1).toLowerCase();
			return tag;
		});
		let ageTag = Object.assign(document.createElement("span"), {
			className: "details__item details__item--tag",
			textContent: this.ageRating,
		});

		switch (this.type) {
			case GameElementType.fullscreen:
				el = (
					<div class='game-card game-card--fullscreen'>
						<h2 class='visually-hidden'>Game card</h2>
						<a class='link-button game-card__go-to' amepad-focus-element href='/library' inner-link>
							<span>
								<i class='icon icon--go-left'></i>
							</span>
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
					</div>
				);
				const playButton = el.querySelector(".link-button--play");
				if (playButton) {
					playButton.addEventListener("click", async () => {
						if (isUserLoggedIn()) {
							const id = playButton.getAttribute("data-gameID")
							if (id) {
								try{
								const responseData = (await roomsApi.create(uuidv4(),id)) as createRoomResponce
								 history.pushState(null, "", "/rooms/" + responseData.uuid);
								 window.location.assign("/rooms/" + responseData.uuid);
								//eventBus.notify("nav", "/rooms/" + responseData.uuid);
								}
								catch(error){
									if (error instanceof ResponseError) {
										alert(error.details);
									}
								}
							}
						} else {
							eventBus.notify("nav", constants.routes.login);
						}
					});
				}

				break;
			case GameElementType.search:
				el = (
					<li class='game-card game-card--search search__item'>
						<a class='game-card__poster' href={constants.routes.library + "/" + this.id} inner-link aria-label={this.name}>
							<img src={this.poster} alt={`Poster of ${this.name}`} />
						</a>
						<div class='game-card__wrapper game-card__wrapper--column'>
							<a class='game-card__header game-card__header--nowrap' href={constants.routes.library + "/" + this.id} inner-link>
								<h3>{this.name}</h3>
							</a>
							<div class='game-card__details'>
								<span>{this.platform}</span>
								<span>{this.release}</span>
								<a
									href={this.url}
									aria-label={`More about ${this.name} on IGDB`}
									target='_blank'
									title={`More about ${this.name} on IGDB`}
								>
									<span class='accent'>IGDB</span>
								</a>
								<i class='icon--star'></i>
								<span>{this.rating}</span>
							</div>
						</div>
					</li>
				);
				break;
			case GameElementType.tile:
				el = (
					<div class='game-card game-card--lib' data-genre={this.tags.join(",")}>
						<a class='game-card__poster' href={constants.routes.library + "/" + this.id} inner-link>
							<img src={this.poster} alt={`Poster of ${this.name}`} />
						</a>
						<div class='game-card__wrapper game-card__wrapper--column'>
							<a
								class='game-card__header game-card__header--fixed-height'
								href={`${constants.routes.library}/uuid:${this.id}`}
								inner-link
							>
								<h3>{this.name}</h3>
								<p>{this.publisher}</p>
							</a>
							<div class='game-card__details details'>
								{ageTag}
								{genresTags}
							</div>
							<div class='game-card__wrapper game-card__wrapper--end'>
								<span>
									<i class='icon--star'></i>({(this.rating / 20).toFixed(1)})
								</span>
								<div class='game-card__actions'>
									<a class='link-button link-button--like' href='#'></a>
								</div>
							</div>
						</div>
					</div>
				);

				break;
		}
		el.setAttribute("data-id", this.id);
		el.setAttribute("data-year", this.release.toString());
		el.setAttribute("data-title", this.name);
		el.setAttribute("data-rating", (this.rating / 20).toFixed(2).toString());

		return el;
	}
}
