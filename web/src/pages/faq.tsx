import Page from "@comp/templates/page";
import { enableGamepadNavigation } from "@/components/templates/input";
import { ButtonElementType ,ActionButton } from "@/components/button";

class FaqPage extends Page {
	
	public title: string = "FAQ | Pixel Cloud";

	constructor() {
		super();
		this.loadTemplate(new ActionButton(ButtonElementType.play,"fdsdf").build());
		this.addConnectedCallback(enableGamepadNavigation);
	}

	async render() {
		return this.container;
	}
}

export default FaqPage;
