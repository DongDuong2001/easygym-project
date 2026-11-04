import { page } from '$app/state';

const EMPTY = { categories: [], products: [], courses: [] };

/** Catalog cards loaded by +layout.server.js. Read it inside $derived so it updates on navigation. */
export const catalog = () => page.data.catalog ?? EMPTY;
