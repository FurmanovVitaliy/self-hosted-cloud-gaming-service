import { Game } from "@/types/types";
import gameApi from "@/api/game";
import log from "@/common/log";
import { GameElementType, GameElementBuilder } from "@/components/game-element";

export type Rooms = {
	id: string;
	game_id: string;
	peer: any;
};

class Cache<T> {
	private cache: Map<string, { value: T; ttl: number; filter: string }>;
	private defaultTTL: number;

	constructor(defaultTTL: number = Infinity) {
		this.cache = new Map();
		this.defaultTTL = defaultTTL;
	}

	set(key: string, value: T, ttl?: number, filter?: string): void {
		const expire = Date.now() + (ttl || this.defaultTTL);
		const queryFilter = filter || "none";
		this.cache.set(key, { value, ttl: expire, filter: queryFilter });
	}

	get(key: string): T | undefined {
		const cached = this.cache.get(key);
		if (!cached) return undefined;
		if (cached.ttl < Date.now()) {
			this.cache.delete(key);
			return undefined;
		}
		return cached.value;
	}
	getMap(): Map<string, { value: T; ttl: number }> {
		return this.cache;
	}
	getArr(): T[] {
		const arr: T[] = [];
		const now = Date.now();
		this.cache.forEach((entry, key) => {
			if (entry.ttl >= now) {
				arr.push(entry.value);
			} else {
				this.cache.delete(key);
			}
		});
		return arr;
	}
	delete(key: string): void {
		this.cache.delete(key);
	}
	clear(): void {
		this.cache.clear();
	}
	has(key: string): boolean {
		return this.cache.has(key);
	}
}

export const roomCacheALL = new Cache<Rooms>(180000);

export const gameCacheALL = new Cache<Game>(Infinity);
export const gameElementsCardsCache = new Cache<HTMLElement>(Infinity);
export const gameElementsSearchCache = new Cache<HTMLElement>(Infinity);
export const gameElementsFullscreenCache = new Cache<HTMLElement>(Infinity);

export async function initializeGameDataCache() {
	try {
		const games = (await gameApi.getAllGames()) as Game[];
		games.forEach((game) => {
			gameCacheALL.set(game.id, game);
		});
		log.debug("Game cache initialized");
	} catch (error) {
		log.error(`Failed to initialize game cache: ${error}`);
		gameCacheALL.set("empty", [] as any); // prevent re-fetching
	}
}
export function initializeGameElemsSearchCache() {
	try {
		const games = gameCacheALL.getArr();
		if (games.length === 0) {
			log.warn("Game cache is empty, skipping search elements initialization");
			return;
		}
		const builder = new GameElementBuilder(GameElementType.search);
		games.forEach((game) => {
			const gameElement = builder.addData(game).build();
			gameElementsSearchCache.set(game.id, gameElement);
		});
		log.debug("Game elements search cache initialized");
	} catch (error) {
		log.error(`Failed to initialize game elements search cache: ${error}`);
	}
}
export function initializeGameElemsTileCache() {
	try {
		const games = gameCacheALL.getArr();
		const builder = new GameElementBuilder(GameElementType.tile);
		games.forEach((game) => {
			const gameElement = builder.addData(game).build();
			gameElementsCardsCache.set(game.id, gameElement);
		});
		log.debug("Game elements slider cache initialized");
	} catch (error) {
		log.error(`Failed to initialize game elements slider cache: ${error}`);
	}
}
export function initializeGameElemsFullscreenCache() {
	try {
		const games = gameCacheALL.getArr();
		const builder = new GameElementBuilder(GameElementType.fullscreen);
		games.forEach((game) => {
			const gameElement = builder.addData(game).build();
			gameElementsFullscreenCache.set(game.id, gameElement);
		});
		log.debug("Game elements fullscreen cache initialized");
	} catch (error) {
		log.error(`Failed to initialize game elements fullscreen cache: ${error}`);
	}
}