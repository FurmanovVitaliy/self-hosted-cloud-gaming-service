import Page from "@comp/templates/page";


class LibraryPage extends Page {
	private  html = (
		<null>
			<h1>This is library page</h1>
		</null>
			
	);
	public title: string = "Library";

	constructor() {
		super();
		this.insertTemplate(this.html)
	}

	async render() {
		return this.container;
	}
}

export default LibraryPage;
