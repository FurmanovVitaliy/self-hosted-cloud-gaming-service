import log from "@/common/log";

abstract class Page {
	protected container: DocumentFragment | HTMLElement;
	title = "Page";
	private connectedCallbacks: Function [] = [];

	constructor() {
		this.container = document.createElement("div");
		this.container.id = "app";
		this.container.classList.add("container");
	}
	
	protected insertTemplate(template: HTMLElement | Node): void {
		this.container.append(template);
	}

	protected addEventListener(selector: string, event: string, fx: Function, all:boolean=false): void {
		if (all) {
			const elements = this.container.querySelectorAll(selector);
			if (elements.length === 0) {
				log.error(`Element not found: ${selector}`);
				return;
			}
			elements.forEach((element) => element.addEventListener(event, (e) => fx(e).bind(this)));
			return;
		}
		const element = this.container.querySelector(selector);
		if (!element) {
			log.error(`Element not found: ${selector}`);
			return;
		}
		element.addEventListener(event, (e) => fx(e));
	}

	protected addConnectedCallback(callback: Function): void {
		this.connectedCallbacks.push(callback);
	}

	callConnectedCallbacks(): void {
		if (this.connectedCallbacks.length === 0) {
			return;
		}
		this.connectedCallbacks.forEach((callback) => callback());
	}


	async render() {
		return this.container;
	}
}

export default Page;
