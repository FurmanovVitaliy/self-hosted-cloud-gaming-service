import Page from "@comp/templates/page";


class HomePage extends Page {
	public title: string = "Cloud Gaming | Home";
	private html = (
		<null>
			<h1>This is home page </h1>
		</null>
	);

	constructor() {
		super();
		this.insertTemplate(this.html);
	}

	async render() {
		return this.container;
	}
}

export default HomePage;
