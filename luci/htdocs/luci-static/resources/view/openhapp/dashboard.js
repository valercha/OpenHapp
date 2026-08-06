'use strict';
'require view';
'require uci';

return view.extend({
	title: _('OpenHapp Dashboard'),

	load: function () {
		return uci.load('openhapp');
	},

	render: function (data) {
		var section = data && data.main ? data.main : null;
		var enabled = section && section.enabled ? section.enabled : '0';
		var engine = section && section.engine ? section.engine : 'xray';
		var mode = section && section.mode ? section.mode : 'proxy';
		var listen = section && section.listen ? section.listen : '127.0.0.1:0';
		var subscription = section && section.subscription ? section.subscription : '';

		return E('div', { 'class': 'cbi-map' }, [
			E('h2', {}, _('OpenHapp Dashboard')),
			E('div', { 'class': 'cbi-section' }, [
				E('p', {}, [ E('strong', {}, _('Runtime: ')), enabled === '1' ? _('enabled') : _('disabled') ]),
				E('p', {}, [ E('strong', {}, _('Engine: ')), engine ]),
				E('p', {}, [ E('strong', {}, _('Mode: ')), mode ]),
				E('p', {}, [ E('strong', {}, _('Listen: ')), listen ]),
				E('p', {}, [ E('strong', {}, _('Subscription: ')), subscription || _('none') ])
			])
		]);
	}
});