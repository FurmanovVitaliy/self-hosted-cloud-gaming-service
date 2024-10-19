import Route from "route-parser";
import { clearFocus } from "@/components/templates/input";
import constants from "@/common/constants";
import eventBus from "@/common/event-bus";

import Page from "@/components/templates/page";
import HomePage from "@/pages/home.tsx";
import RoomPage from "@/pages/room.tsx";
import RoomsPage from "@/pages/rooms.tsx";
import LoginPage from "@/pages/login.tsx";
import SingUpPage from "@/pages/singup.tsx";
import LibraryPage from "@/pages/library";
import GamePage from "@/pages/game";
import FaqPage from "@/pages/faq";
import ErrorPage from "@/pages/errorPage";

/**
 * Represents a router for handling navigation and rendering of pages.
 */
class App {
	private static ROUTER: App | null = null;
	private routes: { page: any; route: Route; params: any }[];
	private appElement: HTMLElement;
	private errorPage: Page;

	/**
	 * Constructs a new instance of the Router class.
	 * The constructor is private to enforce the singleton pattern.
	 * New pages and routes of app should be added here.
	 */
	private constructor() {
		this.appElement = document.body!;
		this.errorPage = new ErrorPage();
		this.routes = [
			{ page: new HomePage(), route: new Route(constants.routes.home), params: {} },
			{ page: new LibraryPage(), route: new Route(constants.routes.library), params: {} },
			{ page: new RoomsPage(), route: new Route(constants.routes.rooms), params: {} },
			{ page: new RoomPage(), route: new Route(constants.routes.room), params: {} },
			{ page: new LoginPage(), route: new Route(constants.routes.login), params: {} },
			{ page: new SingUpPage(), route: new Route(constants.routes.singup), params: {} },
			{ page: new GamePage(), route: new Route(constants.routes.game), params: {} },
			{ page: new FaqPage(), route: new Route(constants.routes.faq), params: {} },
		];
	}
	/**
	 * Initializes the router by setting up event listeners and rendering the initial page.
	 */
	private initRouter() {
		/**
		 * eventBus represent publisher / subscriber pattern,
		 * as way to brake cirkle dependency between router and pages,
		 * for exemple if you need to call  NavigateTo from page
		 *
		 */
		eventBus.execute("nav", this.navigateTo.bind(this));
		window.addEventListener("popstate", () => this.render(new URL(window.location.href).pathname));
		this.appElement.addEventListener("click", (e) => {
			const target = e.target as HTMLElement;
			const link = target.closest("[inner-link]") as HTMLAnchorElement;
			if (link) {
				e.preventDefault();
				this.navigateTo(new URL(link.href).pathname);
			}
		});
		this.render(new URL(window.location.href).pathname);
	}
	/**
	 * Returns the singleton instance of the Router class.
	 * If an instance does not exist, it creates a new one.
	 *
	 * @returns The Router instance.
	 */
	static router() {
		if (!App.ROUTER) {
			App.ROUTER = new App();
			App.ROUTER.initRouter();
		}
		return App.ROUTER;
	}
	/**
	 * Gets the route that matches the given path.
	 * @returns The matched page  or prepared errorpage if no match is found.
	 */
	private getMatchingPage(path: string): { page: Page; params?: Record<string, any> } {
		const targetRoute = this.routes.find((r) => r.route.match(path));
		let params: Record<string, any> | undefined;

		if (targetRoute) {
			// Предполагается, что match возвращает объект с параметрами
			const matchResult = targetRoute.route.match(path);
			if (matchResult !== false) {
				params = matchResult;
				targetRoute.params = params; // Сохраняем параметры в целевом маршруте
			} else {
				params = undefined;
			}
		}

		return {
			page: targetRoute ? targetRoute.page : this.errorPage,
			params, // Возвращаем параметры
		};
	}
	/**
	 * Renders the page associated with the given path.
	 * @param path - The path to render.
	 */
	private async render(path: string) {
		const { page, params } = this.getMatchingPage(path);
		page.setParams(params || {}); // Устанавливаем параметры перед вызовом render
		document.title = page.title;

		// Анимации и дальнейшая логика...

		const ex= (path.match(/\/([^/]+)(?:\/|$)/) || [])[1] || "main";
		const elementsToAnimate = this.appElement.querySelectorAll(`[animated-upd]:not([static-if-next="${ex}"]), [animated-upd-if="${ex}"]`) as NodeListOf<HTMLElement>;
		this.appElement.querySelectorAll(`[animated-upd]:not([static-if-next="${ex}"])`)
		const prev = elementsToAnimate[0]?.getAttribute("animated-upd") || "?";
		if (elementsToAnimate.length > 0) {
			await this.animateOut(elementsToAnimate);
		} else {
			this.appElement.innerHTML = "";
		}

		this.appElement.append(await page.render()); // Здесь рендерится страница с доступными параметрами
		page.callConnectedCallbacks();

		// Анимация появления новых элементов
		const newElementsToAnimate = this.appElement.querySelectorAll(`[animated-upd]:not([static-if-prev="${prev}"]), [animated-if-prev="${prev}"], [animated-upd-if="${prev}"]`) as NodeListOf<HTMLElement>;
		if (newElementsToAnimate.length > 0) {
			this.animateIn(newElementsToAnimate);
		}
	}

	/**
	 * Animates the fade-out effect for the given elements.
	 * @param elements - The elements to animate.
	 * @returns A promise that resolves when the fade-out animation is complete.
	 */
	private async animateOut(elements: NodeListOf<HTMLElement>): Promise<void> {
		return new Promise((resolve) => {
			elements.forEach((element) => {
				element.classList.add("fade-out");
				element.addEventListener(
					"animationend",
					() => {
						element.classList.remove("fade-out"); // Remove the fade-out class
						// Check if all animations of the elements are complete
						if ([...elements].every((el) => !el.classList.contains("fade-out"))) {
							this.appElement.innerHTML = ""; // Clear the content after all animations are complete
							resolve(); // Resolve the promise
						}
					},
					{ once: true },
				);
			});
		});
	}

	/**
	 * Animates the fade-in effect for the given elements.
	 * @param elements - The elements to animate.
	 */
	private animateIn(elements: NodeListOf<HTMLElement>): void {
		elements.forEach((element) => {
			element.classList.add("fade-in");

			// Wait for the fade-in animation to complete
			element.addEventListener(
				"animationend",
				() => {
					element.classList.remove("fade-in"); // Remove the fade-in class
				},
				{ once: true },
			);
		});
	}

	/**
	 * Navigates to the specified path.
	 * @param path - The path to navigate to.
	 */
	public navigateTo(path: string) {
		clearFocus(); //disable gamepad update and navigation
		window.history.pushState({ path }, path, path);
		this.render(path);
	}
}

export default App;
