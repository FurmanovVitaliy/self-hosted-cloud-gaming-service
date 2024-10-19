import Api from "./api";
import constants from "@/common/constants"

class AuthApi extends Api {
    constructor() {
        super(constants.api.server);
    }
     async login(email: string, password: string) {
        return this.post(constants.api.login, { email, password });
    }
     async logout() {
        return this.get(constants.api.logout);
    }
    async signup(email: string, username: string, password: string) {
        return this.post(constants.api.signup, { email, username, password });
    }
}
const authApi = new AuthApi();

export default authApi;