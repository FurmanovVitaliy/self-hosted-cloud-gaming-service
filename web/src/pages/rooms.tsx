import constants from "@/common/constants";
import Page from "@/components/templates/page";


class RoomPage extends Page {
  title: string = "Rooms | Pixel Cloud";

  private html = (
    <main-nav></main-nav>
  )
  constructor() {
    super();
    
    this.loadTemplate(this.html);
  }
  async render() {
    return this.container;
  }
} export default RoomPage;