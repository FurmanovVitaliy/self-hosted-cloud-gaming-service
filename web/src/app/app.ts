import constants from "@/common/constants";
import Route from "route-parser";

import Page from "@/components/templates/page";
import HomePage from "@/pages/home.tsx";
import LoginPage from "@/pages/login.tsx";
import LibraryPage from "@/pages/library";
import ErrorPage from "@/pages/errorPage";


/**
 * Represents a router for handling navigation and rendering of pages.
 */
class App {
	private static ROUTER: App | null = null;
	private routes: { page: any; route: Route }[];
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
			{page: new HomePage(),route: new Route(constants.routes.home)},
			{page: new LibraryPage(),route: new Route(constants.routes.library)},
			{page: new LoginPage(),route: new Route(constants.routes.login)},
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
		//eventBus.execute("nav", this.navigateTo.bind(this));
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
	private getMatchingPage(path: string): Page {
		const targetRoute = this.routes.find((r) => r.route.match(path));
		return targetRoute ? targetRoute.page : this.errorPage;
	}
	/**
	 * Renders the page associated with the given path.
	 * @param path - The path to render.
	 */
	private async render(path: string) {
		const page = this.getMatchingPage(path);
		document.title = page.title;
		this.appElement.innerHTML = "";
		this.appElement.append(await page.render());
		page.callConnectedCallbacks();
	
	}
	/**
	 * Navigates to the specified path.
	 * @param path - The path to navigate to.
	 */
	public navigateTo(path: string) {
		window.history.pushState({ path }, path, path);
		this.render(path);
	}
}

export default App;
