'use strict';
'require view';
'require fs';
'require uci';

return view.extend({
	title: _('OpenHapp Dashboard'),

	load: function () {
		return Promise.all([
			uci.load('openhapp')
		]);
	},

	render: function (data) {
		var section = data[0] && data[0].openhapp ? data[0].openhapp.main : null;
		var enabled = section && section.enabled ? section.enabled : '0';
		var engine = section && section.engine ? section.engine : 'xray';
		var mode = section && section.mode ? section.mode : 'proxy';
		var listen = section && section.listen ? section.listen : '127.0.0.1:0';
		var subscription = section && section.subscription ? section.subscription : '';

		var html = [];
		html.push('<div class="cbi-map">');
		html.push('<h2>OpenHapp Dashboard</h2>');
		html.push('<div class="cbi-section">');
		html.push('<p><strong>Runtime:</strong> ' + (enabled === '1' ? 'enabled' : 'disabled') + '</p>');
		html.push('<p><strong>Engine:</strong> ' + engine + '</p>');
		html.push('<p><strong>Mode:</strong> ' + mode + '</p>');
		html.push('<p><strong>Listen:</strong> ' + listen + '</p>');
		html.push('<p><strong>Subscription:</strong> ' + (subscription || 'none') + '</p>');
		html.push('</div>');
		html.push('</div>');

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
