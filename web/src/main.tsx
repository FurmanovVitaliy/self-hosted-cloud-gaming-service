import App from "@/app/app";
import { initializeGameDataCache, initializeGameElemsSearchCache, initializeGameElemsTileCache, initializeGameElemsFullscreenCache} from "@/common/cache";
//component imports

import "@comp/nav";


import { simpleHandling } from "./components/templates/input";
import { startUpdatingGamepadState } from "./components/templates/input";

//scss imports
import "@css/styles.scss";
//external imports
import "https://unpkg.com/ionicons@7.1.0/dist/ionicons/ionicons.esm.js";


await initializeGameDataCache();
initializeGameElemsTileCache();
initializeGameElemsSearchCache();
initializeGameElemsFullscreenCache();





App.router();
simpleHandling();






