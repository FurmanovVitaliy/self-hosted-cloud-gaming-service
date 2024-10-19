import log from "@/common/log";

abstract class Page {// или более конкретный тип
	public title = "Page";
	protected params: Record<string, any> = {};
	protected container: DocumentFragment | HTMLElement;
	private connectedCallbacks: Function [] = [];

	constructor() {
		this.container = document.createElement("div");
		this.container.id = "app";
		this.container.classList.add("container");
	}
	
	protected loadTemplate(template: HTMLElement | Node): void {
		this.container.append(template);
	}

	protected addEventListener(selector: string, event: string, fx: Function, all: boolean = false): void {
		if (!selector) {
			window.addEventListener(event, (e) => fx.call(this, e)); // Правильная привязка контекста 'this'
			return;
		}
		if (all) {
			const elements = this.container.querySelectorAll(selector);
			if (elements.length === 0) {
				log.error(`Элемент не найден: ${selector}`);
				return;
			}
			elements.forEach((element) => 
				element.addEventListener(event, (e) => fx.call(this, e)) // Правильная привязка контекста 'this'
			);
			return;
		}
	
		const element = this.container.querySelector(selector);
		if (!element) {
			log.error(`Элемент не найден: ${selector}`);
			return;
		}
		element.addEventListener(event, (e) => fx.call(this, e)); // Правильная привязка контекста 'this'
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

	public setParams(params: Record<string, any>): void {
		this.params = params;
	}

	async render(params?: Record<string, any>) {
		this.params = params || {};
		return this.container;
	}
}

export default Page;
