import { MessageJSON, DeviceInfoJSON } from "@/types/types";

//add favorite game to local storage ad array

export const addItemToLocalStorage = (key: string, item: any) => {
    let items = JSON.parse(localStorage.getItem(key) || "[]");
    items.push(item);
    localStorage.setItem(key, JSON.stringify(items));
}

export const removeItemFromLocalStorage = (key: string, item: any) => {
    let items = JSON.parse(localStorage.getItem(key) || "[]");
    items = items.filter((i: any) => i !== item);
    localStorage.setItem(key, JSON.stringify(items));
}

export const isItemInLocalStorage = (key: string, item: any) => {
    let items = JSON.parse(localStorage.getItem(key) || "[]");
    return items.includes(item);
}

// chek is user logged in
export const isUserLoggedIn = () => {
    return localStorage.getItem("access_token") ? true : false; 
}

/**
 * Constructs a message object with the specified content type, content, UUID, and username.
 * If the content is an object, it is encoded using base64 before being added to the message object.
 * 
 * @param tag - The tag associated with the message.
 * @param contentType - The type of the message content.
 * @param content - The content of the message.
 * @param UUID - The UUID associated with the message (optional).
 * @param username - The username associated with the message (optional).
 * @returns A JSON string representation of the message object.
 */

export const msg = (tag : string ,contentType: string,
    content: any | undefined = undefined,
    UUID: string | undefined = undefined,
    username: string | undefined = undefined): string => {
    const message: MessageJSON = {
        tag: tag,
        content_type: contentType,
        content: content,
        UUID: UUID,
        username: username,
    };


    return JSON.stringify(message);
};


interface ClientDeviceInfo {
    display?: {
        height?: number;
        width?: number;
    };
    control?: {
        type?: string;
        vendorID?: string;
        productID?: string;
    };
}

function getGamepadInfo(): { type: string; vendorID?: string; productID?: string } | null {
    const gamepads = navigator.getGamepads();
    for (const gamepad of gamepads) {
        if (gamepad) {
            const id = gamepad.id;
            const vendorProductMatch = id.match(/Vendor:\s*(\w+)\s*Product:\s*(\w+)/i);

            let vendorID, productID;

            if (vendorProductMatch) {
                vendorID = vendorProductMatch[1];
                productID = vendorProductMatch[2];
            } else if (id.includes('-')) {
                [vendorID, productID] = id.split('-');
            }

            return {
                type: 'gamepad',
                vendorID,
                productID
            };
        }
    }
    return null;
}


function getTouchInfo(): {  type: string; vendorID?: string; productID?: string } | null {
    if ('ontouchstart' in window || navigator.maxTouchPoints > 0) {
        return { type: 'touch', 
        vendorID: "default",
        productID: "default"    
    };
    }
    return null;
}

function getKeyboardInfo(): {  type: string; vendorID?: string; productID?: string } {
    return { type: 'keyboard' ,
    vendorID: "default",
    productID: "default"                    
};
}

export function getClientDeviceInfo(): ClientDeviceInfo {
    const displayInfo = {
        height: window.screen.height,
        width: window.screen.width
    };

    let controlInfo = getGamepadInfo();

    if (!controlInfo) {
        controlInfo = getTouchInfo();
    }
    if (!controlInfo) {
        controlInfo = getKeyboardInfo();
    }
    return {
        display: displayInfo,
        control: controlInfo
    };
}


