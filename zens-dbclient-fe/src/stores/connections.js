import { writable } from "svelte/store";

export const connections = writable([]);
export const activeConnectionId = writable(0);
export const isLsConnectionsFetched = writable(false);
