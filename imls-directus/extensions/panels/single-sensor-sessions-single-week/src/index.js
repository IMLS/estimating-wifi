import PanelComponent from './panel.vue';

export default {
	id: 'single-sensor-sessions-single-week',
	name: 'Sessions/Week per sensor',
	icon: 'calendar_month',
	description: 'Single Sensor Sessions in a Single Week',
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
				name: 'Sensor name',
				type: 'string',
				meta: {
					interface: 'select-dropdown',
					options: {
						//TODO: Use collection data as list
						choices: [ 
							{
								text: "springfield",
								value: "ME8675-309"
							},
							{
								text: "in-op",
								value: "GA0027-004"
							},
							{
								text: "rpi03",
								value: "GA0058-005"
							}
						]
					},
				},
			},			
		];
	},
	minWidth: 20,
	minHeight: 15,
};
