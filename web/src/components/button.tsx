import {addItemToLocalStorage} from "@/common/utils";
import {removeItemFromLocalStorage} from "@/common/utils";
import {isItemInLocalStorage} from "@/common/utils";
import eventBus from "@/common/event-bus";
import roomsApi from "@/api/rooms";
import {v4 as uuidv4} from "uuid";

export enum ButtonElementType {
    fav = "favorite",
    play = "play",
}

interface IActionButtonBuilder {
  gameID: string;
	build(): HTMLElement;
}

export class ActionButton implements IActionButtonBuilder {
  public type: ButtonElementType;
  public gameID: string;

  constructor(type: ButtonElementType, gameID: string) {
    this.type = type;
    this.gameID = gameID;
  }
  build(): HTMLElement {
    let el: HTMLElement;
    switch (this.type) {
      case ButtonElementType.play:
        el = (
          <a >
            <span>PLAY NOW</span>
          </a>
        )
        el.addEventListener("click", () => {
        console.log("play button clicked");
        
        });
        break;
      case ButtonElementType.fav:
        el = (
          <a>
            <span>ADD TO LIST</span>
          </a>
        )
    }
    return el;
  }
}