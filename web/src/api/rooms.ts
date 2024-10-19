import Api from "./api"
import constants from "@/common/constants"

class RoomsApi extends Api {
    constructor() {
        super(constants.api.server);
    }
    async create(uuid: string, game_id:string, peer?: any) {
        return  this.post(constants.api.roomCreeate, { uuid, game_id, peer });
    }
    async isAlive(uuid: string,username: string) {
        return this.get(constants.api.roomIsAlive(uuid, username));
    }
}
const roomsApi = new RoomsApi();
export default  roomsApi;