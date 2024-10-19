import Api from "./api";
import constants from "@/common/constants";

class GameApi extends Api {
	constructor() {
		super(constants.api.server);
	}
	async getAllGames() {
		return this.get(constants.api.games);
	}

}
const gameApi = new GameApi();

export default gameApi;