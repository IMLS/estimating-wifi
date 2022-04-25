import PanelComponent from './panel.vue';

export default {
	id: 'wifi-minutes',
	name: 'Wifi',
	icon: 'Wifi',
	description: 'Wi-Fi minutes served',
	component: PanelComponent,
	options: ({ options }) => {
		console.log(options);
		return [
			{
				field: 'tableName',
				type: 'string',
				name: 'Collection',
				meta: {
					interface: 'system-collection',
					options: {
						includeSystem: false,
					},
				},
			},
			{
				field: 'fcfsId',
				name: 'FCFS ID',
				type: 'string',
				meta: {
					interface: 'input',
					options: {
						placeholder: 'CA5678-999',
					},
				},
			},
		];
	},
	minWidth: 20,
	minHeight: 15,
};
