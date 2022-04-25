import { defineDisplay } from '@directus/extensions-sdk';
import TimestampComponent from './display.vue';

export default defineDisplay({
	id: 'unix-timestamp',
	name: 'Unix Timestamp',
	icon: 'box',
	description: 'Convert unix timestamps to human-readable datetimes',
	component: ({ value }) => new Date(value * 1000).toISOString(),
	options: null,
	types: ['string'],
});
