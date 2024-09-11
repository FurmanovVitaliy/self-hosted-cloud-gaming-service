import { jsx } from "jsx-runtime";

declare global {
	module JSX {
		type IntrinsicElements = Record<keyof HTMLElementTagNameMap, Record<string, any>> & {
			"null": Record<string, any>;
			"game-slider": Record<string, any>;
			"room-display": Record<string, any>;
			"main-nav": Record<string, any>;
			"ion-icon": Record<string, any>;
			"main": Record<string, any>;
			"input": Record<string, any>;
			"path": Record<string, any>;
			"svg": Record<string, any>;
			"g": Record<string, any>;
		};
	}
}
