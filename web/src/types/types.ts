export type Game = {
	id: string; // Unique identifier for the game
	name: string; // Name of the game
	url?: string; // URL to the game's page (optional)
	poster?: string; // URL of the game's poster (optional)
	platform?: string; // Platform on which the game is available (optional)
	rating?: number; // Rating of the game (optional)
	summary?: string; // Brief description of the game (optional)
	videos?: string[]; // List of links to video materials (optional)
	release?: number; // Release year of the game (optional)
	images?: string[]; // List of links to images of the game (optional)
	age_rating?: string; // Age rating of the game (optional)
	developer?: string; // Name of the game developer (optional)
	publisher?: string; // Name of the game publisher (optional)
	genres?: string[]; // List of genres of the game (optional)
};

export type MessageJSON = {
	tag: string;
	content_type: string;
	content?: any;
	UUID?: string;
	username?: string;
};

export type RoomStateJSON = {
	bysy: boolean;
	exist: boolean;
	game: string;
	player_quantity: number;
};

export type DeviceInfoJSON = {
	display: {
		height: number;
		width: number;
	};
	control: {
		gamepad: boolean;
		keyboard: boolean;
		mouse: boolean;
		touch: boolean;
	};
};
